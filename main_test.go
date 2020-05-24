package main

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"
)

func TestMerge2Channels1(t *testing.T) {
	f := func(n int) int {
		sleep := rand.Int31n(100)
		// log.Println("run =", n)
		time.Sleep(time.Duration(sleep) * time.Millisecond)
		//time.Sleep(2 * time.Second)
		return (n * n)
	}
	rand.Seed(12000)
	// repeats := rand.Intn(400)
	repeats := 100
	in1 := make(chan int, repeats)
	in2 := make(chan int, repeats)
	out := make(chan int, repeats)
	// log.Println("Seed")
	log.Println("Merge2Channels")
	Merge2Channels(f, in1, in2, out, repeats/2)
	Merge2Channels(f, in1, in2, out, repeats/2)
	results := []int{}

	go func() {
		for i := 0; i < repeats; i++ {
			i1 := rand.Intn(200)
			i2 := rand.Intn(200)
			in1 <- i1
			in2 <- i2
			results = append(results, (i1*i1)+(i2*i2))
		}
	}()
	c := 0
	log.Println("range out")
	for i := range out {
		if i != results[c] {
			t.Errorf("%v != %v", i, results[c])
		}
		//log.Println("result =", results[c])
		//log.Println("out =", i)
		c++
		if c == repeats {
			close(out)
		}
	}
	log.Println("REPEATS=", repeats)
	// log.Println("results =", results)
}

func square(n int) int {
	time.Sleep(time.Duration(rand.Int31n(10)) * time.Millisecond)
	return n * n
}

func TestMerge2Channels2(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	repeats := 30
	done := make(chan struct{}, 2)
	runTimes := 2
	for i := 0; i < runTimes; i++ {
		go func() {
			in1 := make(chan int, 100)
			in2 := make(chan int, 100)
			out := make(chan int, 100)

			var expectedOut []int

			for i := 1; i < 101; i++ {
				in1 <- i
				in2 <- i
				expectedOut = append(expectedOut, square(i)*2)
			}
			Merge2Channels(square, in1, in2, out, repeats)
			go func(expectedResult []int, out <-chan int, done chan<- struct{}) {
				for i := 0; i < repeats; i++ {
					v := expectedOut[i]
					r := <-out
					if v != r {
						t.Error("ОЖИДАЛ:", v, "ПОЛУЧИЛ:", r)
					}
				}
				done <- struct{}{}
			}(expectedOut, out, done)
		}()

	}
	for i := 0; i < runTimes; i++ {
		<-done
	}

}

func fastSquare(a int) int {
	return a * a
}

func slowSquare(a int) int {
	time.Sleep(50 * time.Millisecond)
	return a * a
}

func TestNonBlocking(t *testing.T) {
	capacity := 100

	ch1 := make(chan int, capacity)
	a1 := make([]int, capacity)

	ch2 := make(chan int, capacity)
	a2 := make([]int, capacity)

	chOut := make(chan int, capacity)
	a3 := make([]int, capacity)

	i := 0
	for i < capacity {
		a1[i] = i + 9
		ch1 <- a1[i]

		a2[i] = i*3 + 289
		ch2 <- a2[i]

		a3[i] = fastSquare(a1[i]) + fastSquare(a2[i])
		i++
	}

	done := make(chan struct{})

	portion := 30

	go func() {
		Merge2Channels(slowSquare, ch1, ch2, chOut, portion)
		close(done)
	}()

	select {
	case <-done:

	case <-time.After(time.Millisecond * 100):
		t.Fail()
		panic("Function should be non-blocking")
	}
}

func TestSlowSquare(t *testing.T) {
	capacity := 100

	ch1 := make(chan int, capacity)
	a1 := make([]int, capacity)

	ch2 := make(chan int, capacity)
	a2 := make([]int, capacity)

	chOut := make(chan int, capacity)
	a3 := make([]int, capacity)

	i := 0
	for i < capacity {
		a1[i] = i + 9
		ch1 <- a1[i]

		a2[i] = i*3 + 289
		ch2 <- a2[i]

		a3[i] = fastSquare(a1[i]) + fastSquare(a2[i])
		i++
	}

	done := make(chan struct{})

	portion := 30

	go func() {
		Merge2Channels(slowSquare, ch1, ch2, chOut, portion)
		close(done)
	}()

	<-done

	i = 0
	for i < portion {
		ans, ok := <-chOut
		if !ok {
			t.Fail()
			panic("Output channel closed prematurely")
		}
		if ans != a3[i] {
			t.Fail()
			panic(fmt.Errorf("Got %d from output channel, should be %d", ans, a3[i]))
		}
		i++
	}

	if len(ch1) != capacity-portion {
		t.Fail()
		panic(fmt.Errorf("First channel has %d numbers in it, should have %d", len(ch1), capacity-portion))
	}
	if len(ch2) != capacity-portion {
		t.Fail()
		panic(fmt.Errorf("Second channel has %d numbers in it, should have %d", len(ch2), capacity-portion))
	}
}

func TestClosedInputChannels(t *testing.T) {
	capacity := 100
	in1 := make(chan int, capacity)
	in2 := make(chan int, capacity)
	out := make(chan int, capacity)

	portion := 30

	go func() {
		for i := 0; i < portion; i++ {
			i1 := rand.Intn(200)
			i2 := rand.Intn(200)
			in1 <- i1
			in2 <- i2
		}
		close(in1)
		close(in2)
	}()

	runs := portion * 2
	Merge2Channels(slowSquare, in1, in2, out, runs)

	count := 0
	for count < runs {
		_, ok := <-out
		if !ok {
			break
		} else {
			count++
		}
	}
	if count != portion {
		t.Fail()
		panic(fmt.Errorf("Out channel had %d values instead of expected %d", count, portion))
	}
}
