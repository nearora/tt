package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const (
	numThreads = 5
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	var mutexes [numThreads]*sync.Mutex

	// launch threads
	for i := 0; i < numThreads; i++ {
		mutexes[i] = &sync.Mutex{}
		go threadedRoutine(ctx, mutexes[i], i)
	}

	for i := 0; i < 3; i++ {
		time.Sleep(time.Second * 5)
		for j := 0; j < numThreads; j++ {
			mutexes[j].Lock()
		}

		// Obtained all mutexes
		fmt.Print("In main loop # ", i, " waiting for ten seconds and no threads should run... ")
		time.Sleep(time.Second * 10)
		fmt.Println("done!")

		// Releases all mutexes
		for j := 0; j < numThreads; j++ {
			mutexes[j].Unlock()
		}
	}

	cancel()
}

func threadedRoutine(ctx context.Context, mutex *sync.Mutex, i int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Terminating threadedRoutine # ", i)
			return
		default:
			mutex.Lock()
			fmt.Println("In threadedRoutine # ", i)
			mutex.Unlock()
		}
		time.Sleep(time.Second * 2)
	}
}
