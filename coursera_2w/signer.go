package main

import (
	"fmt"
	"time"
)

var combine = ""

func main() {

	now := time.Now()
	for v := range CombineResults(MultiHash(SingleHash(gen(0,1,2,3,4,5,6)))) {
		fmt.Println(v)
	}
	fmt.Println(time.Since(now))

	//combine := ""
	//
	//for _, v := range inputData {
	//
	//	str := fmt.Sprintf("%v", v)
	//	resHash := DataSignerCrc32(str) + "~" + DataSignerCrc32(DataSignerMd5(str))
	//
	//	mulHash := ""
	//	for i := 0; i <= 5; i++ {
	//		mulHash += DataSignerCrc32(fmt.Sprintf("%v", i) + resHash)
	//	}
	//
	//	if combine == "" {
	//		combine += mulHash
	//	} else {
	//		combine += "_" + mulHash
	//	}
	//}
	//
	//fmt.Println(combine)

	fmt.Scanln()
}

func SingleHash(in <-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		for n := range in {
			// processing input
			str := fmt.Sprintf("%v", n)
			resHash := DataSignerCrc32(str) + "~" + DataSignerCrc32(DataSignerMd5(str))

			//send it to the next handler
			out <- resHash
		}
		close(out)
	}()
	return out
}

func MultiHash(in <-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		for n := range in {
			// processing input
			str := fmt.Sprintf("%v", n)
			mulHash := ""
			for i := 0; i <= 5; i++ {
				mulHash += DataSignerCrc32(fmt.Sprintf("%v", i) + str)
			}
			out <- mulHash
		}
		close(out)
	}()
	return out
}

func CombineResults(in <-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		for n := range in {
			// processing input
			str := fmt.Sprintf("%v", n)
			if combine == "" {
				combine += str
			} else {
				combine += "_" + str
			}
		}
		out <- combine
		close(out)
	}()
	return out
}

func gen(nums ...interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}
