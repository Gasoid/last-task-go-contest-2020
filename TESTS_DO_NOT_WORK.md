Привет, есть ошибочный пример реализации, на котором проходят все тесты, но в интерфейсе контеста всё равно "Wrong Answer". Т
о есть тесты не всё тестируют (( вот этот некорректный вариант (не спойлер, потому что не работает):

package main

import "sync"

var mux = sync.Mutex{}
func Merge2Channels(f func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int) {

	go func() {

		resultChannel1 := make(chan int)
		resultChannel2 := make(chan int)

		for i := 0; i < n; i++ {

			mux.Lock()
			arg1,arg2 := <- in1, <- in2


			go func() {
				resultChannel1 <- f(arg1)
			}()

			go func() {
				resultChannel2 <- f(arg2)
			}()

			result1 := <- resultChannel1
			result2 := <- resultChannel2
			out <- result1 + result2
			mux.Unlock()
		}
	}()
}
