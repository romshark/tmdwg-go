# Timed WaitGroup for Go

[tmdwg-go](https://github.com/qbeon/tmdwg-go) provides a **timed wait group** implementation similar to **[sync.WaitGroup](https://golang.org/pkg/sync/#WaitGroup)**.

It's purpose is simple: it blocks all goroutines that called it's `wg.Wait()` method and frees them when:
- either the timeout is reached...
- or the progress is reached

In case the timeout was reached before the progress `wg.Wait()` will return a timeout error, otherwise it'll return `nil`.

The timed wait group is fully thread safe and may safely be used concurrently from within multiple goroutines.
