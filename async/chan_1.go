package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan int)

	go func(in chan int) {
		val := <-in
		fmt.Println("GO: get from chan", val)
		fmt.Println("GO: after read from chan")
	}(ch1)

	go func() {
		time.Sleep(3 * time.Second)
		ch1 <- 777
		fmt.Println("MAIN: after put to chan")
	}()

	fmt.Scanln()
}
