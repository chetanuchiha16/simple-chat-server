package main

import "fmt"

func fanOut(tasks <-chan int, workers int) []chan int {
	outs := make([]chan int, workers)
	for i := range workers {
		go func() {
			for task := range tasks {
				fmt.Printf("worker %v got task %v\n", i, task)
			}
		}()
	}
	return outs
}

func main() {
	tasks := make(chan int)
	fanOut(tasks, 3)
	for i := range 10 {
		tasks <- i
	}
	close(tasks)

}
