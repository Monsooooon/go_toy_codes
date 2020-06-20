package main

import (
	"fmt"
)

// Demonstrate how to use nil channels to DISABLE a select case
func main() {
	a := make(chan int)
	b := make(chan int)
	out := make(chan int)

	go func() {
		for i := 0; i < 5; i++ {
			a <- 1
		}
		close(a)
	}()

	go func() {
		for i := 0; i < 10; i++ {
			b <- 2
		}
		close(b)
	}()

	go receive_v4(a, b, out)

	for num := range out {
		fmt.Println(num)
	}
}

// Get infinite 0s because when a or b is closed, v will be 0!
func receive_v1(a, b <-chan int, out chan<- int) {
	for {
		select {
		case v := <-a: // after a is closed, v will be 0
			out <- v
		case v := <-b: // after b is closed, v will be 0
			out <- v
		}
	}
}

func receive_v2(a, b <-chan int, out chan<- int) {
	var aClosed, bClosed bool
	for !aClosed && !bClosed {
		select {
		case v, ok := <-a:
			if !ok { // when we find a is closed, set aClosed = true
				aClosed = true
				continue
			}
			out <- v
		case v, ok := <-b:
			if !ok { // when we find a is closed, set bClosed = true
				bClosed = true
				continue
			}
			out <- v
		}
	}
	close(out)
}

// add logging to v2
// busy waiting!!!
func receive_v3(a, b <-chan int, out chan<- int) {
	var aClosed, bClosed bool
	for !aClosed || !bClosed {
		select {
		case v, ok := <-a:
			if !ok {
				aClosed = true
				fmt.Println("a is closed")
				continue
			}
			out <- v
		case v, ok := <-b:
			if !ok {
				bClosed = true
				fmt.Println("b is closed")
				continue
			}
			out <- v
		}
	}
	close(out)
}

// best version!
// since if we call <- on a nil channel, it will block
// we can do this to avoid busy waiting
func receive_v4(a, b <-chan int, out chan<- int) {
	for a != nil || b != nil {
		select {
		case v, ok := <-a:
			if !ok {
				a = nil
				fmt.Println("a is closed")
				continue
			}
			out <- v
		case v, ok := <-b:
			if !ok {
				b = nil
				fmt.Println("b is closed")
				continue
			}
			out <- v
		}
	}
	close(out)
}
