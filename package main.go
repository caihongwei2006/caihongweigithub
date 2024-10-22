package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	c := make(chan int)
	for i := 0; i < 5; i++ {
		go gopher(i, c)
	}
	timeout := time.After(4 * time.Second)
	select {
	case gopherid := <-c:
		fmt.Println("gooher ID:", gopherid, "has finished sleeping")
	case <-timeout:
		{
			fmt.Println("gopher is still sleepy")
		}
	}
}

func gopher(i int, c chan int) {
	k := time.Duration(rand.Intn(4000))
	time.Sleep(k * time.Millisecond)
	fmt.Println("slept for", k, "seconds")
	c <- i
}
