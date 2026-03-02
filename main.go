package main

import (
	"fmt"
	"sync"
	"time"
)

var counter = 0
var mutex = sync.Mutex{}

func incrementCounter() {
	counter++
}

func main() {
	for range 1000 {
		go incrementCounter()
	}

	time.Sleep(time.Second)
	fmt.Println(counter)
}
