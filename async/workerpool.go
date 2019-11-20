package main

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

const goroutinesNum = 3

func startWorker(workerNum int, in <-chan string) {
	for input := range in {
		fmt.Printf(formatWork(workerNum, input))
		runtime.Gosched() // попробуйте закомментировать
	}
	printFinishWork(workerNum)
}

func main() {
	runtime.GOMAXPROCS(0)               // попробуйте с 0 (все доступные) и 1
	worketInput := make(chan string, 2) // попробуйте увеличить размер канала
	for i := 0; i < goroutinesNum; i++ {
		go startWorker(i, worketInput)
	}

	months := []string{"Январь", "Февраль", "Март",
		"Апрель", "Май", "Июнь",
		"Июль", "Август", "Сентябрь",
		"Октябрь", "Ноябрь", "Декабрь",
	}

	years := []string{"2009","2010","2011"}

	for _, monthName := range months {
		worketInput <- monthName
	}

	for _, yearNum := range years {
		worketInput <- yearNum
	}

	close(worketInput) // попробуйте закомментировать, если не закрыть канал воркеров, то цикл в воркере не завершится и все что должно произойти по завершении цикла не произойдет

	time.Sleep(time.Millisecond)
}

func formatWork(in int, input string) string {
	return fmt.Sprintln(strings.Repeat("  ", in), "█",
		strings.Repeat("  ", goroutinesNum-in),
		"th", in,
		"recieved", input)
}

func printFinishWork(in int) {
	fmt.Println(strings.Repeat("  ", in), "█",
		strings.Repeat("  ", goroutinesNum-in),
		"===", in,
		"finished")
}
