# Interfaces (接口)

An interface defines a set of method signatures. A type implements an interface implicitly by implementing all its methods. (接口定义了一组方法签名。类型通过实现其所有方法隐式地实现接口)

**Empty interface (`interface{}` / `any`):** holds any type. (**空接口 (`interface{}` / `any`):** 可容纳任何类型)

```go
type Writer interface {
    Write([]byte) (int, error)
}

type MyWriter struct{}
func (mw MyWriter) Write(data []byte) (int, error) {
    // ...
    return len(data), nil
}

var w Writer = MyWriter{}
```

**Interface values** consist of a **type** and a **value** (the concrete data). Use type assertions or type switches to retrieve the underlying value. (**接口值**包含一个**类型**和一个**值**(具体数据)。使用类型断言或类型切换来检索底层值)

## Navigation

← [Previous: Goroutines and Channels](./05_goroutines-channels.md) | [Back to README](../README.md) | **Next:** [Defer, Panic, and Recover](./07_defer-panic-recover.md) →