# Worker Pool 设计文档

## 概述

实现了一个高性能、线程安全的 Worker Pool，用于并发处理事件。该设计遵循 Go 并发模式最佳实践，支持优雅关闭、错误处理和灵活配置。

## 核心特性

### 1. **并发处理**
- 支持配置多个 worker goroutine 并发处理事件
- 通过 buffered channel 实现事件队列，减少阻塞
- 自动负载均衡，多个 worker 从同一队列消费

### 2. **生命周期管理**
```go
pool := NewWorkerPool(config)
pool.Start()              // 启动所有 worker
defer pool.Stop()         // 优雅关闭，等待所有任务完成
// 或使用超时关闭
pool.StopWithTimeout(5 * time.Second)
```

### 3. **错误处理**
- 独立的错误 channel，不阻塞业务处理
- 每个错误包含 worker ID 和详细信息
- 支持异步监听错误

### 4. **提交策略**
```go
// 非阻塞提交（队列满时立即返回错误）
err := pool.Submit(event)

// 带超时的提交
err := pool.SubmitWithTimeout(event, 100*time.Millisecond)
```

### 5. **线程安全**
- 使用 `sync.RWMutex` 保护状态
- 使用 `sync.WaitGroup` 等待 worker 退出
- 使用 `sync.Once` 确保只关闭一次

## 架构设计

### 组件关系图

```
                        ┌─────────────────┐
                        │  WorkerPool     │
                        ├─────────────────┤
                        │ eventChan       │◄─── Submit()
                        │ processor       │
                        │ errorChan       │───► Errors()
                        │ ctx/cancel      │
                        └────────┬────────┘
                                 │
                    ┌────────────┼────────────┐
                    ▼            ▼            ▼
               ┌─────────┐ ┌─────────┐ ┌─────────┐
               │Worker 1 │ │Worker 2 │ │Worker N │
               └─────────┘ └─────────┘ └─────────┘
                    │            │            │
                    └────────────┼────────────┘
                                 ▼
                          Process Event
                                 │
                        ┌────────┴────────┐
                        ▼                 ▼
                    Success          errorChan
```

### 数据流

1. **事件提交**: `Submit()` → `eventChan` (buffered)
2. **事件处理**: Worker 从 `eventChan` 取事件 → 执行 `processor`
3. **错误处理**: 处理失败 → `errorChan` (non-blocking)
4. **优雅关闭**: `Stop()` → close(`eventChan`) → workers 处理完剩余事件 → `wg.Wait()`

## 使用示例

### 基础用法

```go
// 定义事件处理器
processor := func(event Event) error {
    fmt.Printf("Processing: %s\n", event.Type)
    // 业务逻辑
    return nil
}

// 配置 worker pool
config := WorkerPoolConfig{
    WorkerNum:     5,      // 5个并发worker
    QueueSize:     100,    // 队列容量100
    ErrorChanSize: 10,     // 错误缓冲10
    Processor:     processor,
}

pool := NewWorkerPool(config)
pool.Start()
defer pool.Stop()

// 提交事件
event := Event{
    Type:    "user.created",
    From:    "api",
    Payload: userData,
}
if err := pool.Submit(event); err != nil {
    log.Printf("Submit failed: %v", err)
}
```

### 错误监听

```go
// 启动错误监听 goroutine
go func() {
    for err := range pool.Errors() {
        log.Printf("Worker error: %v", err)
        // 可以进行告警、重试等操作
    }
}()
```

### 优雅关闭

```go
// 方式1: 等待所有任务完成
pool.Stop()

// 方式2: 带超时的关闭
if err := pool.StopWithTimeout(5 * time.Second); err != nil {
    log.Printf("Forced shutdown: %v", err)
}
```

## 配置说明

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `WorkerNum` | int | 5 | 并发 worker 数量 |
| `QueueSize` | int | 100 | 事件队列缓冲大小 |
| `ErrorChanSize` | int | 10 | 错误通道缓冲大小 |
| `Processor` | Processor | defaultProcessor | 事件处理函数 |

### 性能调优建议

- **CPU 密集型**: `WorkerNum = NumCPU`
- **IO 密集型**: `WorkerNum = 2 * NumCPU` 或更高
- **QueueSize**: 根据事件产生速率设置，避免频繁阻塞
- **ErrorChanSize**: 通常设置为 `QueueSize / 10`

## 设计优势

### ✅ 相比原设计的改进

| 方面 | 原设计 | 新设计 |
|------|--------|--------|
| **Processor 管理** | ❌ 通过 channel 传递，会被消费 | ✅ 直接字段，可复用 |
| **并发能力** | ❌ 单 goroutine | ✅ 多 worker 并发 |
| **资源管理** | ❌ 无优雅关闭 | ✅ WaitGroup + Context |
| **错误处理** | ❌ 空实现 | ✅ 独立错误通道 |
| **API 完整性** | ❌ 缺少提交接口 | ✅ Submit/SubmitWithTimeout |
| **状态管理** | ❌ 无状态保护 | ✅ RWMutex + Once |
| **监控能力** | ❌ 无 | ✅ QueueLen/IsStarted/Errors |

### 🚀 核心优势

1. **高性能**: 多 worker 并行处理，充分利用多核
2. **可靠性**: 优雅关闭，确保事件不丢失
3. **可观测**: 提供队列长度、错误监听等监控指标
4. **易用性**: 简洁的 API，灵活的配置
5. **安全性**: 线程安全，防止资源泄漏

## 测试覆盖

- ✅ 基本功能测试
- ✅ 错误处理测试
- ✅ 并发提交测试
- ✅ 超时机制测试
- ✅ 优雅关闭测试
- ✅ 边界条件测试
- ✅ 性能基准测试

运行测试：
```bash
go test ./internal/event/... -v
go test ./internal/event/... -bench=. -benchmem
```

## 最佳实践

### 1. 合理设置 Worker 数量

```go
import "runtime"

config := DefaultConfig()
config.WorkerNum = runtime.NumCPU() // CPU密集型
// 或
config.WorkerNum = runtime.NumCPU() * 2 // IO密集型
```

### 2. 监控队列长度

```go
ticker := time.NewTicker(10 * time.Second)
go func() {
    for range ticker.C {
        qLen := pool.QueueLen()
        if qLen > config.QueueSize * 0.8 {
            log.Warn("Queue is nearly full: %d/%d", qLen, config.QueueSize)
        }
    }
}()
```

### 3. 实现重试机制

```go
processor := func(event Event) error {
    const maxRetries = 3
    var err error
    
    for i := 0; i < maxRetries; i++ {
        if err = processEvent(event); err == nil {
            return nil
        }
        time.Sleep(time.Second * time.Duration(i+1))
    }
    
    return fmt.Errorf("failed after %d retries: %w", maxRetries, err)
}
```

### 4. 上下文传递

```go
type Event struct {
    Type     string
    Payload  any
    From     string
    SendAt   time.Time
    HandleAt time.Time
    Context  context.Context // 添加上下文
}

processor := func(event Event) error {
    ctx := event.Context
    if ctx == nil {
        ctx = context.Background()
    }
    
    // 使用 context 进行超时控制、取消等
    return doWorkWithContext(ctx, event)
}
```

## 未来改进方向

- [ ] 添加 Metrics (处理延迟、吞吐量等)
- [ ] 支持优先级队列
- [ ] 支持动态调整 worker 数量
- [ ] 支持事件重试队列（DLQ）
- [ ] 添加分布式追踪支持

## 许可

MIT License
