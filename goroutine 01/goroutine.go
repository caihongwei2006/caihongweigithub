package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func main() {

	wg.Add(10001)
	for i := 0; i < 10000; i++ {
		go func(i int) {
			fmt.Println("hello", i)
			wg.Done()
		}(i)
	}
	go func() {
		fmt.Println("hello")
		wg.Done()
	}()
	wg.Wait()

	fmt.Println("end")
}
