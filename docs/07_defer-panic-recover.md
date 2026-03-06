# Defer, Panic, and Recover (Defer、Panic 和 Recover)

- **`defer`** schedules a function call to run when the surrounding function returns (LIFO order). (**`defer`** 安排函数调用在周围函数返回时执行（LIFO顺序）)
- **`panic`** stops normal execution and begins panicking. (**`panic`** 停止正常执行并开始 panic)
- **`recover`** regains control of a panicking goroutine; only useful inside deferred functions. (**`recover`** 重新获得 panic goroutine 的控制；仅在defer函数中有用)

```go
func safeDivide(a, b int) (result int, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic recovered: %v", r) // 恐慌恢复：%v
        }
    }()
    return a / b, nil
}
```

## Navigation

← [Previous: Interfaces](./06_interfaces.md) | [Back to README](../README.md) | **Next:** [Memory Management and GC](./08_memory-gc.md) →