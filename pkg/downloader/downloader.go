package downloader

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// DownloadTask 下载任务
type DownloadTask struct {
	IP  string // 指定的IP地址
	URL string // 下载链接
}

// Downloader 下载器
type Downloader struct {
	client          *http.Client
	tasks           []DownloadTask
	goroutines      int
	bytesDownloaded atomic.Int64 // 已下载的字节数
	speedFile       *os.File     // 速度文件
	mu              sync.Mutex
}

// New 创建新的下载器
func New(goroutines int) *Downloader {
	return &Downloader{
		client:     &http.Client{Timeout: 30 * time.Second},
		goroutines: goroutines,
	}
}

// LoadTasksFromAPI 从API加载下载任务
func (d *Downloader) LoadTasksFromAPI(apiURL string) error {
	// 使用 http 请求API
	resp, err := d.client.Get(apiURL)
	if err != nil {
		return fmt.Errorf("请求API失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %w", err)
	}

	// 解析响应内容
	return d.parseTasksFromContent(string(body))
}

// LoadTasksFromFile 从文件加载下载任务
func (d *Downloader) LoadTasksFromFile(filepath string) error {
	// 读取文件
	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("读取文件失败: %w", err)
	}

	return d.parseTasksFromContent(string(data))
}

// parseTasksFromContent 从内容解析下载任务
func (d *Downloader) parseTasksFromContent(content string) error {
	scanner := bufio.NewScanner(strings.NewReader(content))
	var tasks []DownloadTask

	// 逐行解析
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// 按逗号分割 IP 和 URL
		parts := strings.SplitN(line, ",", 2)
		if len(parts) != 2 {
			continue
		}

		tasks = append(tasks, DownloadTask{
			IP:  strings.TrimSpace(parts[0]),
			URL: strings.TrimSpace(parts[1]),
		})
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("解析内容失败: %w", err)
	}

	d.tasks = tasks
	return nil
}

// Start 开始下载（循环模式）
func (d *Downloader) Start(ctx context.Context) error {
	if len(d.tasks) == 0 {
		return fmt.Errorf("没有下载任务")
	}

	// 打开速度文件
	var err error
	d.speedFile, err = os.Create("./speed")
	if err != nil {
		return fmt.Errorf("创建速度文件失败: %w", err)
	}
	defer d.speedFile.Close()

	// 启动速度统计协程
	go d.reportSpeed(ctx)

	// 创建任务通道（带缓冲，用于循环发送任务）
	taskChan := make(chan DownloadTask, d.goroutines*2)

	// 启动任务分发协程（循环发送任务）
	go func() {
		defer close(taskChan)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				// 循环发送所有任务
				for _, task := range d.tasks {
					select {
					case <-ctx.Done():
						return
					case taskChan <- task:
						// 任务已发送，继续
					}
				}
			}
		}
	}()

	// 启动工作协程
	var wg sync.WaitGroup
	for i := 0; i < d.goroutines; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			d.worker(ctx, workerID, taskChan)
		}(i)
	}

	// 等待所有工作协程完成（只有在 ctx 被取消时才会结束）
	wg.Wait()

	// 输出最终统计
	fmt.Println("\n正在保存最终统计数据...")
	d.printFinalStats()

	return nil
}

// printFinalStats 输出最终统计信息
func (d *Downloader) printFinalStats() {
	totalBytes := d.bytesDownloaded.Load()
	totalMB := float64(totalBytes) / 1024 / 1024

	d.mu.Lock()
	defer d.mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	finalLine := fmt.Sprintf("\n%s | ========== 下载结束 ==========\n", timestamp)
	finalLine += fmt.Sprintf("%s | 总下载量: %.2f MB (%.2f GB)\n", timestamp, totalMB, totalMB/1024)

	fmt.Print(finalLine)
	d.speedFile.WriteString(finalLine)
	d.speedFile.Sync()
}

// worker 工作协程
func (d *Downloader) worker(ctx context.Context, workerID int, taskChan <-chan DownloadTask) {
	for task := range taskChan {
		// 不使用 ctx 来中断当前任务，让任务自然完成
		err := d.downloadTask(task)
		if err != nil {
			// 检查是否是状态码错误
			if strings.Contains(err.Error(), "HTTP状态码错误") {
				// 静默跳过非200状态码，不输出错误
				continue
			}
			// 其他错误正常输出
			fmt.Printf("[Worker %d] 下载失败 %s: %v\n", workerID, task.URL, err)
		} else {
			fmt.Printf("[Worker %d] 下载完成 %s\n", workerID, task.URL)
		}
	}
}

// downloadTask 下载单个任务
func (d *Downloader) downloadTask(task DownloadTask) error {
	// 解析 URL 获取域名
	parsedURL, err := url.Parse(task.URL)
	if err != nil {
		return fmt.Errorf("解析URL失败: %w", err)
	}

	// 创建自定义的 HTTP Transport，将域名解析到指定IP
	transport := &http.Transport{
		DialContext: func(dialCtx context.Context, network, addr string) (net.Conn, error) {
			// 获取端口
			_, port, err := net.SplitHostPort(addr)
			if err != nil {
				// 如果没有端口，使用默认端口
				if parsedURL.Scheme == "https" {
					port = "443"
				} else {
					port = "80"
				}
			}

			// 使用指定的IP和端口
			addr = net.JoinHostPort(task.IP, port)
			dialer := &net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}
			return dialer.DialContext(dialCtx, network, addr)
		},
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
	}

	// 创建自定义的 HTTP 客户端（直接使用 http 包，不使用 resty）
	httpClient := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second, // 设置30秒超时
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("GET", task.URL, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	// 发送请求
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP状态码错误: %d", resp.StatusCode)
	}

	// 读取响应体，但不保存到硬盘
	buf := make([]byte, 64*1024) // 64KB 缓冲区
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			// 累加下载字节数
			d.bytesDownloaded.Add(int64(n))
		}
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("读取响应失败: %w", err)
		}
	}
}

// reportSpeed 报告下载速度
func (d *Downloader) reportSpeed(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	var lastBytes int64
	startTime := time.Now()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			currentBytes := d.bytesDownloaded.Load()

			// 计算本秒下载的字节数
			bytesThisSecond := currentBytes - lastBytes
			lastBytes = currentBytes

			// 转换为 MB/s
			speedMBps := float64(bytesThisSecond) / 1024 / 1024

			// 计算总体平均速度
			elapsed := time.Since(startTime).Seconds()
			avgSpeedMBps := float64(currentBytes) / 1024 / 1024 / elapsed

			// 输出到控制台
			fmt.Printf("[速度统计] 当前速度: %.2f MB/s | 平均速度: %.2f MB/s | 总下载: %.2f MB\n",
				speedMBps, avgSpeedMBps, float64(currentBytes)/1024/1024)

			// 写入文件（覆盖模式，只保留最新的统计）
			d.mu.Lock()
			timestamp := time.Now().Format("2006-01-02 15:04:05")
			content := fmt.Sprintf("%s | 当前速度: %.2f MB/s | 平均速度: %.2f MB/s | 总下载: %.2f MB\n",
				timestamp, speedMBps, avgSpeedMBps, float64(currentBytes)/1024/1024)

			// 清空文件并写入新内容
			d.speedFile.Seek(0, 0)
			d.speedFile.Truncate(0)
			d.speedFile.WriteString(content)
			d.speedFile.Sync() // 立即刷新到磁盘
			d.mu.Unlock()
		}
	}
}

// GetTasks 获取任务列表
func (d *Downloader) GetTasks() []DownloadTask {
	return d.tasks
}
