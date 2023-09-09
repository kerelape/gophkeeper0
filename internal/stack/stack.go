package stack

// Stack is a convenience wrapper for slice.
type Stack[T any] []T

// Pop returns top element of the stack.
func (s *Stack[T]) Pop() T {
	var element = (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return element
}

// Push pushes value to the top of the stack.
func (s *Stack[T]) Push(value T) {
	*s = append(*s, value)
}
