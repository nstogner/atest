package atest_test

import (
	"math/rand"
	"testing"

	"github.com/nstogner/atest"

	"github.com/stretchr/testify/require"
)

func ExampleTestEventually(t *testing.T) {
	rand.Seed(1)
	r := rand.Intn(100)

	atest.Eventually(t, func(t atest.T) {
		require.Less(t, r, 50)
	})
}

func ExampleTestConsistently(t *testing.T) {
	rand.Seed(1)
	r := rand.Intn(100)

	atest.Consistently(t, func(t atest.T) {
		require.Less(t, r, 50)
	})
}
