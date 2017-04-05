package main

import (
	"context"
	"fmt"
	"time"
)

const (
	numThreads = 5
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	var mutexes [numThreads]chan bool

	// launch threads
	for i := 0; i < numThreads; i++ {
		mutexes[i] = make(chan bool, 1)
		mutexes[i] <- true
		go threadedRoutine(ctx, mutexes[i], i)
	}

	for i := 0; i < 3; i++ {
		time.Sleep(time.Second * 5)
		for j := 0; j < numThreads-1; j++ {
			<-mutexes[j]
		}

		// Obtained all mutexes
		fmt.Print("In main loop # ", i, " waiting for ten seconds and no threads should run... ")
		time.Sleep(time.Second * 10)
		fmt.Println("done!")

		// Releases all mutexes
		for j := 0; j < numThreads-1; j++ {
			mutexes[j] <- true
		}
	}

	cancel()
}

func threadedRoutine(ctx context.Context, mutex chan bool, i int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Terminating threadedRoutine # ", i)
			return
		default:
			<-mutex
			fmt.Println("In threadedRoutine # ", i)
			mutex <- true
		}
		time.Sleep(time.Second * 2)
	}
}
