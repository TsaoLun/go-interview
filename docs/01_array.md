# Array

- Elements length and type are fixed (元素长度与类型固定)
- Default allocation on stack; escapes to heap when size exceeds `MaxStackVarSize` (默认在栈上分配；当大小超过 `MaxStackVarSize` 时逃逸到堆上)
- Go 1.25: `MaxStackVarSize = int64(128 * 1024)` (older versions: `10 * 1024 * 1024`) (Go 1.25: `MaxStackVarSize = int64(128 * 1024)` (旧版本为 `10 * 1024 * 1024`))

**Escape analysis (逃逸分析):**
```bash
go build -gcflags="-m" .
```

## Navigation

← [Back to README](../README.md) | **Next:** [Mutex and Race Conditions](./02_mutex.md) →