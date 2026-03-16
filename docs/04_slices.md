# 切片（Slices）

## 切片概述

切片是 Go 语言中比数组更灵活的动态数据结构。与数组不同，切片的长度和容量可以在运行时改变。切片本质上是底层数组的一个视图（view），由三部分组成：指向底层数组的指针、长度（len）和容量（cap）。

### 创建方式

```go
// 1. nil 切片
var a []byte

// 2. 切片字面量
b := []byte{1, 2, 3}

// 3. 从数组创建
array := [6]byte{0, 1, 2, 3, 4, 5}
c := array[1:3]  // [1 2]

// 4. 使用 make
d := make([]byte, 1, 3)  // 长度=1，容量=3

// 5. 使用 new（不常见）
e := *new([]byte)
```

## 切片结构

### 内部表示

切片在运行时是一个包含三个字段的结构体：

```go
type slice struct {
    array unsafe.Pointer  // 指向底层数组的指针
    len   int             // 切片长度
    cap   int             // 切片容量
}
```

### 查看切片起始位置

有三种方法可以查看切片指向的底层数组起始位置：

#### 方法一：使用 `println`

```go
array := [6]byte{0, 1, 2, 3, 4, 5}
slice := array[1:3]
println("array:", &array)
println("slice:", slice)  // 输出: [2/5]0x1400004e6f3
// 切片指向 array[1]，地址为 0x1400004e6f3
```

#### 方法二：使用 `unsafe.SliceData`

```go
array := [6]byte{0, 1, 2, 3, 4, 5}
slice := array[1:3]
arrPtr := unsafe.SliceData(slice)
println("array[1]:", &array[1])      // 0x1400004e6f3
println("slice.array:", arrPtr)      // 0x1400004e6f3
```

**注意**：对于空切片（长度和容量为0但非nil），`unsafe.SliceData`返回一个指向特殊变量`zerobase`的指针。`zerobase`是所有零大小类型（如`struct{}`、`[0]int`）共享的内存地址。

#### 方法三：使用 sliceHeader 结构

```go
type sliceHeader struct {
    array unsafe.Pointer
    len   int
    cap   int
}

array := [6]byte{0, 1, 2, 3, 4, 5}
slice := array[1:3]
header := (*sliceHeader)(unsafe.Pointer(&slice))
println("sliceHeader:", header.array, header.len, header.cap)
// 输出: 0x1400004e6f3 2 5
```

**注意**：`reflect` 包中的 `reflect.SliceHeader` 类型已弃用，不应在生产代码中使用。上面的 `sliceHeader` 结构仅用于演示切片内部结构。

### 长度与容量

- **长度（len）**：切片中当前包含的元素数量
- **容量（cap）**：切片在不重新分配内存的情况下可以容纳的最大元素数量。对于从数组创建的切片，通常是从切片起始位置到底层数组末尾的元素数量。

```go
array := [6]int{0, 1, 2, 3, 4, 5}
slice := array[1:3]
fmt.Println(slice, len(slice), cap(slice))  // [1 2] 2 5
```

容量可以通过三参数切片操作指定：

```go
array := [6]int{0, 1, 2, 3, 4, 5}
slice := array[1:3:4]  // 容量 = 4 - 1 = 3
fmt.Println(len(slice), cap(slice))  // 2 3
```

## 切片增长

### append 操作

使用 `append()` 函数可以向切片添加元素。如果添加的元素不超过容量，则直接修改底层数组；否则会创建新的底层数组。

```go
array := [6]int{0, 1, 2, 3, 4, 5}
slice := array[1:3:4]  // [1 2], len=2, cap=3

// 容量内添加元素
slice = append(slice, 6)  // 修改底层数组，array[3] = 6
fmt.Println(slice)        // [1 2 6]
fmt.Println(array)        // [0 1 2 6 4 5]

// 超出容量，创建新数组
slice = append(slice, 7)  // 创建新底层数组
fmt.Println(slice)        // [1 2 6 7]
fmt.Println(array)        // [0 1 2 6 4 5] (未改变)
```

### 常见陷阱：函数参数中的切片

```go
func changeSlice(slice []int) {
    slice[0] = 100           // 修改原切片
    slice = append(slice, 400, 500)  // 可能创建新底层数组
}

func main() {
    slice := []int{1, 2, 3}
    changeSlice(slice)
    fmt.Println(slice)  // [100 2 3]，不是 [100 2 3 400 500]
}
```

### 容量增长策略

当切片需要增长时，Go 会分配新的底层数组。增长策略如下：

1. **小切片（容量 < 256）**：容量翻倍
2. **大切片（容量 ≥ 256）**：使用公式 `新容量 = 旧容量 + (旧容量 + 3*256)/4`
   - 近似于 `1.25 * 旧容量 + 192`
   - 实际容量还会根据内存对齐和大小类进行调整

#### 容量增长示例表

| 容量 | []int8 | []int32 | []int64 |
| ---- | ------ | ------- | ------- |
| 0    | 0      | 0       | 0       |
| 1    |        |         | 1       |
| 2    |        | 2       | 2       |
| 4    |        | 4       | 4       |
| 8    | 8      | 8       | 8       |
| 16   | 16     | 16      | 16      |
| 32   | 32     | 32      | 32      |
| 64   | 64     | 64      | 64      |
| 128  | 128    | 128     | 128     |
| 256  | 256    | 256     | 256     |
| 512  | 512    | 512     | 512     |
| 848  |        |         | 848     |
| 864  |        | 864     |         |
| 896  | 896    |         |         |
| 1280 |        |         | 1280    |
| 1344 |        | 1344    |         |
| 1408 | 1408   |         |         |
| 1792 |        |         | 1792    |
| 2048 | 2048   | 2048    |         |
| 2560 |        |         | 2560    |
| 3072 | 3072   | 3072    | 3072    |
| 3408 |        |         | 3408    |
| 4096 | 4096   | 4096    | 4096    |
| 5120 |        |         | 5120    |

### 扩展切片长度

如果不需要添加新元素，可以通过切片操作扩展长度：

```go
array := [6]int{0, 1, 2, 3, 4, 5}
slice := array[1:3]            // [1 2]
slice = slice[:len(slice)+1]   // [1 2 3]
```

**注意**：新长度不能超过切片容量。

## 内存分配

切片的内存分配涉及两部分：

1. **切片头（slice header）**：通常分配在栈上
2. **底层数组**：可能分配在栈或堆上，取决于大小和使用情况

### 栈分配情况

如果切片大小在编译时已知且较小，底层数组可能分配在栈上：

```go
func doSomething() {
    a := byte(1)
    println("a's address:", &a)           // 栈地址

    s := make([]byte, 1)
    println("slice's address:", &s)       // 栈地址
    println("underlying array:", s)       // 栈地址，接近 a 的地址
}
```

### 堆分配情况

以下情况底层数组会分配在堆上：

1. **容量超过 64KB**

```go
sliceA := make([]byte, 64*1024)      // 栈分配（刚好 64KB）
sliceB := make([]byte, 64*1024+1)    // 堆分配
```

2. **动态大小**

```go
func arrayOnHeap(n int) {
    slice := make([]int, n)  // 底层数组在堆上分配
}
```

3. **切片增长导致重新分配**

```go
slice := make([]int, 0, 3)
slice = append(slice, 1, 2, 3)    // 栈分配
slice = append(slice, 4)          // 堆分配（容量增长）
```

### 优化建议

1. **预分配容量**：使用 `make([]T, 0, estimatedCapacity)` 减少重新分配
2. **重用底层数组**：使用 `sync.Pool` 复用切片

   ```go
   var pool = sync.Pool{
       New: func() interface{} {
           return make([]byte, 0, 1024)
       },
   }

   // 获取切片
   slice := pool.Get().([]byte)
   slice = slice[:0]  // 重置长度

   // 使用后放回
   slice = append(slice, data...)
   pool.Put(slice)
   ```

## 重要注意事项

1. **切片共享底层数组**：多个切片可能共享同一个底层数组，修改一个切片会影响其他切片
2. **append 可能断开连接**：当切片容量不足时，`append()` 会创建新底层数组，断开与原数组的连接
3. **函数参数传递**：切片作为函数参数传递时，传递的是切片头的拷贝，但底层数组是共享的
4. **空切片 vs nil 切片**
   - **nil 切片**：切片头中的指针为nil，长度和容量都为0
   - **空切片**：切片头中的指针指向`zerobase`（一个特殊的零大小内存地址），长度和容量为0

   ```go
   var nilSlice []int        // nil 切片，pointer=nil, len=0, cap=0
   emptySlice := []int{}     // 空切片，pointer=zerobase, len=0, cap=0
   makeSlice := make([]int, 0)  // 空切片，pointer=zerobase, len=0, cap=0
   ```

   **zerobase 特性**：
   - 所有零大小类型（`struct{}`、`[0]int`）都共享同一个`zerobase`地址
   - `unsafe.SliceData([]int{})`返回`zerobase`地址
   - 空切片的底层数组不占用额外内存

## 性能建议

1. **避免频繁重新分配**：合理预估容量，使用 `make` 预分配
2. **大切片考虑堆分配影响**：超过 64KB 的切片会在堆上分配
3. **热点路径优化**：在性能关键路径上，考虑复用切片或使用对象池
4. **注意内存泄漏**：大切片即使长度缩小，底层数组仍然占用内存

## Navigation

← [Previous: Maps](./03_maps.md) | [Back to README](../README.md) | **Next:** [Goroutines and Channels](./05_goroutines-channels.md) →
