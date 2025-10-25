# 统计数据上报 API 文档

## 概述

NetFlood 支持将统计数据自动上报到指定的 HTTP API 接口，方便集中监控和管理多台服务器的带宽测试情况。

## 功能特性

- ✅ 每10秒自动上报一次
- ✅ 包含主机名、平均速度、总下载量、时间范围
- ✅ JSON 格式数据
- ✅ HTTP POST 请求
- ✅ 失败自动重试（下次周期）
- ✅ 可选启用（不设置则不上报）

## 命令行参数

```bash
-stats-api <URL>    # 统计上报API地址（完整形式）
-s <URL>            # 统计上报API地址（简写形式）
```

## 使用示例

### 基础使用

```bash
# 启用统计上报
./netflood -demo -stats-api http://your-server.com/api/stats

# 使用简写
./netflood -d -s http://your-server.com/api/stats
```

### 组合使用

```bash
# 带时间段控制
./netflood -d -t "09:00-18:00" -s http://your-server.com/api/stats

# 完整参数
./netflood -api https://cdn.example.com/list -g 30 -t "09:00-12:00,14:00-18:00" -s http://monitor.example.com/stats
```

## API 接口规范

### 请求

**方法：** `POST`

**Content-Type：** `application/json`

**请求体格式：**

```json
{
  "name": "主机名称",
  "speed": 15.5,
  "total": 1024.0,
  "time": "12:00-13:00, 14:00-15:00"
}
```

**字段说明：**

| 字段 | 类型 | 说明 | 示例 |
|------|------|------|------|
| `name` | string | 主机名称（自动获取） | `"my-server"` |
| `speed` | float64 | 平均下载速度（MB/s） | `15.5` |
| `total` | float64 | 总下载量（MB） | `1024.0` |
| `time` | string | 时间范围（来自 -time 参数） | `"12:00-13:00"` 或 `"全天候"` |

### 响应

**成功响应（推荐）：**

```json
{
  "status": "success",
  "message": "Statistics received"
}
```

**状态码：**
- `200 OK` - 成功接收
- `400 Bad Request` - 请求格式错误
- `401 Unauthorized` - 认证失败（如果需要）
- `500 Internal Server Error` - 服务器错误

## 示例：接收服务器实现

### Go 语言实现

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "log"
)

type StatsData struct {
    Name  string  `json:"name"`
    Speed float64 `json:"speed"`
    Total float64 `json:"total"`
    Time  string  `json:"time"`
}

func handleStats(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var stats StatsData
    if err := json.NewDecoder(r.Body).Decode(&stats); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    // 处理统计数据（保存到数据库、发送告警等）
    fmt.Printf("收到统计: %s - %.2f MB/s, %.2f MB\n", 
        stats.Name, stats.Speed, stats.Total)

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "status": "success",
    })
}

func main() {
    http.HandleFunc("/stats", handleStats)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### Python (Flask) 实现

```python
from flask import Flask, request, jsonify
from datetime import datetime

app = Flask(__name__)

@app.route('/stats', methods=['POST'])
def receive_stats():
    data = request.get_json()
    
    # 验证数据
    if not all(k in data for k in ['name', 'speed', 'total', 'time']):
        return jsonify({'error': 'Missing fields'}), 400
    
    # 处理统计数据
    print(f"[{datetime.now()}] 收到统计:")
    print(f"  主机: {data['name']}")
    print(f"  速度: {data['speed']} MB/s")
    print(f"  总量: {data['total']} MB")
    print(f"  时段: {data['time']}")
    
    # 保存到数据库...
    # save_to_database(data)
    
    return jsonify({'status': 'success'}), 200

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8080)
```

### Node.js (Express) 实现

```javascript
const express = require('express');
const app = express();

app.use(express.json());

app.post('/stats', (req, res) => {
    const { name, speed, total, time } = req.body;
    
    // 验证数据
    if (!name || speed === undefined || total === undefined || !time) {
        return res.status(400).json({ error: 'Missing fields' });
    }
    
    // 处理统计数据
    console.log(`[${new Date().toISOString()}] 收到统计:`);
    console.log(`  主机: ${name}`);
    console.log(`  速度: ${speed} MB/s`);
    console.log(`  总量: ${total} MB`);
    console.log(`  时段: ${time}`);
    
    // 保存到数据库...
    // saveToDatabase(req.body);
    
    res.json({ status: 'success' });
});

app.listen(8080, () => {
    console.log('服务器运行在 http://localhost:8080');
});
```

## 数据库存储建议

### MySQL 表结构

```sql
CREATE TABLE netflood_stats (
    id INT AUTO_INCREMENT PRIMARY KEY,
    hostname VARCHAR(255) NOT NULL,
    avg_speed DECIMAL(10, 2) NOT NULL,
    total_mb DECIMAL(15, 2) NOT NULL,
    time_range VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_hostname (hostname),
    INDEX idx_created_at (created_at)
);
```

### PostgreSQL 表结构

```sql
CREATE TABLE netflood_stats (
    id SERIAL PRIMARY KEY,
    hostname VARCHAR(255) NOT NULL,
    avg_speed NUMERIC(10, 2) NOT NULL,
    total_mb NUMERIC(15, 2) NOT NULL,
    time_range VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_hostname ON netflood_stats(hostname);
CREATE INDEX idx_created_at ON netflood_stats(created_at);
```

### MongoDB 文档结构

```javascript
{
  "_id": ObjectId("..."),
  "hostname": "my-server",
  "avgSpeed": 15.5,
  "totalMB": 1024.0,
  "timeRange": "12:00-13:00",
  "createdAt": ISODate("2025-10-25T14:30:15Z")
}
```

## 监控和告警

### 监控指标

1. **速度监控**：检测平均速度是否低于阈值
2. **稳定性监控**：检测速度波动是否过大
3. **可用性监控**：检测是否正常上报
4. **总量监控**：统计每日/每周/每月总下载量

### 告警示例

```python
# 伪代码示例
def check_and_alert(stats):
    # 速度过低告警
    if stats['speed'] < 5.0:
        send_alert(f"速度过低: {stats['name']} - {stats['speed']} MB/s")
    
    # 长时间未上报告警
    if time_since_last_report(stats['name']) > 60:
        send_alert(f"服务器失联: {stats['name']}")
    
    # 异常波动告警
    if speed_variance(stats['name']) > 50:
        send_alert(f"速度波动异常: {stats['name']}")
```

## 安全建议

### 1. 身份验证

添加 API Key 验证：

```go
func handleStats(w http.ResponseWriter, r *http.Request) {
    apiKey := r.Header.Get("X-API-Key")
    if apiKey != "your-secret-key" {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    // ... 处理请求
}
```

### 2. IP 白名单

```go
allowedIPs := []string{"192.168.1.100", "10.0.0.50"}

func handleStats(w http.ResponseWriter, r *http.Request) {
    clientIP := r.RemoteAddr
    if !isAllowed(clientIP, allowedIPs) {
        http.Error(w, "Forbidden", http.StatusForbidden)
        return
    }
    // ... 处理请求
}
```

### 3. HTTPS

生产环境必须使用 HTTPS：

```bash
./netflood -d -s https://monitor.example.com/stats
```

### 4. 请求限流

防止恶意请求：

```go
import "golang.org/x/time/rate"

limiter := rate.NewLimiter(10, 20) // 每秒10个请求，突发20个

func rateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if !limiter.Allow() {
            http.Error(w, "Too many requests", http.StatusTooManyRequests)
            return
        }
        next(w, r)
    }
}
```

## 测试

### 使用 curl 测试

```bash
curl -X POST http://localhost:8080/stats \
  -H "Content-Type: application/json" \
  -d '{
    "name": "test-server",
    "speed": 25.5,
    "total": 512.0,
    "time": "09:00-18:00"
  }'
```

### 使用 NetFlood 测试

```bash
# 确保有 demo.txt 文件
echo "8.8.8.8,https://speed.cloudflare.com/__down?bytes=100000000" > demo.txt

# 运行测试
./netflood -demo -stats-api http://localhost:8080/stats
```

## 故障排查

### 常见问题

**1. 上报失败**

```
[统计上报] 失败: Post "http://...": dial tcp: connection refused
```

**解决方法：**
- 检查API服务器是否运行
- 检查URL是否正确
- 检查网络连接

**2. 无法解析主机名**

```
[统计上报] 失败: no such host
```

**解决方法：**
- 检查DNS设置
- 使用IP地址代替域名

**3. 超时**

```
[统计上报] 失败: context deadline exceeded
```

**解决方法：**
- 检查网络延迟
- 优化API服务器响应时间

## 完整示例

在项目的 `examples/stats-server/` 目录中提供了一个完整的统计接收服务器示例。

编译并运行：

```bash
cd examples/stats-server
go build -o stats-server
./stats-server
```

然后在另一个终端运行 NetFlood：

```bash
./netflood -demo -stats-api http://localhost:8080/stats
```

## 扩展应用

1. **Grafana 集成**：将数据导入 InfluxDB，使用 Grafana 可视化
2. **Prometheus 集成**：导出为 Prometheus metrics
3. **钉钉/企业微信告警**：速度异常时自动推送消息
4. **自动伸缩**：根据带宽使用情况自动调整资源
5. **报表生成**：生成日报、周报、月报

## 相关文档

- [README.md](README.md) - 项目主文档
- [EXAMPLES.md](EXAMPLES.md) - 使用示例
- [CHANGELOG.md](CHANGELOG.md) - 更新日志
- [examples/stats-server/](examples/stats-server/) - 示例服务器

