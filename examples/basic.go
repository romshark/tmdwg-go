package main

import (
	"log"
	"time"

	"github.com/qbeon/tmdwg-go"
)

func main() {
	waitTime := 1 * time.Second

	// Create a new wait group that times out after 1 second
	// and sets a progress target of 4
	wg := tmdwg.NewTimedWaitGroup(4, waitTime)

	// Progress the wait group from another goroutine by 1
	go func() {
		wg.Progress(1)
	}()

	// Progress the wait group from another goroutine by 3
	go func() {
		wg.Progress(3)
	}()

	// Wait for the wait group to either timeout or complete
	if err := wg.Wait(); err != nil {
		log.Printf("Wait group timed out after %s", waitTime)
	}
	log.Printf(
		"Wait group completed (%d) within %s",
		wg.CurrentProgress(),
		waitTime,
	)
}
