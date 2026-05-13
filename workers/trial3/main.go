package main

import (
	"fmt"
	"sync"
)

func fanOut(tasks <-chan int, workers int) []chan int {
	outs := make([]chan int, workers)
	for i := range workers {
		out := make(chan int)
		go func(out chan int) {
			defer close(out)
			for task := range tasks {
				fmt.Printf("worker %v got task %v\n", i, task)
				result := task * 2
				out <- result
			}
		}(out)
		outs[i] = out
	}
	return outs
}

func fanIn(outs []chan int, result chan int, wg *sync.WaitGroup) {
	for _, out := range outs {
		wg.Add(1)
		go func(out chan int) {
			defer wg.Done()
			for o := range out {
				result <- o
			}
		}(out)
	}
}
func main() {
	tasks := make(chan int)
	outs := fanOut(tasks, 3)
	go func() {

		for i := range 10 { // this was sequential bottom go routine never reached, it was still waiting for worker out to be read by the bottom go routine that never ran
			tasks <- i
		}
		close(tasks)
	}()

	var wg sync.WaitGroup
	result := make(chan int)
	fanIn(outs, result, &wg)

	go func() { // has to be in a go routine or else it will block and result channel is never consumed
		wg.Wait()
		close(result)
	}()
	for res := range result { // let this block or else main will end while go routines are still running
		fmt.Println(res)
	}
	fmt.Println("end")

}
