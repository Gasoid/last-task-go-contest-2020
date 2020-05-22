package main

func Merge2Channels(f func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int) {
	for i := 0; i < n; i++ {
		x1, x2 := <-in1, <-in2
		out <- f(x1) + f(x2)
	}
}
