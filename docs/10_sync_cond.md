# sync.Cond, the Most Overlooked Sync Mechanism

`sync.Cond` 是 Go 语言中用于 goroutine 间条件变量同步的原语。它允许一组 goroutine 在某个条件满足之前进行等待，并在条件变化时被唤醒。尽管功能强大，但 `sync.Cond` 在实际开发中使用较少，主要是因为它的使用模式相对复杂，且大多数场景可以通过 channel 或 WaitGroup 解决。

## What is sync.Cond?（什么是 sync.Cond？）

`sync.Cond` 是基于条件变量的同步机制，包含三个核心方法：

- **Wait()**: 阻塞当前 goroutine，直到被 Signal 或 Broadcast 唤醒
- **Signal()**: 唤醒一个等待的 goroutine（如果有）
- **Broadcast()**: 唤醒所有等待的 goroutine

每个 `sync.Cond` 都与一个 `sync.Locker`（通常是 `sync.Mutex` 或 `sync.RWMutex`）关联，用于保护共享条件。

### 基本结构

```go
type Cond struct {
    noCopy noCopy
    L      Locker
    notify notifyList
    checker copyChecker
}
```

## 为什么需要 sync.Cond？

考虑以下场景：一个 goroutine 需要等待某个条件成立（如共享数据达到特定状态）。简单的轮询（循环检查）会浪费 CPU 资源：

```go
for !condition() {
    time.Sleep(100 * time.Millisecond)  // 浪费资源
}
```

而 `sync.Cond` 提供了高效的等待-通知机制，避免忙等待。

## 基本用法示例

### 生产者-消费者模式

以下示例展示了如何使用 `sync.Cond` 实现生产者-消费者模式，其中消费者等待特定条件（抓到 Pikachu）：

```go
package main

import (
    "math/rand"
    "sync"
    "time"
    "log"
)

var pokemonList = []string{"Pikachu", "Charmander", "Squirtle", "Bulbasaur", "Jigglypuff"}
var cond = sync.NewCond(&sync.Mutex{})
var pokemon = ""

func main() {
    now := time.Now().UnixMilli()
    // Consumer
    go func() {
        cond.L.Lock()
        defer cond.L.Unlock()

        // waits until Pikachu appears
        for pokemon != "Pikachu" {
            cond.Wait()
        }
        log.Printf("Caught %s, takes %d ms", pokemon, time.Now().UnixMilli()-now)
        pokemon = ""
    }()

    // Producer
    go func() {
        // Every 1ms, a random Pokémon appears
        for i := 0; i < 100; i++ {
            time.Sleep(time.Millisecond)

            cond.L.Lock()
            pokemon = pokemonList[rand.Intn(len(pokemonList))]
            cond.L.Unlock()

            cond.Signal()
        }
    }()

    time.Sleep(100 * time.Millisecond) // lazy wait
}
```

### 关键点

1. **条件检查必须在循环中**：因为 `Wait()` 返回时条件可能仍未满足（虚假唤醒或条件被其他 goroutine 改变）
2. **必须在持有锁的情况下调用 `Wait()`**：`Wait()` 会自动释放锁并阻塞，被唤醒后会重新获取锁
3. **Signal() vs Broadcast()**：`Signal()` 只唤醒一个等待者，`Broadcast()` 唤醒所有等待者

## 内部实现

### Wait() 方法的工作原理

`Wait()` 方法的内部逻辑如下：

```go
func (c *Cond) Wait() {
    // 1. 检查 Cond 是否被复制
    c.checker.check()
    
    // 2. 获取票据号
    t := runtime_notifyListAdd(&c.notify)
    
    // 3. 释放锁
    c.L.Unlock()
    
    // 4. 挂起 goroutine 直到被唤醒
    runtime_notifyListWait(&c.notify, t)
    
    // 5. 重新获取锁
    c.L.Lock()
}
```

**关键行为**：
- `Wait()` 原子性地释放锁并挂起当前 goroutine
- 被唤醒后，`Wait()` 在返回前会重新获取锁
- 这意味着其他 goroutine 可以在第一个 goroutine 等待期间获取相同的锁

### Signal() 和 Broadcast() 的实现

```go
func (c *Cond) Signal() {
    c.checker.check()
    runtime_notifyListNotifyOne(&c.notify)
}

func (c *Cond) Broadcast() {
    c.checker.check()
    runtime_notifyListNotifyAll(&c.notify)
}
```

### notifyList 结构

`notifyList` 是运行时内部结构，用于管理等待的 goroutine：

```go
type notifyList struct {
    wait   uint32      // 等待者计数
    notify uint32      // 已通知的票据号
    lock   uintptr     // 内部锁
    head   *sudog      // 等待队列头
    tail   *sudog      // 等待队列尾
}
```

## 使用模式

### 1. 典型条件等待模式

```go
cond.L.Lock()
for !condition() {
    cond.Wait()
}
// ... 使用条件 ...
cond.L.Unlock()
```

### 2. 多个等待者

```go
cond := sync.NewCond(&sync.Mutex{})
for i := range 10 {
    go func(i int) {
        cond.L.Lock()
        defer cond.L.Unlock()
        cond.Wait()
        
        fmt.Println(i)
    }(i)
}

time.Sleep(100 * time.Millisecond) // wait for goroutines to be ready
cond.Signal()  // 只唤醒一个
// 或 cond.Broadcast() 唤醒所有
```

## 常见陷阱

### 1. 忘记条件检查循环

```go
// 错误：可能虚假唤醒
cond.L.Lock()
cond.Wait()  // 可能在其他条件未满足时被唤醒
useCondition()
cond.L.Unlock()

// 正确：始终在循环中检查条件
cond.L.Lock()
for !condition() {
    cond.Wait()
}
useCondition()
cond.L.Unlock()
```

### 2. 复制 Cond

`sync.Cond` 包含 `noCopy` 字段，复制会导致运行时 panic：

```go
var c1 sync.Cond
c2 := c1  // 错误！panic: sync.Cond is copied
```

### 3. 锁管理不当

必须在调用 `Wait()` 前持有锁，`Wait()` 返回后仍然持有锁：

```go
cond.L.Lock()
// 必须在锁保护下检查条件
for !condition() {
    cond.Wait()  // 自动释放锁，唤醒后重新获取
}
// 仍然持有锁
cond.L.Unlock()
```

### 4. Signal/Broadcast 时机

通常应在持有锁的情况下发送信号，以确保条件变化的原子性：

```go
cond.L.Lock()
// 修改条件
condition = true
cond.L.Unlock()
cond.Signal()  // 可以在锁外调用，但需要确保条件变化对其他 goroutine 可见
```

## sync.Cond vs Channel

| 特性 | sync.Cond | Channel |
|------|-----------|---------|
| 等待多个 goroutine | 支持（Broadcast） | 需要多个 channel 或广播模式 |
| 条件检查 | 内置条件检查循环 | 需要外部逻辑 |
| 性能 | 更低开销，直接调度 | 更高的内存分配和调度开销 |
| 使用复杂度 | 较高，需要管理锁和条件 | 较低，直观 |
| 适用场景 | 高性能、复杂条件同步 | 简单同步、数据传递 |

**推荐**：除非需要高性能或复杂条件同步，否则优先使用 channel。

## 实际应用场景

### 1. 资源池（连接池、对象池）

当资源池为空时，请求者等待；当资源返回池中时，唤醒等待者。

### 2. 工作队列

多个工作者等待任务，生产者添加任务后唤醒工作者。

### 3. 状态同步

多个 goroutine 等待系统达到特定状态（如初始化完成、配置加载等）。

### 4. 限流器

当达到速率限制时，请求等待直到下一个时间窗口。

## 性能考虑

1. **锁竞争**：大量 goroutine 等待同一条件时，Signal/Broadcast 可能引起锁竞争
2. **虚假唤醒**：`Wait()` 可能在没有 Signal/Broadcast 的情况下返回，因此必须循环检查条件
3. **Broadcast 开销**：唤醒所有等待者可能引起 "惊群效应"，导致大量 goroutine 竞争锁

## 总结

1. **核心功能**：`sync.Cond` 提供了基于条件变量的高效等待-通知机制
2. **关键模式**：条件检查必须在循环中，必须在持有锁的情况下调用 `Wait()`
3. **内部实现**：使用 `notifyList` 管理等待队列，`Wait()` 自动释放和重新获取锁
4. **适用场景**：高性能同步、复杂条件等待、资源池等场景
5. **替代方案**：简单场景优先考虑 channel，`sync.Cond` 是性能关键场景的高级工具

`sync.Cond` 是 Go 并发工具箱中的强大但常被忽视的工具。正确使用时，它可以解决 channel 难以处理的高性能同步问题。然而，其复杂性也意味着更高的出错风险，因此建议仅在确实需要时使用。

## Navigation

← [Previous: sync.WaitGroup](./09_sync_waitgroup.md) | [Back to README](../README.md) | **Next:** [References](./11_references.md) →