# NetFlood 统计接收服务器示例

这是一个简单的 HTTP 服务器示例，用于接收和显示 NetFlood 的统计数据上报。

## 功能

- 接收 NetFlood 统计数据
- 实时在控制台显示接收到的数据
- 提供友好的 Web 界面说明

## 编译

```bash
cd examples/stats-server
go build -o stats-server
```

## 运行

```bash
# 默认监听 8080 端口
./stats-server
```

**输出：**
```
===========================================
📊 NetFlood 统计接收服务器
===========================================
监听地址: http://localhost:8080
统计接口: http://localhost:8080/stats

使用方法:
  ./netflood -demo -stats-api http://localhost:8080/stats

按 Ctrl+C 停止服务器
===========================================
```

## 使用

### 1. 启动统计服务器

在一个终端窗口中运行：

```bash
./stats-server
```

### 2. 运行 NetFlood 并启用统计上报

在另一个终端窗口中运行：

```bash
cd ../..
./build/netflood -demo -stats-api http://localhost:8080/stats
```

### 3. 查看统计数据

统计服务器会每 10 秒接收一次数据，并在控制台显示：

```
[2025-10-25 14:30:15] 📊 收到统计数据:
  主机名: my-server
  平均速度: 15.23 MB/s
  总下载量: 245.60 MB (0.24 GB)
  时间段: 全天候
-------------------------------------------
[2025-10-25 14:30:25] 📊 收到统计数据:
  主机名: my-server
  平均速度: 18.45 MB/s
  总下载量: 430.20 MB (0.42 GB)
  时间段: 全天候
-------------------------------------------
```

## Web 界面

在浏览器中访问 `http://localhost:8080` 可以看到使用说明和 API 文档。

## API 接口

### POST /stats

接收统计数据

**请求头：**
```
Content-Type: application/json
```

**请求体：**
```json
{
  "name": "my-server",
  "speed": 15.5,
  "total": 1024.0,
  "time": "12:00-13:00"
}
```

**响应：**
```json
{
  "status": "success",
  "message": "Statistics received"
}
```

## 测试

使用 curl 测试：

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

## 扩展

这个示例服务器可以扩展为：

1. **数据库存储**：将统计数据保存到数据库
2. **实时监控面板**：使用 WebSocket 实时推送数据到前端
3. **告警系统**：当速度低于阈值时发送告警
4. **数据分析**：生成图表和报表
5. **多服务器监控**：汇总多台服务器的统计数据

## 示例：完整测试流程

```bash
# 终端 1: 启动统计服务器
cd examples/stats-server
go run main.go

# 终端 2: 运行 NetFlood
cd ../..
./build/netflood -demo -stats-api http://localhost:8080/stats

# 终端 3: 手动测试 API
curl -X POST http://localhost:8080/stats \
  -H "Content-Type: application/json" \
  -d '{"name":"test","speed":20.5,"total":500,"time":"全天候"}'
```

## 生产环境部署建议

1. **使用环境变量配置端口**
2. **添加身份验证**（API Key 或 Token）
3. **添加 HTTPS 支持**
4. **添加请求日志**
5. **添加数据验证和清洗**
6. **使用数据库持久化**
7. **添加监控和告警**

