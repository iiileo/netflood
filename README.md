# NetFlood - 带宽测试工具

NetFlood 是一个高性能的带宽测试工具，支持多协程并发下载，并实时统计下载速度。

## 功能特性

- ✅ 从配置文件或 API 获取下载任务
- ✅ 支持指定 IP 地址进行下载（绕过 DNS 解析）
- ✅ 多协程并发下载
- ✅ 实时速度统计（每秒更新）
- ✅ 不保存文件到硬盘（仅用于带宽测试）
- ✅ 优雅退出（Ctrl+C）
- ✅ 使用 resty 库实现

## 配置文件

配置文件位于 `config/config.yml`：

```yaml
# API 接口地址，用于获取下载链接列表
api: http://yd.xingshangyun.com/bh.php?bh=2

# 同时下载的协程数量
goroutines: 5
```

## 下载链接格式

从 API 返回或 demo.txt 文件中的格式为：

```
IP地址,下载链接URL
```

示例：
```
183.214.139.130,https://imtt2.dd.qq.com/sjy.00008/sjy.00001/16891/apk/A3607F63DD5A13C26A276D5141032ED0.apk
111.62.48.158,https://s2.g.mi.com/523b71ac1ec2f923aeb500f167760b08/1761574912/download/AppStore/com.tencent.hyrzol.apk
```

## 使用方法

### 编译

**使用 Makefile（推荐）：**

查看所有可用命令：
```bash
make help
```

编译当前平台版本：
```bash
make build
```

编译特定平台版本：
```bash
make linux-amd64      # Linux AMD64
make linux-arm64      # Linux ARM64
make darwin-amd64     # macOS Intel
make darwin-arm64     # macOS Apple Silicon
make windows-amd64    # Windows AMD64
```

一键编译所有平台版本：
```bash
make cross-compile
```

**使用 Go 命令：**
```bash
go build -o netflood ./cmd
```

### 运行

**从 API 获取下载任务：**
```bash
./netflood
```

**使用自定义配置文件：**
```bash
./netflood -config /path/to/config.yml
```

**使用 demo.txt 文件：**
```bash
./netflood -demo
```

### 命令行参数

- `-config <path>`: 指定配置文件路径（默认：`config/config.yml`）
- `-demo`: 使用 demo.txt 文件而不是 API 获取下载任务

## 输出

### 控制台输出

程序会实时显示：
- 当前下载速度（MB/s）
- 平均下载速度（MB/s）
- 总下载量（MB）
- 每个工作协程的下载完成情况

示例：
```
配置加载成功: API=http://yd.xingshangyun.com/bh.php?bh=2, 协程数=5
从 API 加载下载任务: http://yd.xingshangyun.com/bh.php?bh=2
成功加载 3 个下载任务
  任务 1: IP=183.214.139.130, URL=https://imtt2.dd.qq.com/sjy.00008/sjy.00001/16891/...
  任务 2: IP=111.62.48.158, URL=https://s2.g.mi.com/523b71ac1ec2f923aeb500f167760b08/...
  任务 3: IP=39.134.236.159, URL=https://apkverifywr-v6dl.vivo.com.cn/appstore/...

开始下载，使用 5 个协程...
速度统计将保存到 ./speed 文件
按 Ctrl+C 停止下载

[速度统计] 当前速度: 15.42 MB/s | 平均速度: 12.58 MB/s | 总下载: 125.80 MB
[Worker 0] 下载完成 https://imtt2.dd.qq.com/...
[Worker 1] 下载完成 https://s2.g.mi.com/...
```

### 速度文件

速度统计会实时保存到 `./speed` 文件，格式为：

```
2025-10-23 14:30:01 | 当前速度: 15.42 MB/s | 平均速度: 12.58 MB/s | 总下载: 125.80 MB
2025-10-23 14:30:02 | 当前速度: 18.91 MB/s | 平均速度: 13.45 MB/s | 总下载: 144.71 MB
```

## 技术实现

- **并发控制**: 使用 Go 协程池实现并发下载
- **IP 绑定**: 通过自定义 HTTP Transport 的 DialContext 实现指定 IP 访问
- **速度统计**: 使用 atomic.Int64 原子操作统计下载字节数
- **内存优化**: 使用流式读取，不将文件保存到硬盘
- **优雅退出**: 使用 context 实现信号处理

## 依赖

- Go 1.24.2+
- resty.dev/v3 - HTTP 客户端库
- gopkg.in/yaml.v3 - YAML 解析库

## 许可证

MIT License

