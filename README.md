# Async Test Helper

**Archived**: Take a look at the Eventually() in https://github.com/stretchr/testify
**PS**: There is a issue open for adding a Consistently(): https://github.com/stretchr/testify/issues/1087

Helpers for testing async behavior.

Something that should eventually succeed:

```go
import "github.com/stretchr/testify/require"

func TestSomething(t *testing.T) {
	rand.Seed(1)
	r := rand.Intn(100)

	atest.Eventually(t, func(t atest.T) {
		require.Less(t, r, 50)
	})
}
```

Something that should consistently succeed:

```go
import "github.com/stretchr/testify/require"

func TestSomething(t *testing.T) {
	rand.Seed(1)
	r := rand.Intn(100)

	atest.Consistently(t, func(t atest.T) {
		require.Less(t, r, 50)
	})
}
```

Global configuration:

```go
func init() {
	atest.Interval = time.Second / 1000
	atest.Duration = time.Second / 100
	atest.Timeout = time.Second / 100
}
```
