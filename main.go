package main

import (
	"sync"
)

var mutex sync.Mutex

func Merge2Channels(f func(int) int, in1 <-chan int, in2 <- chan int, out chan<- int, n int){

	go func(){
		mutex.Lock()
		defer mutex.Unlock()

		res := make([]int, n)
		var wg sync.WaitGroup

		for i := 0; i < n; i++ {
			wg.Add(1)
			go func(x1 int, x2 int, ix int ) {
				res[ix] = f(x1) + f(x2)
              	//fmt.Printf("[%d] == (f(%d)+f(%d)) == %d\n", ix, x1, x2, res[ix])
				wg.Done()
			}(<- in1, <- in2, i)
		}
		wg.Wait()
		for i := 0; i < n; i++{
			out <- res[i]
		}
	}()
}


