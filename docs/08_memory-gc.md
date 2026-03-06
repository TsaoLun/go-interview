# Memory Management and GC (内存管理与垃圾回收)

- Go uses a concurrent, tri‑color mark‑and‑sweep garbage collector. (Go使用并发的三色标记-清除垃圾回收器)
- **Stack allocation** for small, short‑lived objects. (**栈分配** 用于小型、短寿命的对象)
- **Heap allocation** for objects that escape or are large. (**堆分配** 用于逃逸或大型的对象)
- **Escape analysis** determines whether a variable can be allocated on the stack. (**逃逸分析** 确定变量是否可以在栈上分配)

**GC Tuning (GC调优):**
- `GOGC` environment variable sets the garbage collection target percentage (default 100). (`GOGC` 环境变量设置垃圾回收目标百分比（默认100）)
- Use `runtime.ReadMemStats` to monitor memory usage. (使用 `runtime.ReadMemStats` 监控内存使用情况)

## Navigation

← [Previous: Defer, Panic, and Recover](./07_defer-panic-recover.md) | [Back to README](../README.md) | **Next:** [Concurrency Patterns](./09_concurrency-patterns.md) →