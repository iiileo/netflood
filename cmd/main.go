package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dora-exku/netflood/pkg/downloader"
	"github.com/dora-exku/netflood/pkg/timerange"
)

func main() {
	// 定义命令行参数（支持简写）
	api := flag.String("api", "", "API 接口地址")
	apiShort := flag.String("a", "", "API 接口地址（简写）")

	goroutines := flag.Int("goroutines", 12, "同时下载的协程数量")
	goroutinesShort := flag.Int("g", 12, "同时下载的协程数量（简写）")

	useDemo := flag.Bool("demo", false, "使用demo.txt文件而不是API")
	useDemoShort := flag.Bool("d", false, "使用demo.txt文件而不是API（简写）")

	timeRangeStr := flag.String("time", "", "下载时间段，格式: HH:MM-HH:MM,HH:MM-HH:MM (例如: 12:00-13:00,14:00-15:00)")
	timeRangeShort := flag.String("t", "", "下载时间段（简写）")

	flag.Parse()

	// 使用简写参数值（如果设置了简写，优先使用简写）
	finalAPI := *api
	if *apiShort != "" {
		finalAPI = *apiShort
	}

	finalGoroutines := *goroutines
	if *goroutinesShort != 12 {
		finalGoroutines = *goroutinesShort
	}

	finalUseDemo := *useDemo || *useDemoShort

	finalTimeRange := *timeRangeStr
	if *timeRangeShort != "" {
		finalTimeRange = *timeRangeShort
	}

	fmt.Printf("配置参数: API=%s, 协程数=%d\n", finalAPI, finalGoroutines)

	// 解析时间段
	trm, err := timerange.NewTimeRangeManager(finalTimeRange)
	if err != nil {
		fmt.Printf("解析时间段失败: %v\n", err)
		fmt.Println("时间段格式示例: -time 12:00-13:00,14:00-15:00")
		os.Exit(1)
	}
	if trm.IsEnabled() {
		fmt.Printf("下载时间段: %s (每天重复)\n", trm.String())
	} else {
		fmt.Println("下载时间段: 全天候运行")
	}

	// 创建下载器
	dl := downloader.New(finalGoroutines)

	// 设置时间段管理器
	dl.SetTimeRangeManager(trm)

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
		fmt.Println("\n收到退出信号，等待当前下载任务完成...")
		fmt.Println("⚠️  如需强制退出，请再次按 Ctrl+C")
		cancel()

		// 等待第二次信号，强制退出
		<-sigChan
		fmt.Println("\n收到强制退出信号，立即退出...")
		os.Exit(1)
	}()

	// 开始下载
	fmt.Printf("\n开始下载，使用 %d 个协程...\n", finalGoroutines)
	fmt.Println("速度统计将保存到 ./speed 文件")
	fmt.Println("⚡ 循环下载模式：协程将不停下载任务")
	if trm.IsEnabled() {
		fmt.Printf("⏰ 时间段控制已启用：%s (每天重复)\n", trm.String())
	}
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
