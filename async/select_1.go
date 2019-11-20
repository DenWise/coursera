package main

import (
	"fmt"
)

func main() {
	ch1 := make(chan int, 2)
	ch2 := make(chan int)
	ch1 <- 1
	ch1 <- 4

	select {
	case val := <-ch1:
		fmt.Println("ch1 val", val)
	case ch2 <- 1:
		fmt.Println("put val to ch2")
		close(ch2)
	default:
		fmt.Println("default case")
	}

	//v := <-ch2
	//fmt.Println(v)
}
