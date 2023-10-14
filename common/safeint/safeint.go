package safeint

import "sync"

// SafeInt is a thread safe int
type SafeInt struct {
	val int
	sync.Mutex
}

// Get the value inside SafeInt
func (s *SafeInt) Get() int {
	s.Lock()
	defer s.Unlock()
	return s.val
}

// Add a value to the SafeInt
func (s *SafeInt) Add(addend int) int {
	s.Lock()
	defer s.Unlock()
	s.val += addend
	return s.val
}

// Set the value of the SafeInt
func (s *SafeInt) Set(val int) int {
	s.Lock()
	defer s.Unlock()
	s.val = val
	return s.val
}

// Set the value to zero
func (s *SafeInt) SetZero() int {
	return s.Set(0)
}

// Increment
func (s *SafeInt) Increment() int {
	return s.Add(1)
}

func New(val int) *SafeInt {
	return &SafeInt{
		val: val,
	}
}
