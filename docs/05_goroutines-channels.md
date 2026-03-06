# Goroutines and Channels (Goroutines 和通道)

## Goroutines (Goroutines)

Lightweight threads managed by the Go runtime. (由Go运行时管理的轻量级线程)

```go
go func() {
    fmt.Println("Running in a goroutine")
}()
```

## Channels (通道)

Provide communication and synchronization between goroutines. (提供goroutine之间的通信和同步)

- **Unbuffered channel**: `make(chan T)` – synchronous; sender blocks until receiver is ready. (**无缓冲通道**: `make(chan T)` - 同步的；发送者阻塞直到接收者就绪)
- **Buffered channel**: `make(chan T, n)` – asynchronous up to the buffer size. (**缓冲通道**: `make(chan T, n)` - 异步的，最多到缓冲区大小)

```go
ch := make(chan int, 2)
ch <- 1
ch <- 2
fmt.Println(<-ch) // 1
fmt.Println(<-ch) // 2
```

**Select statement (Select语句):**
```go
select {
case msg := <-ch1:
    fmt.Println("Received from ch1:", msg)
case ch2 <- 42:
    fmt.Println("Sent 42 to ch2")
default:
    fmt.Println("No communication ready")
}
```

## Navigation

← [Previous: Slices](./04_slices.md) | [Back to README](../README.md) | **Next:** [Interfaces](./06_interfaces.md) →