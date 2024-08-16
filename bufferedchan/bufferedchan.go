package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var ch = make(chan int)
	var wg = new(sync.WaitGroup)

	for i:=0; i<5; i++ {
		wg.Add(1)
		go func(i int, ch chan int, wg *sync.WaitGroup) {
			defer wg.Done()
			fmt.Println("blocking...")
			ch<-i
			fmt.Println("sent")
		}(i, ch, wg)
	}

	go func(wg *sync.WaitGroup, ch chan int) {
		fmt.Println("waiting")
		wg.Wait()
		fmt.Println("done waiting")
		close(ch)

	}(wg, ch)

	for msg := range ch {
		wg.Add(1)
		go func(msg int) {
			defer wg.Done()
			time.Sleep(2000 * time.Millisecond)
			fmt.Println("> ", msg)
		}(msg)
	}
}