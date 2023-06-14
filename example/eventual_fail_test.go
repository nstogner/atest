package main

import (
	"testing"
	"time"

	"github.com/nstogner/atest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	atest.Interval = time.Second / 1000
	atest.Duration = time.Second / 100
	atest.Timeout = time.Second / 100
}

func TestEventualFailNow(t *testing.T) {
	t.Log("Before")

	atest.Eventually(t, func(t atest.T) {
		require.Equal(t, true, false)
	})

	t.Log("After: This should NOT be logged")
}

func TestEventualFailContinue(t *testing.T) {
	t.Log("Before")

	atest.Eventually(t, func(t atest.T) {
		assert.Equal(t, true, false)
	})

	t.Log("After: This SHOULD be logged")
}
