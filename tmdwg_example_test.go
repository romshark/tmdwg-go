package tmdwg_test

import (
	"fmt"
	"time"

	"github.com/romshark/tmdwg-go"
)

func Example() {
	waitTime := 1 * time.Second

	// Create a new wait group that times out after 1 second
	// and sets a progress target of 4
	wg := tmdwg.NewTimedWaitGroup(4, waitTime)

	// Progress the wait group from another goroutine by 1
	go func() { wg.Progress(1) }()

	// Progress the wait group from another goroutine by 3
	go func() { wg.Progress(3) }()

	// Wait for the wait group to either timeout or complete
	if err := wg.Wait(); err != nil {
		fmt.Printf("Wait group timed out after %s\n", waitTime)
	}
	fmt.Printf(
		"Wait group completed (%d) within %s\n",
		wg.CurrentProgress(),
		waitTime,
	)

	// Output:
	// Wait group completed (4) within 1s
}
