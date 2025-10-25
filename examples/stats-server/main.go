package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// StatsData 统计数据结构
type StatsData struct {
	Name  string  `json:"name"`
	Speed float64 `json:"speed"`
	Total float64 `json:"total"`
	Time  string  `json:"time"`
}

func main() {
	http.HandleFunc("/stats", handleStats)
	http.HandleFunc("/", handleHome)

	fmt.Println("===========================================")
	fmt.Println("📊 NetFlood 统计接收服务器")
	fmt.Println("===========================================")
	fmt.Println("监听地址: http://localhost:8080")
	fmt.Println("统计接口: http://localhost:8080/stats")
	fmt.Println("")
	fmt.Println("使用方法:")
	fmt.Println("  ./netflood -demo -stats-api http://localhost:8080/stats")
	fmt.Println("")
	fmt.Println("按 Ctrl+C 停止服务器")
	fmt.Println("===========================================")
	fmt.Println("")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>NetFlood 统计服务器</title>
    <meta charset="UTF-8">
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 800px;
            margin: 50px auto;
            padding: 20px;
            background: #f5f5f5;
        }
        .container {
            background: white;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        h1 { color: #333; }
        code {
            background: #f0f0f0;
            padding: 2px 8px;
            border-radius: 3px;
            font-family: monospace;
        }
        .endpoint {
            background: #e8f5e9;
            padding: 15px;
            margin: 20px 0;
            border-left: 4px solid #4caf50;
            border-radius: 3px;
        }
        .example {
            background: #e3f2fd;
            padding: 15px;
            margin: 20px 0;
            border-left: 4px solid #2196f3;
            border-radius: 3px;
        }
        pre {
            background: #f5f5f5;
            padding: 15px;
            border-radius: 5px;
            overflow-x: auto;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>📊 NetFlood 统计接收服务器</h1>
        <p>此服务器用于接收 NetFlood 的统计数据上报。</p>
        
        <div class="endpoint">
            <h3>📍 统计接口</h3>
            <p><strong>URL:</strong> <code>POST /stats</code></p>
            <p><strong>Content-Type:</strong> <code>application/json</code></p>
        </div>

        <div class="example">
            <h3>📝 请求示例</h3>
            <pre>{
  "name": "my-server",
  "speed": 15.5,
  "total": 1024.0,
  "time": "12:00-13:00"
}</pre>
        </div>

        <div class="example">
            <h3>🚀 使用方法</h3>
            <pre>./netflood -demo -stats-api http://localhost:8080/stats</pre>
            <p>或者带时间段控制：</p>
            <pre>./netflood -d -g 20 -t "09:00-18:00" -s http://localhost:8080/stats</pre>
        </div>

        <h3>📋 字段说明</h3>
        <ul>
            <li><code>name</code> - 主机名称</li>
            <li><code>speed</code> - 平均下载速度（MB/s）</li>
            <li><code>total</code> - 总下载量（MB）</li>
            <li><code>time</code> - 时间范围</li>
        </ul>

        <p style="color: #666; margin-top: 30px;">服务器将在控制台实时显示接收到的统计数据。</p>
    </div>
</body>
</html>
`
	w.Write([]byte(html))
}

func handleStats(w http.ResponseWriter, r *http.Request) {
	// 只接受 POST 请求
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 读取请求体
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("❌ 读取请求体失败: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// 解析 JSON
	var stats StatsData
	if err := json.Unmarshal(body, &stats); err != nil {
		log.Printf("❌ 解析 JSON 失败: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// 打印统计信息
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] 📊 收到统计数据:\n", timestamp)
	fmt.Printf("  主机名: %s\n", stats.Name)
	fmt.Printf("  平均速度: %.2f MB/s\n", stats.Speed)
	fmt.Printf("  总下载量: %.2f MB (%.2f GB)\n", stats.Total, stats.Total/1024)
	fmt.Printf("  时间段: %s\n", stats.Time)
	fmt.Println("-------------------------------------------")

	// 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"status":  "success",
		"message": "Statistics received",
	}
	json.NewEncoder(w).Encode(response)
}
