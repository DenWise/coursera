package main

import (
	"fmt"
	"time"
)

func sayHelloo() {
	fmt.Println("Hello")
}

func main() {
	timer := time.AfterFunc(2 * time.Second, sayHelloo)

	fmt.Scanln()
	timer.Stop()
	fmt.Println(timer.C)
}
