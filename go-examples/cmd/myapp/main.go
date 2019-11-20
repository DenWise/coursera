package main

import (
	"fmt"
	"os"
	"bufio"
	"io"
)

func defFunc (s string) {
	fmt.Printf("we deferring at the end, and value is %q\n",s)
}

func worker() int {
	str := "Start value"
	defer func() {
		defFunc(str)
	}()
	str += " is changed"
	return 1
}

type MySlice []int

func (sl *MySlice) Add(val... int) {
	for _, v := range val {
		*sl = append(*sl,v)
	}
}

func (sl *MySlice) Count() int {
	return len(*sl)
}

func uniq(input io.Reader, output io.Writer) error {
	in := bufio.NewScanner(input)
	var prev string
	for in.Scan() {
		txt := in.Text()

		if txt == prev {
			continue
		}

		if txt < prev {
			return fmt.Errorf("file not sorted")
		}

		prev = txt

		fmt.Fprintf(output,"%v\n",txt)
	}
	return nil
}

func main() {
	worker()
	sl := MySlice([]int{1, 2, 3})
	sl.Add(4)
	sl.Add()
	fmt.Println(sl.Count(),sl)

	err := uniq(os.Stdin,os.Stdout)
	if err != nil {
		fmt.Println("There is some error while running")
	}
}