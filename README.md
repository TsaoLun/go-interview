# Go Interview Brief

A comprehensive collection of Go (Golang) concepts and examples for interview preparation.

This documentation is organized into individual numbered sections for focused learning and sequential navigation. Each section can be accessed independently, and all sections include simplified navigation links (previous, next, and back to index).

## Table of Contents

- [Array](./docs/01_array.md) - Fixed-length sequences and stack allocation
- [Mutex and Race Conditions](./docs/02_mutex.md) - Synchronization and concurrency control
- [Maps](./docs/03_maps.md) - Hash table implementation and behavior
- [Slices](./docs/04_slices.md) - Dynamic array segments and memory management
- [Goroutines and Channels](./docs/05_goroutines-channels.md) - Concurrency primitives
- [Interfaces](./docs/06_interfaces.md) - Polymorphism and type abstraction
- [Defer, Panic, and Recover](./docs/07_defer-panic-recover.md) - Error handling and flow control
- [Memory Management and GC](./docs/08_memory-gc.md) - Memory allocation and garbage collection
- [Concurrency Patterns](./docs/09_concurrency-patterns.md) - Common patterns for concurrent programming
- [Examples](./docs/10_examples.md) - Practical code examples and demonstrations
- [References](./docs/11_references.md) - Official documentation and learning resources

## Quick Overview

### Array
- Elements length and type are fixed
- Default stack allocation; escapes to heap when size exceeds `MaxStackVarSize`
- Go 1.25: `MaxStackVarSize = int64(128 * 1024)` (older versions: `10 * 1024 * 1024`)

### Mutex and Race Conditions
- Race conditions occur when multiple goroutines access shared data without synchronization
- Use `sync.Mutex` or `sync.RWMutex` for synchronization
- Mutex implements fast path (atomic CAS) and slow path (spinning/queueing)

### Maps
- Go maps are hash tables composed of buckets
- Maps are not reference types; they are pointers to underlying `hmap` structure
- Growth strategies: double bucket count or redistribute entries

### Slices
- Slices are descriptors of contiguous array segments
- Contain pointer to array, length, and capacity
- Slicing does not copy data; changes affect shared underlying array

### Goroutines and Channels
- Goroutines are lightweight threads managed by Go runtime
- Channels provide communication and synchronization between goroutines
- Select statement for non-blocking communication

### Interfaces
- Define method signatures; types implement interfaces implicitly
- Empty interface (`interface{}` / `any`) holds any type
- Interface values consist of type and value components

### Defer, Panic, and Recover
- `defer` schedules function calls for execution when surrounding function returns
- `panic` stops normal execution and begins panicking
- `recover` regains control of a panicking goroutine (only in deferred functions)

### Memory Management and GC
- Concurrent, tri-color mark-and-sweep garbage collector
- Stack allocation for small, short-lived objects
- Heap allocation for objects that escape or are large
- Escape analysis determines stack vs heap allocation

### Concurrency Patterns
- Worker pool for distributing work among goroutines
- Pipeline pattern for processing data through stages

### Examples
- Example code in `main.go` demonstrates maps are not reference types
- Shows reassigning map inside function doesn't affect caller's variable
- Output: `true` indicating the map remains `nil`

### References
- Official Go language specification and documentation
- Go by Example tutorials and effective Go practices
- Memory model and concurrency patterns references
- Technical articles about maps, mutex, and Go internals

## Navigation

Each section file includes:
- **Back to README** link to return to this index
- **Previous** and **Next** links for sequential reading
- **Numbered sequence**: Files are numbered 01-11 for logical progression (01_array.md, 02_mutex.md, etc.)
- Consistent structure with clear headings and bilingual explanations

To navigate from any section back to this index, click the "Back to README" link at the bottom of each section page.

## Getting Started

Explore the concepts by clicking on any topic in the table of contents. The documentation includes both English and Chinese explanations for key concepts.

To run the example code:
```bash
go run main.go
```

## Contributing

This is a living collection of Go concepts. Contributions and corrections are welcome. Each section is maintained as a separate Markdown file in the `docs/` directory.

---
*Documentation split into sections for focused learning.*