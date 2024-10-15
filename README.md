[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Coverage Status](https://coveralls.io/repos/github/romshark/tmdwg-go/badge.svg?branch=master)](https://coveralls.io/github/romshark/tmdwg-go?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/romshark/tmdwg-go)](https://goreportcard.com/report/github.com/romshark/tmdwg-go)

# Timed WaitGroup for Go

[tmdwg-go](https://github.com/romshark/tmdwg-go) provides a **timed wait group** implementation similar to **[sync.WaitGroup](https://golang.org/pkg/sync/#WaitGroup)**.


It's purpose is simple: it blocks all goroutines that called it's `wg.Wait()` method and frees them when:
- either the timeout is reached...
- or the progress is reached

In case the timeout was reached before the progress `wg.Wait()` will return a timeout error, otherwise it'll return `nil`.

The timed wait group is fully thread safe and may safely be used concurrently from within multiple goroutines.

🧟‍♂️ Written in May 2018, revived in October 2024.
