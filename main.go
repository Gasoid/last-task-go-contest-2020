package main

import (
	"sync"
)

func main() {

}

var lock sync.Mutex

func Merge2Channels(f func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int) {
	go func(f func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int) {
		done := make(chan bool, n*2)
		results := make([]*int, n*2)

		go func() {
			nDone := 0

			for <-done {
				for i := nDone; i < n; i++ {
					if results[i] != nil && results[i+n] != nil {
						out <- (*results[i] + *results[i+n])
						if nDone++; nDone == n {
							lock.Unlock()
							return
						}
					} else {
						break
					}
				}
			}
		}()
		input := func(ch <-chan int, results []*int) {
			for i := 0; i < n; i++ {
				x := <-ch
				go func(i int, x int) {
					result := f(x)
					results[i] = &result
					done <- true
				}(i, x)
			}
		}
		lock.Lock()
		go input(in1, results[:n])
		go input(in2, results[n:])

	}(f, in1, in2, out)
}
