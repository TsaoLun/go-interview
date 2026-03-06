# Mutex and Race Conditions

## Race Condition (竞争条件)

Occurs when multiple goroutines access and modify shared data concurrently without proper synchronization. (多个goroutine在没有适当同步的情况下并发访问和修改共享数据)

```go
var counter = 0

func incrementCounter() {
    counter++
}

func main() {
    for range 1000 {
        go incrementCounter()
    }

    time.Sleep(time.Second)
    fmt.Println(counter) // Output: 9xx (varies)
}
```

Use `sync.Mutex` (or `sync.RWMutex`) to synchronize access. (使用 `sync.Mutex` 或 `sync.RWMutex` 来同步访问)

## `sync.Mutex` Internals (`sync.Mutex` 内部结构)

```go
package sync

type Mutex struct {
    state int32
    sema  uint32 // semaphore identifier; maps to a global hash table of queues (信号量标识符；映射到全局哈希表中的队列)
}
```

**`state` bits (int32) (`state` 位，int32):**
- **bit 0 (Locked)**: 1 = mutex is locked (1 = 互斥锁已锁定)
- **bit 1 (Woken)**: 1 = a goroutine has been woken and is trying to acquire the mutex (1 = 有goroutine被唤醒并尝试获取互斥锁)
- **bit 2 (Starving)**: 1 = starvation mode is active (1 = 饥饿模式激活)
- **bits 3–31 (Waiter)**: number of goroutines waiting to acquire the mutex (等待获取互斥锁的goroutine数量)

### `Lock()` method (`Lock()` 方法)

```go
func (m *Mutex) Lock() {
    // Fast path: grab unlocked mutex. (快速路径：获取未锁定的互斥锁)
    if atomic.CompareAndSwapInt32(&m.state, 0, 1) {
        if race.Enabled {
            race.Acquire(unsafe.Pointer(m))
        }
        return
    }
    // Slow path (spinning to avoid sleep‑wake overhead) (慢路径：自旋以避免睡眠-唤醒开销)
    m.lockSlow()
}
```

## Starving Mode (饥饿模式)

- Spinning does not work in starving mode. (在饥饿模式下自旋无效)
- In normal mode, new goroutines may race and acquire the mutex before a queued goroutine wakes up. (正常模式下，新的goroutine可能在队列中的goroutine唤醒前竞争并获取互斥锁)
- Starvation mode activates when a goroutine fails to acquire the lock for more than 1 ms. (当goroutine在超过1毫秒内无法获取锁时，饥饿模式激活)
- New goroutines do not try to acquire the mutex; they join the end of the waiting queue. (新的goroutine不尝试获取互斥锁；它们加入等待队列的末尾)
- Starvation mode persists until the waiting queue is empty or the wait time drops below 1 ms. (饥饿模式持续直到等待队列为空或等待时间低于1毫秒)

## Navigation

← [Previous: Array](./01_array.md) | [Back to README](../README.md) | **Next:** [Maps](./03_maps.md) →