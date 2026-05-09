package main

import (
	"fmt"
	"sync"
)

func fanOut(tasks <-chan int, workers int, wg *sync.WaitGroup) []chan int {
	outs := make([]chan int, workers)
	for i := range workers {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for task := range tasks {
				fmt.Printf("worker %v got task %v\n", id, task)
			}
		}(i)
	}
	return outs
}

func main() {
	tasks := make(chan int)
	var wg sync.WaitGroup
	fanOut(tasks, 3, &wg)
	for i := range 10 {
		tasks <- i
	}
	close(tasks)
	wg.Wait()

}
