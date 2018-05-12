package tmdwg_test

// TestInvalidTarget tests passing invalid zero target to
import (
	"sync"
	"testing"
	"time"

	"github.com/qbeon/tmdwg-go"
)

// TestInvalidTarget tests the constructor with an invalid target
func TestInvalidTarget(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Fatal("Expected a panic after passing invalid target")
		}
	}()

	tmdwg.NewTimedWaitGroup(0, 1*time.Second)
}

// TestInvalidTimeout tests the constructor with an invalid timeout dur
func TestInvalidTimeout(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Fatal("Expected a panic after passing invalid timeout duration")
		}
	}()

	tmdwg.NewTimedWaitGroup(0, 0)
}

// TestWaitAfterCompletion tests calling Wait after completion
func TestWaitAfterCompletion(t *testing.T) {
	wg := tmdwg.NewTimedWaitGroup(1, 1*time.Second)
	wg.Progress(1)

	if err := wg.Wait(); err != nil {
		t.Fatalf(
			"Expected completed wait group to have no error, got: %s",
			err,
		)
	}

	if !wg.IsCompleted() {
		t.Fatalf("Expected wait group to be completed")
	}

	currentProgress := wg.CurrentProgress()
	if currentProgress != 1 {
		t.Fatalf("Unexpected progress: %d", currentProgress)
	}
}

// TestWaitAfterTimeout tests calling Wait after timeout
func TestWaitAfterTimeout(t *testing.T) {
	wg := tmdwg.NewTimedWaitGroup(1, 10*time.Millisecond)

	// Wait for the timeout
	time.Sleep(50 * time.Millisecond)

	if err := wg.Wait(); err == nil {
		t.Fatal("Expected wait group to have timed out before")
	}
}

// TestWaitTimeout tests the Wait method
// expecting it to return a timeout error on timeout
func TestWaitTimeout(t *testing.T) {
	wg := tmdwg.NewTimedWaitGroup(1, 1*time.Millisecond)

	// Wait for 1 millisecond, then timeout and return err
	if err := wg.Wait(); err == nil {
		t.Fatal("Expected wait group to return an error")
	}
}

// TestProgressHigherThanTarget tests calling Progress with the delta parameter
// being higher than the actual target
func TestProgressHigherThanTarget(t *testing.T) {
	wg := tmdwg.NewTimedWaitGroup(1, 1*time.Second)
	wg.Progress(5)

	if err := wg.Wait(); err != nil {
		t.Fatalf(
			"Expected completed wait group to have no error, got: %s",
			err,
		)
	}

	if !wg.IsCompleted() {
		t.Fatalf("Expected wait group to be completed")
	}

	if wg.CurrentProgress() != 5 {
		t.Fatalf("Expected current progress to be 5")
	}
}

// TestProgressUpdates tests the CurrentProgress method
func TestProgressUpdates(t *testing.T) {
	wg := tmdwg.NewTimedWaitGroup(3, 1*time.Second)

	currentProgress0 := wg.CurrentProgress()
	if currentProgress0 != 0 {
		t.Fatalf("Unexpected progress: %d", currentProgress0)
	}
	if wg.IsCompleted() {
		t.Fatal("Expected wait group to not have yet been completed")
	}

	progress1Result := wg.Progress(1)
	if progress1Result != 1 {
		t.Fatalf("Unexpected progress: %d", progress1Result)
	}

	currentProgress1 := wg.CurrentProgress()
	if currentProgress1 != 1 {
		t.Fatalf("Unexpected progress: %d", currentProgress1)
	}
	if wg.IsCompleted() {
		t.Fatal("Expected wait group to not have yet been completed")
	}

	progress2Result := wg.Progress(2)
	if progress2Result != 3 {
		t.Fatalf("Unexpected progress: %d", progress2Result)
	}

	currentProgress3 := wg.CurrentProgress()
	if currentProgress3 != 3 {
		t.Fatalf("Unexpected progress: %d", currentProgress3)
	}
	if !wg.IsCompleted() {
		t.Fatal("Expected wait group to be completed")
	}
}

// TestProgressNotReachedTimeout tests unreached progress timeout
func TestProgressNotReachedTimeout(t *testing.T) {
	wg := tmdwg.NewTimedWaitGroup(2, 1*time.Millisecond)
	wg.Progress(1)

	// Expect wait to directly return without any error
	if err := wg.Wait(); err == nil {
		t.Fatal("Expected wait group to return an error")
	}
}

// TestConcurrentProgress tests multiple concurrent Progress calls
func TestConcurrentProgress(t *testing.T) {
	concurrentProgressers := 16

	wg := tmdwg.NewTimedWaitGroup(concurrentProgressers, 1*time.Second)

	for i := 0; i < concurrentProgressers; i++ {
		go func() {
			wg.Progress(1)
		}()
	}

	// Wait for 1 second and expect not to receive a timeout error
	if err := wg.Wait(); err != nil {
		t.Fatal("Expected wait group to have been completed")
	}
}

// TestConcurrentWait tests multiple concurrent Wait calls
func TestConcurrentWait(t *testing.T) {
	concurrentWaiters := 16

	var swg sync.WaitGroup
	swg.Add(concurrentWaiters)

	wg := tmdwg.NewTimedWaitGroup(1, 1*time.Second)

	for i := 0; i < concurrentWaiters; i++ {
		go func() {
			// Wait for 1 millisecond, then timeout and return err
			if err := wg.Wait(); err != nil {
				t.Error("Expected wait group to have been completed")
			}
			swg.Done()
		}()
	}

	wg.Progress(1)

	// Wait for all concurrent waiters to have passed Wait
	swg.Wait()
}

// TestConcurrentTimeout tests timeout of multiple concurrent Wait calls
func TestConcurrentTimeout(t *testing.T) {
	concurrentWaiters := 16
	errsLock := sync.Mutex{}
	errs := make([]error, 0, concurrentWaiters)

	var swg sync.WaitGroup
	swg.Add(concurrentWaiters)

	wg := tmdwg.NewTimedWaitGroup(1, 1*time.Millisecond)

	for i := 0; i < concurrentWaiters; i++ {
		go func() {
			if err := wg.Wait(); err != nil {
				errsLock.Lock()
				errs = append(errs, err)
				errsLock.Unlock()
			}
			swg.Done()
		}()
	}

	// Wait for all waiters to have received a timeout error
	swg.Wait()

	// Expect all Wait calls to have received a timeout error
	if len(errs) != concurrentWaiters {
		t.Fatalf(
			"Only %d of %d have returned an error",
			len(errs),
			concurrentWaiters,
		)
	}
}
