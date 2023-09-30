package safeset

import "sync"

// copy by https://gist.github.com/bbengfort/2470a7b3174a2142417b75ade73edf41

//===========================================================================
// Base Set
//===========================================================================

// Provides a common internal structure for thread-safe and non-thread-safe
// Set objects, as well as internal helper methods used by those objects.
//
// Items is a mapping of an interface (the set member) to an empty value,
// struct{}, which doesn't take up any space.
type set struct {
	items map[interface{}]struct{}
}

// Used as the value in the items map; the key represents the set member.
var exists = struct{}{}

// Initialize the internal data structure
func (s *set) init() *set {
	s.items = make(map[interface{}]struct{})
	return s
}

// Add one or more items to the set.
func (s *set) add(items ...interface{}) {
	for _, item := range items {
		s.items[item] = exists
	}
}

// Remove one or more items from the set.
func (s *set) remove(items ...interface{}) {
	for _, item := range items {
		delete(s.items, item)
	}
}

// Check for the existence of the item in the set.
func (s *set) contains(item interface{}) bool {
	_, contains := s.items[item]
	return contains
}

//===========================================================================
// Non-Thread-Safe Set
//===========================================================================

// SetNonTS defines a non-thread safe set data structure.
type SetNonTS struct {
	set
}

// NewNonTS creates and initializes a new non-thread-safe Set. It accepts a
// variable number of arguments to populate the initial set with.
func NewNonTS(items ...interface{}) *SetNonTS {
	s := new(SetNonTS)
	s.init()
	s.add(items...)
	return s
}

// Add one or more items to the set. If no items, silently exit.
func (s *SetNonTS) Add(items ...interface{}) {
	if len(items) == 0 {
		return
	}

	s.add(items...)
}

// Remove one or more items from the set. If no items, silently exit.
func (s *SetNonTS) Remove(items ...interface{}) {
	if len(items) == 0 {
		return
	}

	s.remove(items...)
}

// Contains returns true if the item is in the set.
func (s *SetNonTS) Contains(item interface{}) bool {
	return s.contains(item)
}

//===========================================================================
// Thread Safe Set
//===========================================================================

// Set defines a thread-safe set data structure.
type Set struct {
	set
	sync.RWMutex
}

// New creates and initializes a new thread-safe Set. It accepts a
// variable number of arguments to populate the initial set with.
func New(items ...interface{}) *Set {
	s := new(Set)
	s.Lock()
	defer s.Unlock()

	s.init()
	s.add(items...)
	return s
}

// Add one or more items to the set. If no items, silently exit.
func (s *Set) Add(items ...interface{}) {
	if len(items) == 0 {
		return
	}

	s.Lock()
	defer s.Unlock()
	s.add(items...)
}

// Remove one or more items from the set. If no items, silently exit.
func (s *Set) Remove(items ...interface{}) {
	if len(items) == 0 {
		return
	}

	s.Lock()
	defer s.Unlock()
	s.remove(items...)
}

// Contains returns true if the item is in the set.
func (s *Set) Contains(item interface{}) bool {
	s.RLock()
	defer s.RUnlock()
	return s.contains(item)
}
