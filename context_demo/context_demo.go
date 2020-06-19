package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func wait_group_example() {
	fmt.Println("---WaitGroup Example---")
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		fmt.Println("Task 1 starts")
		time.Sleep(100)
		fmt.Println("Task 1 finishes!")
		wg.Done()
	}()

	go func() {
		fmt.Println("Task 2 starts")
		time.Sleep(90)
		fmt.Println("Task 2 finishes!")
		wg.Done()
	}()

	fmt.Println("Waiting for tasks to finish...")
	wg.Wait()
	fmt.Println("All tasks are done!")
}

func channel_example() {
	fmt.Println("---Channel Example---")
	stop := make(chan bool)

	go func() {
		for {
			select {
			case <-stop:
				fmt.Println("[Daemon] Receive stop signal, exit")
				return

			default:
				fmt.Println("[Daemon] Working...")
				time.Sleep(2 * time.Second)
			}
		}
	}()

	time.Sleep(5 * time.Second)
	stop <- true
	fmt.Println("[Main] Daemon has been stopped.")
	time.Sleep(3 * time.Second)
}

func context_example() {
	ctx, cancel := context.WithCancel(context.Background())
	go watch(ctx, "W1")
	go watch(ctx, "W2")
	go watch(ctx, "W3")

	time.Sleep(5 * time.Second)
	fmt.Println("Now cancel all watching")

	// cancel a group of child goroutines
	cancel()
	time.Sleep(3 * time.Second)
}

func watch(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("[%v] watching is done\n", name)
			return
		default:
			fmt.Printf("[%v] watching is ongoing\n", name)
			time.Sleep(time.Second)
		}
	}
}

func main() {
	// wait_group()
	// channel()
	context_example()
}
