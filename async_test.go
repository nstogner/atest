package atest_test

import (
	"atest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func init() {
	atest.Interval = time.Second / 1000
	atest.Duration = time.Second / 100
	atest.Timeout = time.Second / 100
}

type mockT struct {
	failnow      bool
	errorf       bool
	errorfFormat string
	errorfArgs   []interface{}
}

func (t *mockT) Errorf(format string, args ...interface{}) {
	t.errorf = true
	t.errorfFormat = format
	t.errorfArgs = args
}
func (t *mockT) FailNow() {
	t.failnow = true
}

func TestEventuallyFailNowSuccess(t *testing.T) {
	i := 0

	mt := &mockT{}
	atest.Eventually(mt, func(t atest.T) {
		i++
		if i < 3 {
			t.FailNow()
		}
	})

	require.Equal(t, 3, i)
	require.False(t, mt.failnow)
	require.False(t, mt.errorf)
}

func TestEventuallyErrorfSuccess(t *testing.T) {
	i := 0

	mt := &mockT{}
	atest.Eventually(mt, func(t atest.T) {
		i++
		if i < 3 {
			t.Errorf("utoh")
		}
	})

	require.Equal(t, 3, i)
	require.False(t, mt.failnow)
	require.False(t, mt.errorf)
}

func TestEventuallyFailNowFailure(t *testing.T) {
	i := 0

	mt := &mockT{}
	atest.Eventually(mt, func(t atest.T) {
		i++
		t.FailNow()
	})

	require.Greater(t, i, 3)
	require.True(t, mt.failnow)
	require.False(t, mt.errorf)
}

func TestEventuallyErrorfFailure(t *testing.T) {
	i := 0

	mt := &mockT{}
	atest.Eventually(mt, func(t atest.T) {
		i++
		t.Errorf("utoh")
	})

	require.Greater(t, i, 3)
	require.True(t, mt.errorf)
	require.False(t, mt.failnow)
}

func TestConsistentlyFailNowFailure(t *testing.T) {
	i := 0

	mt := &mockT{}
	atest.Consistently(mt, func(t atest.T) {
		i++
		if i >= 3 {
			t.FailNow()
		}
	})

	require.Equal(t, 3, i)
	require.True(t, mt.failnow)
	require.False(t, mt.errorf)
}

func TestConsistentlyErrorfFailure(t *testing.T) {
	i := 0

	mt := &mockT{}
	atest.Consistently(mt, func(t atest.T) {
		i++
		if i >= 3 {
			t.Errorf("fail: %v", 123)
		}
	})

	require.Equal(t, 3, i)
	require.True(t, mt.errorf)
	require.Equal(t, "fail: %v", mt.errorfFormat)
	require.Equal(t, []interface{}{123}, mt.errorfArgs)
	require.False(t, mt.failnow)
}

func TestConsistentlyFailNowSuccess(t *testing.T) {
	i := 0

	mt := &mockT{}
	atest.Consistently(mt, func(t atest.T) {
		i++
	})

	require.Greater(t, i, 3)
	require.False(t, mt.failnow)
	require.False(t, mt.errorf)
}

func TestConsistentlyErrorfSuccess(t *testing.T) {
	i := 0

	mt := &mockT{}
	atest.Consistently(mt, func(t atest.T) {
		i++
	})

	require.Greater(t, i, 3)
	require.False(t, mt.errorf)
	require.False(t, mt.failnow)
}
