package main

import (
	"fmt"
	"sync"
)

func main() {
	var pool = make(chan chan int)
	var worker = make(chan int)
	var wg sync.WaitGroup

	go func() {
		pool<-worker
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println(<-worker)
	}()

	w := <-pool
	w<-688

	wg.Wait()
}