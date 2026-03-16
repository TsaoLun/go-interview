# Array

## 什么是数组？

数组是具有固定大小、存储相同类型元素的连续内存区域。数组的长度是其类型的一部分，因此 `[5]byte` 和 `[4]byte` 是不同的类型。

### 内存布局

```go
func main() {
    arr := [5]byte{0, 1, 2, 3, 4}
    println("arr", &arr)

    for i := range arr {
        println(i, &arr[i])
    }
}
// 输出：
// arr 0x1400005072b
// 0 0x1400005072b
// 1 0x1400005072c
// 2 0x1400005072d
// 3 0x1400005072e
// 4 0x1400005072f
```

注意：

- 数组的地址与第一个元素的地址相同。
- 每个元素的地址相隔 1 字节（因为元素类型是 `byte`）。

## 数组字面量

有多种初始化数组的方式：

```go
var arr1 [10]int                    // [0 0 0 0 0 0 0 0 0 0]
arr2 := [...]int{1, 2, 3, 4, 5}    // [1 2 3 4 5]
arr3 := [...]int{11: 3}            // [0 0 0 0 0 0 0 0 0 0 0 3]
arr4 := [5]int{1, 4: 5}            // [1 0 0 0 5]
arr5 := [5]int{2: 3, 4, 4: 5}      // [0 0 3 4 5]
```

### 初始化策略

- 元素数量 ≤ 4：采用 **局部代码初始化**（local-code initialization），编译器生成逐个赋值的指令。
- 元素数量 > 4：采用 **静态初始化**（static initialization），数组值被嵌入到二进制文件的只读段中。

### 栈与堆分配

默认情况下，数组在栈上分配。如果数组的大小超过 `MaxStackVarSize`（Go 1.25 之前为 10 MB，Go 1.25 改为 128 KB），则会逃逸到堆上。

```go
func main() {
    a := [10 * 1024 * 1024]byte{}   // 栈分配
    b := [10 * 1024 * 1024 + 1]byte{} // 堆分配
}
```

## 数组操作

### 长度与容量

数组的长度是类型的一部分，因此 `len(arr)` 和 `cap(arr)` 在编译时是常量。

```go
a := [5]int{1, 2, 3} // [1 2 3 0 0]
println(len(a)) // 5
println(cap(a)) // 5
```

### 切片（Slicing）

可以从数组创建切片，语法为 `[start:end:capacity]`，其中 `start` 包含，`end` 不包含，`capacity` 可选。

```go
a := [5]int{0, 1, 2, 3, 4}
b := a[1:3]   // [1 2]
c := a[:3]    // [0 1 2]
d := a[1:]    // [1 2 3 4]
```

切片规则：

- `start` 默认为 0
- `end` 默认为数组长度
- `capacity` 默认为数组长度 - `start`

边界检查：`0 <= start <= end <= cap <= 实际容量`。允许 `start` 等于数组长度，此时会生成空切片。

## 数组是值类型

Go 中的数组是值类型，赋值或传参时会复制整个数组。

```go
func doSomething(a [5]byte) {
    a[0] = 1
}

func main() {
    a := [5]byte{}
    doSomething(a)
    fmt.Println(a) // [0 0 0 0 0]
}
```

### for-range 的陷阱

使用 `for-range` 遍历数组时，Go 会创建数组的副本，迭代变量 `v` 访问的是副本。

```go
func main() {
    a := [3]int{1, 2, 3}
    b := [3]int{4, 5, 6}

    for i, v := range a {
        if i == 1 {
            a = b
        }
        fmt.Println(v)
    }
    // 输出：1 2 3（而不是 1 2 6）
}
```

如果希望避免复制，可以使用指针：

```go
for i, v := range &a {
    // 此时修改 a 会影响迭代，输出 1 2 6
}
```

## 逃逸分析

```bash
go build -gcflags="-m" .
```

- Go 1.25: `MaxStackVarSize = int64(128 * 1024)`（旧版本为 `10 * 1024 * 1024`）

## Navigation

← [Back to README](../README.md) | **Next:** [Mutex and Race Conditions](./02_mutex.md) →
