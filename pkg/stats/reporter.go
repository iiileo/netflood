package stats

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// StatsData 统计数据结构
type StatsData struct {
	Name  string  `json:"name"`  // 主机名称
	Speed float64 `json:"speed"` // 平均下载速度（MB/s）
	Total float64 `json:"total"` // 总下载量（MB）
	Time  string  `json:"time"`  // 时间范围
}

// Reporter 统计数据上报器
type Reporter struct {
	apiURL   string
	hostname string
	client   *http.Client
}

// NewReporter 创建统计上报器
func NewReporter(apiURL string) (*Reporter, error) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	return &Reporter{
		apiURL:   apiURL,
		hostname: hostname,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

// Report 上报统计数据
func (r *Reporter) Report(avgSpeed, totalMB float64, timeRange string) error {
	data := StatsData{
		Name:  r.hostname,
		Speed: avgSpeed,
		Total: totalMB,
		Time:  timeRange,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("序列化数据失败: %w", err)
	}

	req, err := http.NewRequest("POST", r.apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("服务器返回错误状态码: %d", resp.StatusCode)
	}

	return nil
}

// StartReporting 启动定期上报
func (r *Reporter) StartReporting(ctx context.Context, getBytesDownloaded func() int64, getStartTime func() time.Time, getTimeRange func() string) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// 获取统计数据
			totalBytes := getBytesDownloaded()
			startTime := getStartTime()
			timeRange := getTimeRange()

			// 计算统计
			elapsed := time.Since(startTime).Seconds()
			if elapsed < 1 {
				elapsed = 1
			}

			totalMB := float64(totalBytes) / 1024 / 1024
			avgSpeedMBps := totalMB / elapsed

			// 上报数据
			if err := r.Report(avgSpeedMBps, totalMB, timeRange); err != nil {
				fmt.Printf("[统计上报] 失败: %v\n", err)
			} else {
				fmt.Printf("[统计上报] 成功: 主机=%s, 平均速度=%.2f MB/s, 总下载=%.2f MB, 时间段=%s\n",
					r.hostname, avgSpeedMBps, totalMB, timeRange)
			}
		}
	}
}

// GetHostname 获取主机名
func (r *Reporter) GetHostname() string {
	return r.hostname
}
