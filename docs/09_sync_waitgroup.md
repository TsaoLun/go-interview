# sync.WaitGroup and The Alignment Problem

`sync.WaitGroup` 是 Go 语言中用于等待一组 goroutine 完成的标准同步原语。当我们需要并发执行多个独立任务，并等待所有任务完成后才能继续执行主 goroutine 时，WaitGroup 提供了简洁的解决方案。

## What is sync.WaitGroup?（什么是 sync.WaitGroup？）

WaitGroup 的基本工作模式：

- **Add(delta int)**: 增加待完成的 goroutine 数量（正数 delta），或减少已完成的 goroutine 数量（负数 delta）
- **Done()**: 减少计数器的值，相当于 `Add(-1)`
- **Wait()**: 阻塞当前 goroutine，直到计数器归零

### 基本用法示例

```go
func main() {
    var wg sync.WaitGroup

    wg.Add(10)
    for i := 0; i < 10; i++ {
        go func(i int) {
            defer wg.Done()
            fmt.Println("Task", i)
        }(i)
    }

    wg.Wait()
    fmt.Println("Done")
}
```

### 常见使用模式

1. **使用 `wg.Add(1)` 避免错误**（推荐）：

```go
for i := 0; i < 10; i++ {
    wg.Add(1)  // 在启动 goroutine 前增加计数
    go func() {
        defer wg.Done()
        // ...
    }()
}
```

2. **避免的错误用法**：

```go
// 错误：Add 在 goroutine 内部调用
for i := 0; i < 10; i++ {
    go func() {
        wg.Add(1)  // 可能先执行 wg.Wait()
        defer wg.Done()
        // ...
    }()
}
```

## 内部结构

### WaitGroup 结构定义（Go 1.23）

```go
type WaitGroup struct {
    noCopy noCopy

    state atomic.Uint64
    sema  uint32
}
```

### 1. noCopy 机制

`noCopy` 是一个空结构体，用于防止 WaitGroup 被意外复制：

```go
type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}
```

`go vet` 工具会检测包含 `noCopy` 字段的结构体是否被复制，并在发现复制时发出警告：

```go
func main() {
    var a sync.WaitGroup
    b := a  // go vet: assignment copies lock value to b

    fmt.Println(a, b)  // go vet: call of fmt.Println copies lock value
}
```

### 2. 内部状态（state）

`state` 是一个 `atomic.Uint64` 类型，其中：

- **高 32 位**：计数器（counter），跟踪待完成的 goroutine 数量
- **低 32 位**：等待者计数（waiter），跟踪调用 `Wait()` 的 goroutine 数量

```
+----------------------+----------------------+
|      Counter         |       Waiter         |
|    (高 32 位)         |    (低 32 位)         |
+----------------------+----------------------+
```

### 3. 信号量（sema）

`sema` 是一个内部信号量，用于阻塞和唤醒等待的 goroutine：

- 当 goroutine 调用 `wg.Wait()` 且计数器不为零时，会调用 `runtime_Semacquire(&wg.sema)` 阻塞
- 当计数器归零时，会调用 `runtime_Semrelease(&wg.sema)` 唤醒所有等待的 goroutine

## 对齐问题（The Alignment Problem）

### 问题背景

在 32 位架构（ARM、386、32-bit MIPS）上，编译器不保证 64 位值在 8 字节边界对齐，可能只在 4 字节边界对齐。然而，`atomic` 包要求 64 位原子操作的变量必须 8 字节对齐，否则可能导致程序崩溃。

```go
// atomic 包文档说明：
// On ARM, 386, and 32-bit MIPS, it is the caller's responsibility
// to arrange for 64-bit alignment of 64-bit words accessed atomically.
```

### 历史解决方案演变

#### Go 1.5: state1 [12]byte

```go
type WaitGroup struct {
    state1 [12]byte  // 12 字节数组，确保有足够的空间找到 8 字节对齐的区域
    sema   uint32
}

func (wg *WaitGroup) state() *uint64 {
    if uintptr(unsafe.Pointer(&wg.state1))%8 == 0 {
        return (*uint64)(unsafe.Pointer(&wg.state1))
    } else {
        return (*uint64)(unsafe.Pointer(&wg.state1[4]))
    }
}
```

**原理**：使用 12 字节数组，如果起始地址不是 8 字节对齐，则偏移 4 字节找到对齐位置。

#### Go 1.11: state1 [3]uint32

```go
type WaitGroup struct {
    noCopy noCopy
    state1 [3]uint32  // 合并 state 和 sema 为 3 个 uint32
}

func (wg *WaitGroup) state() (statep *uint64, semap *uint32) {
    if uintptr(unsafe.Pointer(&wg.state1))%8 == 0 {
        return (*uint64)(unsafe.Pointer(&wg.state1)), &wg.state1[2]
    } else {
        return (*uint64)(unsafe.Pointer(&wg.state1[1])), &wg.state1[0]
    }
}
```

**原理**：根据对齐情况重新排列 3 个 uint32，确保 state 部分（counter + waiter）8 字节对齐。

#### Go 1.18: state1 uint64; state2 uint32

```go
type WaitGroup struct {
    noCopy noCopy
    state1 uint64
    state2 uint32
}

func (wg *WaitGroup) state() (statep *uint64, semap *uint32) {
    if unsafe.Alignof(wg.state1) == 8 || uintptr(unsafe.Pointer(&wg.state1))%8 == 0 {
        // 64 位对齐：直接使用
        return &wg.state1, &wg.state2
    } else {
        // 32 位系统：转换为 [3]uint32 并重新排列
        state := (*[3]uint32)(unsafe.Pointer(&wg.state1))
        return (*uint64)(unsafe.Pointer(&state[1])), &state[0]
    }
}
```

#### Go 1.20+: atomic.Uint64（当前方案）

```go
type Uint64 struct {
    _ noCopy
    _ align64  // 编译器识别的特殊标记，强制 8 字节对齐
    v uint64
}

type align64 struct{}  // 编译器特殊处理，不占用空间

type WaitGroup struct {
    noCopy noCopy
    state atomic.Uint64  // 自动 8 字节对齐
    sema  uint32
}
```

**align64 的魔力**：`align64` 是一个空结构体，但编译器会识别它并强制整个 `atomic.Uint64` 结构体在 8 字节边界对齐。这是编译器层面的特殊处理，不是运行时机制。

## 内部工作原理

### Add(delta int) 方法

```go
func (wg *WaitGroup) Add(delta int) {
    // 1. 原子地将 delta 添加到 counter（高 32 位）
    state := wg.state.Add(uint64(delta) << 32)

    // 2. 提取 counter 和 waiter
    v := int32(state >> 32)  // counter
    w := uint32(state)       // waiter

    // 3. 检查 counter 不能为负数
    if v < 0 {
        panic("sync: negative WaitGroup counter")
    }

    // 4. 如果 counter > 0 或没有等待者，直接返回
    if v > 0 || w == 0 {
        return
    }

    // 5. counter == 0 且有等待者，唤醒所有等待的 goroutine
    wg.state.Store(0)
    for ; w != 0; w-- {
        runtime_Semrelease(&wg.sema, false, 0)
    }
}
```

**关键点**：

1. `Done()` 就是 `Add(-1)`
2. 如果 delta 使 counter 变为负数，会 panic
3. 必须在所有 `wg.Add(positive)` 调用完成后才能调用 `wg.Wait()`
4. 可以在任何时候调用 `wg.Add(negative)`，只要 counter 不变成负数

### Wait() 方法

```go
func (wg *WaitGroup) Wait() {
    for {
        state := wg.state.Load()
        v := int32(state >> 32)  // counter
        w := uint32(state)       // waiter

        // counter 为 0，无需等待
        if v == 0 {
            return
        }

        // 尝试增加 waiter 计数
        if wg.state.CompareAndSwap(state, state+1) {
            // 成功增加 waiter，阻塞当前 goroutine
            runtime_Semacquire(&wg.sema)

            // 被唤醒后，检查 WaitGroup 是否被重用
            if wg.state.Load() != 0 {
                panic("sync: WaitGroup is reused before previous Wait has returned")
            }
            return
        }
        // CAS 失败，重试
    }
}
```

**关键点**：

1. 使用 CAS（Compare-And-Swap）原子操作增加 waiter 计数
2. 如果 CAS 失败（状态被其他 goroutine 修改），重试
3. 被唤醒后检查状态，防止 WaitGroup 被重用

## 常见陷阱和最佳实践

### 1. 不要复制 WaitGroup

```go
var wg1 sync.WaitGroup
wg2 := wg1  // 错误！go vet 会警告
```

### 2. Add 必须在 Wait 之前调用

```go
// 正确
wg.Add(1)
go func() { defer wg.Done(); /* ... */ }()
wg.Wait()

// 错误：Add 在 goroutine 内部
go func() {
    wg.Add(1)  // 可能先执行下面的 wg.Wait()
    defer wg.Done()
}()
wg.Wait()  // 可能立即返回
```

### 3. 使用 defer wg.Done()

```go
go func() {
    defer wg.Done()  // 确保 Done 一定会被调用
    // 可能 panic 或有多条 return 路径
    if err != nil {
        return
    }
    // ...
}()
```

### 4. 不要重用 WaitGroup

WaitGroup 设计为单次使用。如果必须重用，确保：

1. 所有 `Wait()` 调用都已完成
2. counter 已归零
3. 没有 goroutine 还在引用旧的 WaitGroup

### 5. 性能考虑

- 使用 `wg.Add(1)` 而不是 `wg.Add(n)` 更安全，但有小性能开销
- 在大多数 64 位系统上，atomic.Uint64 确保对齐，性能最佳

## 总结

1. **核心功能**：WaitGroup 用于等待一组 goroutine 完成
2. **对齐问题**：历史版本通过各种技巧（12 字节数组、3 uint32 数组等）解决 32 位系统的对齐问题
3. **当前方案**：使用 `atomic.Uint64` 和 `align64` 标记，编译器强制 8 字节对齐
4. **内部实现**：无锁算法，使用原子操作和信号量
5. **最佳实践**：不复制、Add 在 Wait 之前、使用 defer Done、不重用

WaitGroup 的设计展示了 Go 语言在并发原语实现上的精妙平衡：既保证了正确性，又通过编译器和运行时的协作解决了底层硬件对齐问题。

## Navigation

← [Previous: Memory Management and GC](./08_memory-gc.md) | [Back to README](../README.md) | **Next:** [sync.Cond](./10_sync_cond.md) →
