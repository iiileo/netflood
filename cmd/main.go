package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dora-exku/netflood/pkg/downloader"
)

func main() {
	// 定义命令行参数（支持简写）
	api := flag.String("api", "https://api.market.huajuhe.com/a", "API 接口地址")
	apiShort := flag.String("a", "https://api.market.huajuhe.com/a", "API 接口地址（简写）")

	goroutines := flag.Int("goroutines", 12, "同时下载的协程数量")
	goroutinesShort := flag.Int("g", 12, "同时下载的协程数量（简写）")

	useDemo := flag.Bool("demo", false, "使用demo.txt文件而不是API")
	useDemoShort := flag.Bool("d", false, "使用demo.txt文件而不是API（简写）")

	flag.Parse()

	// 使用简写参数值（如果设置了简写，优先使用简写）
	finalAPI := *api
	if *apiShort != "https://api.market.huajuhe.com/a" {
		finalAPI = *apiShort
	}

	finalGoroutines := *goroutines
	if *goroutinesShort != 12 {
		finalGoroutines = *goroutinesShort
	}

	finalUseDemo := *useDemo || *useDemoShort

	fmt.Printf("配置参数: API=%s, 协程数=%d\n", finalAPI, finalGoroutines)

	// 创建下载器
	dl := downloader.New(finalGoroutines)

	// 加载下载任务
	if finalUseDemo {
		// 从demo.txt文件加载
		fmt.Println("从 demo.txt 文件加载下载任务...")
		if err := dl.LoadTasksFromFile("demo.txt"); err != nil {
			fmt.Printf("加载任务失败: %v\n", err)
			os.Exit(1)
		}
	} else {
		// 从API加载
		fmt.Printf("从 API 加载下载任务: %s\n", finalAPI)
		if err := dl.LoadTasksFromAPI(finalAPI); err != nil {
			fmt.Printf("加载任务失败: %v\n", err)
			os.Exit(1)
		}
	}

	tasks := dl.GetTasks()
	fmt.Printf("成功加载 %d 个下载任务\n", len(tasks))

	// 显示任务列表
	for i, task := range tasks {
		fmt.Printf("  任务 %d: IP=%s, URL=%s\n", i+1, task.IP, task.URL[:min(60, len(task.URL))]+"...")
	}

	// 创建上下文，用于优雅退出
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 监听退出信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\n收到退出信号，正在停止...")
		cancel()
	}()

	// 开始下载
	fmt.Printf("\n开始下载，使用 %d 个协程...\n", finalGoroutines)
	fmt.Println("速度统计将保存到 ./speed 文件")
	fmt.Println("⚡ 循环下载模式：协程将不停下载任务")
	fmt.Println("⚠️  按 Ctrl+C 优雅退出")
	fmt.Println()

	if err := dl.Start(ctx); err != nil {
		if err != context.Canceled {
			fmt.Printf("\n❌ 下载过程出错: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Println("\n✅ 下载已停止，程序退出")
}

// min 返回两个整数中的最小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
