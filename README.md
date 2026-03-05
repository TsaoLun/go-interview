# Go Interview Brief

### Array

元素长度与类型固定，默认栈上分配，当大小超过 MaxStackVarSize 时会逃逸到堆上
> 可以通过指令分析 `go build -gcflags="-m" .` Go 1.25 版本 MaxStackVarSize = int64(128 * 1024) 旧版本为 10 * 1024 * 1024

### Mutex

#### race condition
multiple goroutines try to access and change shared data at the same time without proper synchronization.
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
	fmt.Println(counter) // Output: 9xx
}
```

To avoid race condition, we can use mutex to synchronize access to shared data.

```go
package sync

type Mutex struct {
	state int32
	sema  uint32 // 信号量标识符 key -> 全局哈希表对应队列
}
```

接下来分析 state 字段(int32)
* Locked(bit_0): If it’s set to 1, the mutex is locked and no other goroutine can grab it.
* Woken(bit_1): Set to 1 if any goroutine has been woken up and is trying to acquire the mutex.
* Starving(bit_2): The starvation mode.
* Waiter(bit_3-31): How many goroutines are waiting to acquire the mutex.

```go
func (m *Mutex) Lock() {
	// Fast path: grab unlocked mutex.
	// Compare And Swap operation fail means the state locked
	if atomic.CompareAndSwapInt32(&m.state, 0, 1) {
		if race.Enabled {
			race.Acquire(unsafe.Pointer(m))
		}
		return
	}
	// Slow path (spining to avoid the overhead of a sleep-wake cycle)
	// Spinning, repeatedly checking the mutex state and then go to sleep(sema->queue)
	m.lockSlow()
}
```
#### Starvation Mode

Spinning doesn’t work in Starvation mode.

In normal mode, new goroutines can quickly try to grab the mutex, while the queued goroutine is still waking up (lose the race to the new contenders and get put back at the front of the queue.)

Starvation mode kicks in if a goroutine fails to acquire the lock for more than 1 millisecond. New goroutines don’t even try to acquire mutex and just join the end of the waiting queue.

Starvation mode continues until the waiting queue is empty or the goroutine wait for less than one millisecond.




### Maps

Go map is composed of many smaller units called "buckets":

```go
type hmap struct {
  ...
  buckets unsafe.Pointer // point to the bucket array.
  ...
}
````

When you assign a map to a variable or pass it to a function, both the variable and the function’s argument are sharing the same map pointer. But maps are pointers to the hmap under the hood, they aren’t reference types.
> 见 https://dave.cheney.net/2017/04/29/there-is-no-pass-by-reference-in-go

#### Buckets(bucket array)

By hashing key ("hello" -> hash("hello", seed)) to a number, then it takes that number and mods it by the number of buckets.



### TODO LIST
