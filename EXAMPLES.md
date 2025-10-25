# NetFlood 使用示例

## 基础用法

### 1. 全天候下载（不限制时间）

```bash
# 使用 demo.txt 文件，全天候运行
./netflood -demo

# 使用 API，全天候运行
./netflood -api https://api.example.com/download
```

**输出示例：**
```
配置参数: API=, 协程数=12
下载时间段: 全天候运行
从 demo.txt 文件加载下载任务...
成功加载 3 个下载任务

开始下载，使用 12 个协程...
⚡ 循环下载模式：协程将不停下载任务
⚠️  按 Ctrl+C 优雅退出
```

---

## 时间段控制

### 2. 单个时间段下载

```bash
# 每天 12:00-13:00 之间下载
./netflood -demo -time 12:00-13:00

# 使用简写
./netflood -d -t 12:00-13:00
```

**输出示例（当前不在时间段内）：**
```
配置参数: API=, 协程数=12
下载时间段: 12:00-13:00 (每天重复)

⏰ 当前不在下载时间段内，等待到 12:00:00 (等待 2h30m45s)
```

**输出示例（进入时间段）：**
```
✅ 进入下载时间段，开始下载...
[速度统计] 当前速度: 18.45 MB/s | 平均速度: 15.23 MB/s | 总下载: 245.60 MB
[Worker 0] 下载完成 https://example.com/file1.apk
```

**输出示例（离开时间段）：**
```
⏰ 已超出下载时间段，停止分发新任务，等待当前任务完成...
⏰ 当前不在下载时间段内，等待到 12:00:00 (等待 23h15m30s)
```

### 3. 多个时间段下载

```bash
# 每天 9:00-12:00 和 14:00-18:00 下载
./netflood -demo -time "09:00-12:00,14:00-18:00"

# 三个时间段
./netflood -d -t "09:00-10:00,14:00-16:00,20:00-22:00"
```

**输出示例：**
```
配置参数: API=, 协程数=12
下载时间段: 09:00-12:00, 14:00-18:00 (每天重复)

⏰ 时间段控制已启用：09:00-12:00, 14:00-18:00 (每天重复)
```

### 4. 午休时间下载

```bash
# 中午休息时间下载
./netflood -demo -time 12:00-14:00
```

### 5. 夜间下载（跨天）

```bash
# 晚上 10 点到凌晨 2 点下载
./netflood -demo -time 22:00-02:00
```

---

## 协程数量控制

### 6. 自定义协程数量

```bash
# 使用 20 个协程
./netflood -demo -goroutines 20

# 使用简写
./netflood -d -g 20

# 低速测试（2 个协程）
./netflood -d -g 2

# 高速测试（50 个协程）
./netflood -d -g 50
```

---

## 统计数据上报

### 7. 启用统计上报

```bash
# 每10秒上报统计数据到API
./netflood -demo -stats-api https://api.example.com/stats

# 使用简写
./netflood -d -s https://api.example.com/stats
```

**输出示例：**
```
配置参数: API=, 协程数=12
统计上报API: https://api.example.com/stats (每10秒上报一次)
从 demo.txt 文件加载下载任务...

[速度统计] 当前速度: 18.45 MB/s | 平均速度: 15.23 MB/s | 总下载: 245.60 MB
[统计上报] 成功: 主机=my-server, 平均速度=15.23 MB/s, 总下载=245.60 MB, 时间段=全天候
```

**上报数据格式（JSON）：**
```json
{
  "name": "my-server",
  "speed": 15.23,
  "total": 245.60,
  "time": "全天候"
}
```

### 8. 时间段控制 + 统计上报

```bash
# 9:00-18:00 下载，并上报统计数据
./netflood -d -g 20 -t "09:00-18:00" -s https://api.example.com/stats
```

**输出示例：**
```
配置参数: API=, 协程数=20
下载时间段: 09:00-18:00 (每天重复)
统计上报API: https://api.example.com/stats (每10秒上报一次)

[统计上报] 成功: 主机=my-server, 平均速度=25.50 MB/s, 总下载=1500.00 MB, 时间段=09:00-18:00
```

---

## 组合使用

### 9. 工作时间带宽测试

```bash
# 工作日 9:00-12:00 和 14:00-18:00，使用 30 个协程
./netflood -api https://api.example.com/download -g 30 -t "09:00-12:00,14:00-18:00"
```

### 10. 夜间低速测试

```bash
# 每天晚上 22:00-06:00，使用 5 个协程
./netflood -d -g 5 -t 22:00-06:00
```

### 11. 午休时间高速测试

```bash
# 中午 12:00-14:00，使用 50 个协程
./netflood -d -g 50 -t 12:00-14:00
```

### 12. 完整功能组合

```bash
# 时间段 + 统计上报 + 自定义协程
./netflood -api https://api.example.com/download -g 30 -t "09:00-18:00" -s https://api.example.com/stats
```

---

## 优雅退出

### 13. 等待任务完成后退出

```bash
# 运行程序
./netflood -demo

# 按 Ctrl+C（第一次）
^C
收到退出信号，等待当前下载任务完成...
⚠️  如需强制退出，请再次按 Ctrl+C

# 等待当前任务完成
✅ 下载已停止，程序退出
```

### 14. 强制立即退出

```bash
# 运行程序
./netflood -demo

# 按 Ctrl+C（第一次）
^C
收到退出信号，等待当前下载任务完成...
⚠️  如需强制退出，请再次按 Ctrl+C

# 按 Ctrl+C（第二次）
^C
收到强制退出信号，立即退出...
```

---

## 实际应用场景

### 场景 1：办公室带宽测试

在办公室上班时间测试网络带宽：

```bash
./netflood -api https://api.example.com/download -g 20 -t "09:00-18:00"
```

### 场景 2：家庭夜间测试

利用夜间空闲时间测试家庭网络：

```bash
./netflood -demo -g 30 -t "00:00-06:00"
```

### 场景 3：分时段测试

避开高峰期，在低峰期测试：

```bash
./netflood -d -g 15 -t "02:00-05:00,14:00-16:00,22:00-23:00"
```

### 场景 4：服务器带宽监控

服务器全天候运行，监控带宽稳定性：

```bash
./netflood -api https://api.example.com/download -g 10
```

### 场景 5：限时压力测试

在指定时间段进行高强度压力测试：

```bash
./netflood -d -g 100 -t "03:00-04:00"
```

### 场景 6：集中监控多服务器

在多台服务器上部署，统一上报到监控中心：

```bash
# 服务器A
./netflood -d -g 20 -t "09:00-18:00" -s https://monitor.example.com/api/stats

# 服务器B
./netflood -d -g 20 -t "09:00-18:00" -s https://monitor.example.com/api/stats

# 服务器C
./netflood -d -g 20 -t "09:00-18:00" -s https://monitor.example.com/api/stats
```

监控中心可以通过主机名区分不同服务器的统计数据。

### 场景 7：性能监控和告警

配合监控系统使用，自动告警：

```bash
# 启用统计上报
./netflood -api https://api.example.com/download -g 30 -s https://monitor.example.com/api/stats
```

监控系统可以：
- 实时监控各服务器下载速度
- 检测网络异常（速度突降）
- 统计总带宽使用情况
- 生成报表和图表

---

## 速度文件监控

程序运行时，速度统计会实时保存到 `./speed` 文件：

```bash
# 实时监控速度
watch -n 1 cat speed

# 或使用 tail
tail -f speed
```

**输出示例：**
```
2025-10-25 14:30:01 | 当前速度: 18.45 MB/s | 平均速度: 15.23 MB/s | 总下载: 245.60 MB
```

---

## 常见错误处理

### 错误 1：无效的时间段格式

```bash
./netflood -demo -time 12:00
```

**错误输出：**
```
解析时间段失败: 无效的时间段格式: 12:00 (应为 HH:MM-HH:MM)
时间段格式示例: -time 12:00-13:00,14:00-15:00
```

**正确用法：**
```bash
./netflood -demo -time 12:00-13:00
```

### 错误 2：小时超出范围

```bash
./netflood -demo -time 25:00-26:00
```

**错误输出：**
```
解析时间段失败: 解析开始时间失败 25:00: 小时必须在 0-23 之间: 25
```

### 错误 3：分钟超出范围

```bash
./netflood -demo -time 12:70-13:00
```

**错误输出：**
```
解析时间段失败: 解析开始时间失败 12:70: 分钟必须在 0-59 之间: 70
```

---

## 高级技巧

### 技巧 1：后台运行

```bash
# 使用 nohup 后台运行
nohup ./netflood -d -g 20 -t "09:00-18:00" > netflood.log 2>&1 &

# 查看日志
tail -f netflood.log
```

### 技巧 2：使用 systemd 服务（Linux）

创建服务文件 `/etc/systemd/system/netflood.service`：

```ini
[Unit]
Description=NetFlood Bandwidth Test
After=network.target

[Service]
Type=simple
User=your-user
WorkingDirectory=/path/to/netflood
ExecStart=/path/to/netflood/netflood -demo -g 20 -t "09:00-18:00"
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

启动服务：
```bash
sudo systemctl daemon-reload
sudo systemctl enable netflood
sudo systemctl start netflood
sudo systemctl status netflood
```

### 技巧 3：Docker 部署

创建 `Dockerfile`：

```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o netflood ./cmd

FROM alpine:latest
COPY --from=builder /app/netflood /usr/local/bin/
COPY demo.txt /app/
WORKDIR /app
CMD ["netflood", "-demo", "-g", "20", "-t", "09:00-18:00"]
```

构建并运行：
```bash
docker build -t netflood .
docker run -d --name netflood netflood
```

---

## 性能建议

1. **协程数量**：根据网络带宽调整，通常 10-30 个协程足够
2. **时间段设置**：避免高峰期，选择网络空闲时段
3. **监控资源**：定期检查 CPU 和内存使用情况
4. **日志管理**：定期清理日志文件，避免磁盘空间不足

---

## 常见问题

**Q: 程序会消耗多少带宽？**
A: 取决于协程数量和网络速度，通常可以跑满带宽。

**Q: 下载的文件保存在哪里？**
A: 不保存，只用于带宽测试，数据直接丢弃。

**Q: 可以全天候运行吗？**
A: 可以，不设置 `-time` 参数即可全天候运行。

**Q: 时间段每天会自动重复吗？**
A: 是的，设置的时间段每天自动重复，无需手动重启。

**Q: 如何修改时间段？**
A: 停止程序（Ctrl+C），然后用新的时间段参数重新启动。

