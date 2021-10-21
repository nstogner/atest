package atest

import (
	"runtime"
	"time"
)

var (
	Interval = time.Second / 10
	Timeout  = 10 * time.Second
	Duration = 10 * time.Second
)

type T interface {
	Errorf(format string, args ...interface{})
	FailNow()
}

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
			if remember.failnow {
				t.FailNow()
			} else {
				t.Errorf(remember.errorfFormat, remember.errorfArgs...)
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

func (t *rememberT) failed() bool {
	return t.failnow || t.errorf
}
