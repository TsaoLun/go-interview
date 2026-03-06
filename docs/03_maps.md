# Maps (映射/字典)

A Go map is a hash table composed of many **buckets**. (Go的map是由许多**桶**组成的哈希表)

```go
type hmap struct {
    // ...
    buckets unsafe.Pointer // pointer to the bucket array (指向桶数组的指针)
    // ...
}
```

When you assign a map to a variable or pass it to a function, both share the same map pointer. However, maps are **not reference types**; they are pointers to the underlying `hmap` structure. (当你将map赋值给变量或传递给函数时，两者共享相同的map指针。然而，map**不是引用类型**；它们是指向底层`hmap`结构的指针)

> See: [There is no pass-by-reference in Go](https://dave.cheney.net/2017/04/29/there-is-no-pass-by-reference-in-go) (参见：[Go中没有按引用传递](https://dave.cheney.net/2017/04/29/there-is-no-pass-by-reference-in-go))

**Example (示例):**
```go
package main

import "fmt"

func fn(m map[int]int) {
    m = make(map[int]int) // creates a new map; does NOT affect the caller's map (创建新map；不会影响调用者的map)
}

func main() {
    var m map[int]int
    fn(m)
    fmt.Println(m == nil) // true
}
```

**Buckets (bucket array) (桶，桶数组)**

1. Hash the key (e.g., `hash("hello", seed)`) to a number. (将键哈希为数字，例如：`hash("hello", seed)`)
2. Modulo the hash by the number of buckets to locate the target bucket. (将哈希值与桶数量取模以定位目标桶)

**Growth strategies (扩容策略):**

- **Double the bucket count** – when load factor exceeds ~80% (i.e., too many entries). (**双倍增加桶数量** - 当负载因子超过约80%时，即条目过多)
- **Redistribute entries** – when there are too many overflow buckets. (**重新分配条目** - 当有太多溢出桶时)

## Navigation

← [Previous: Mutex and Race Conditions](./02_mutex.md) | [Back to README](../README.md) | **Next:** [Slices](./04_slices.md) →