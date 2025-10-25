# 更新日志

## [v2.1.0] - 2025-10-25

### 🎉 新增功能

#### 统计数据上报
- ✅ 支持将统计数据自动上报到指定API
- ✅ 每10秒自动上报一次
- ✅ 包含主机名、平均速度、总下载量、时间范围
- ✅ 新增 `-stats-api` / `-s` 参数

**使用示例：**
```bash
# 启用统计上报
./netflood -demo -stats-api https://api.example.com/stats

# 组合使用
./netflood -d -g 20 -t "09:00-18:00" -s https://api.example.com/stats
```

**上报数据格式（JSON）：**
```json
{
  "name": "主机名称",
  "speed": 15.5,     // 平均下载速度（MB/s）
  "total": 1024.0,   // 总下载量（MB）
  "time": "12:00-13:00"  // 时间范围
}
```

### 🔧 改进

1. **监控能力增强**
   - 支持集中监控多台服务器
   - 便于统计分析和告警

2. **文档完善**
   - 更新 README.md 添加统计上报说明
   - 更新 EXAMPLES.md 添加统计上报示例
   - 新增实际应用场景

### 📝 技术实现

**统计上报架构：**
```
Reporter (pkg/stats)
├── 获取主机名
├── 每10秒上报一次
├── HTTP POST 发送 JSON 数据
└── 错误处理和重试

Downloader
├── 集成 Reporter
├── 提供统计数据回调
└── 自动启动上报协程
```

### 🧪 测试

新增单元测试：
- ✅ Reporter 创建测试
- ✅ 数据上报测试
- ✅ 服务器错误处理测试
- ✅ 数据格式验证测试

**测试结果：**
```
=== RUN   TestNewReporter
--- PASS: TestNewReporter (0.00s)
=== RUN   TestReporter_Report
--- PASS: TestReporter_Report (0.00s)
=== RUN   TestReporter_Report_ServerError
--- PASS: TestReporter_Report_ServerError (0.00s)
PASS
```

### 📦 依赖变化

- 无新增外部依赖
- 仅使用 Go 标准库

---

## [v2.0.0] - 2025-10-25

### 🎉 新增功能

#### 1. 时间段控制
- ✅ 支持设置多个下载时间段
- ✅ 时间段每天自动重复
- ✅ 不在时间段内自动休眠，到时间自动唤醒
- ✅ 支持跨天时间段（如 23:00-01:00）
- ✅ 新增 `-time` / `-t` 参数

**使用示例：**
```bash
# 单个时间段
./netflood -demo -time 12:00-13:00

# 多个时间段
./netflood -demo -time "09:00-12:00,14:00-18:00"

# 跨天时间段
./netflood -demo -time 22:00-02:00
```

#### 2. 优雅退出增强
- ✅ 第一次 Ctrl+C：等待当前任务完成后退出
- ✅ 第二次 Ctrl+C：立即强制退出
- ✅ 退出时保存最终统计数据

**工作流程：**
1. 按下第一次 Ctrl+C → 停止接收新任务，等待当前任务完成
2. 如需立即退出，按下第二次 Ctrl+C → 强制中断并退出

#### 3. 新模块
- ✅ 新增 `pkg/timerange` 模块，处理时间段解析和判断
- ✅ 包含完整的单元测试

### 🔧 改进

1. **下载逻辑优化**
   - 移除了下载任务中的强制中断
   - 改为等待任务自然完成
   - 提高数据完整性

2. **命令行参数优化**
   - 所有参数都支持简写形式
   - 更友好的错误提示

3. **文档完善**
   - 更新 README.md
   - 新增 EXAMPLES.md 使用示例
   - 新增 CHANGELOG.md 更新日志

### 📝 技术实现

**时间段控制架构：**
```
TimeRangeManager
├── 解析时间段字符串
├── 检查当前是否在时间段内
├── 计算到下一个时间段的等待时间
└── 支持多时间段和跨天场景

Downloader
├── Start() - 主循环
│   ├── 检查时间段
│   ├── 等待到时间段开始
│   └── 运行下载会话
└── runDownloadSession() - 单次会话
    ├── 启动协程池
    ├── 分发任务
    ├── 每秒检查时间段
    └── 时间段结束时停止
```

### 🧪 测试

新增单元测试覆盖：
- ✅ 时间段解析功能
- ✅ 时间段判断逻辑
- ✅ 跨天场景处理
- ✅ 错误处理

**测试结果：**
```
=== RUN   TestParseTimeRanges
--- PASS: TestParseTimeRanges (0.00s)
=== RUN   TestNewTimeRangeManager
--- PASS: TestNewTimeRangeManager (0.00s)
=== RUN   TestTimeRangeManager_IsInRange
--- PASS: TestTimeRangeManager_IsInRange (0.00s)
PASS
```

### 📊 性能

- CPU 使用：无明显增加
- 内存使用：无明显增加
- 时间检查：每秒一次，对性能影响可忽略

### 🐛 修复

- 修复了退出时任务被强制中断的问题
- 修复了速度文件可能不完整的问题

### 📦 依赖变化

- 无新增外部依赖
- 仅使用 Go 标准库

---

## [v1.0.0] - 2025-10-23

### 初始版本

- ✅ 多协程并发下载
- ✅ 支持指定 IP 地址
- ✅ 实时速度统计
- ✅ 循环下载模式
- ✅ 从 API 或文件加载任务
- ✅ 不保存文件到硬盘

