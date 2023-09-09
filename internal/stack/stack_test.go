package stack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	t.Run("Push one element to an empty stack", func(t *testing.T) {
		var s = make(Stack[string], 0)
		s.Push("Hello, World!")
		assert.Len(t, s, 1, "expected stack length to be equal to 1")
		assert.Contains(t, s, "Hello, World!", "expected stack to contain pushed element")
	})
	t.Run("Pop one element from a unit length stack", func(t *testing.T) {
		var s = make(Stack[string], 0)
		s.Push("Hello, World!")
		assert.Equal(t, "Hello, World!", s.Pop(), "popped element doesn't match the expected")
		assert.Len(t, s, 0, "stack must be 0 len after popping its only element")
	})
}
