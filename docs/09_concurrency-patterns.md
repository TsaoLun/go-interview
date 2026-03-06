# Concurrency Patterns (并发模式)

## Worker Pool (工作池)

```go
func worker(id int, jobs <-chan int, results chan<- int) {
    for job := range jobs {
        results <- job * 2
    }
}

func main() {
    jobs := make(chan int, 100)
    results := make(chan int, 100)

    // Start workers
    for w := 1; w <= 3; w++ {
        go worker(w, jobs, results)
    }

    // Send jobs
    for j := 1; j <= 9; j++ {
        jobs <- j
    }
    close(jobs)

    // Collect results
    for r := 1; r <= 9; r++ {
        fmt.Println(<-results)
    }
}
```

## Pipeline (管道)

```go
func gen(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()
    return out
}

func sq(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * n
        }
        close(out)
    }()
    return out
}

func main() {
    c := gen(2, 3, 4)
    out := sq(c)
    for n := range out {
        fmt.Println(n) // 4, 9, 16
    }
}
```

## Navigation

← [Previous: Memory Management and GC](./08_memory-gc.md) | [Back to README](../README.md) | **Next:** [Examples](./10_examples.md) →