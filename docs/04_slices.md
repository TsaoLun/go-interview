# Slices (切片)

- A slice is a descriptor of a contiguous segment of an underlying array. (切片是底层数组连续段的描述符)
- Contains a pointer to the array, length, and capacity. (包含指向数组的指针、长度和容量)
- Slicing does **not** copy the data; changes affect the shared underlying array. (切片操作**不会**复制数据；修改会影响共享的底层数组)
- Use `copy(dst, src)` to create independent slices. (使用 `copy(dst, src)` 创建独立的切片)

```go
slice := make([]int, 5, 10)      // length=5, capacity=10
slice2 := slice[1:3]             // shares the same backing array
slice3 := make([]int, len(slice))
copy(slice3, slice)              // independent copy
```

## Navigation

← [Previous: Maps](./03_maps.md) | [Back to README](../README.md) | **Next:** [Goroutines and Channels](./05_goroutines-channels.md) →