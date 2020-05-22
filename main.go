package main

import "sync"

func main() {

}

func Merge2Channels(f func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int) {
	go func(f func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int) {
		var lock1, lock2 sync.Mutex
		done := make(chan bool, n*2)
		results := make([]*int, n*2)

		go func() {
			nDone := 0
			defer close(out)
			for <-done {
				for i := nDone; i < n; i++ {
					if results[i] != nil && results[i+n] != nil {
						out <- (*results[i] + *results[i+n])
						if nDone++; nDone == n {
							return
						}
					} else {
						break
					}
				}
			}
		}()
		input := func(ch <-chan int, results []*int, lock *sync.Mutex) {
			for i := 0; i < n; i++ {
				x := <-ch
				go func(i int, x int) {
					lock.Lock()
					result := f(x)
					lock.Unlock()
					results[i] = &result
					done <- true
				}(i, x)
			}
		}
		go input(in1, results[:n], &lock1)
		go input(in2, results[n:], &lock2)

	}(f, in1, in2, out)
}
