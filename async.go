package atest

import (
	"runtime"
	"time"
)

var (
	// Interval for retries.
	Interval = time.Second / 10
	// Timeout for eventual tests.
	Timeout = 10 * time.Second
	// Duration for consistent tests.
	Duration = 10 * time.Second
)

// T is an interface targeted at testing.T from the standard library.
type T interface {
	Errorf(format string, args ...interface{})
	FailNow()
	Log(args ...interface{})
	Logf(format string, args ...interface{})
}

// Eventually runs a test repeatedly (based on async.Interval)
// until it passes or until async.Timeout is reached. If Timeout
// is reached the entire test fails.
func Eventually(t T, f func(t T)) {
	t0 := time.Now()

	var (
		errorf       bool
		errorfFormat string
		errorfArgs   []interface{}
	)
	for time.Since(t0) <= Timeout {
		remember := &rememberT{t: t}

		done := make(chan bool)
		go func() {
			defer func() { done <- true }()
			f(remember)
		}()
		<-done

		if !remember.failed() {
			return
		}
		if remember.errorf {
			errorf = true
			errorfFormat = remember.errorfFormat
			errorfArgs = remember.errorfArgs
		}

		time.Sleep(Interval)
	}

	if errorf {
		t.Errorf(errorfFormat, errorfArgs...)
	} else {
		t.FailNow()
	}

	return
}

// Eventually runs a test repeatedly (based on async.Interval)
// until it fails or until async.Duration is reached. If Duration
// is reached the entire test passes.
func Consistently(t T, f func(t T)) {
	t0 := time.Now()

	remember := &rememberT{t: t}

	for time.Since(t0) <= Duration {

		done := make(chan bool)
		go func() {
			defer func() { done <- true }()
			f(remember)
		}()
		<-done

		if remember.failed() {
			if remember.errorfFormat != "" {
				t.Errorf(remember.errorfFormat, remember.errorfArgs...)
			}

			if remember.failnow {
				t.FailNow()
			}
			return
		}

		time.Sleep(Interval)
	}

	return
}

type rememberT struct {
	t T

	failnow      bool
	errorf       bool
	errorfFormat string
	errorfArgs   []interface{}
}

func (t *rememberT) Errorf(format string, args ...interface{}) {
	t.errorf = true
	t.errorfFormat = format
	t.errorfArgs = args
}
func (t *rememberT) FailNow() {
	t.failnow = true
	runtime.Goexit()
}
func (t *rememberT) Log(args ...interface{}) {
	t.t.Log(args...)
}
func (t *rememberT) Logf(format string, args ...interface{}) {
	t.t.Logf(format, args...)
}

func (t *rememberT) failed() bool {
	return t.failnow || t.errorf
}
