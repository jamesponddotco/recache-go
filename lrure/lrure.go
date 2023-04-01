// Package lrure implements a thread-safe LRU cache for [Go's standard regex
// package] that complies with the [recache.Cache] interface.
//
// [Go's standard regex package]: https://godocs.io/regexp
// [recache.Cache]: https://godocs.io/git.sr.ht/~jamesponddotco/recache-go#Cache
package lrure

import (
	"container/list"
	"context"
	"fmt"
	"regexp"
	"sync"

	"git.sr.ht/~jamesponddotco/recache-go"
)

// Cache is a thread-safe LRU cache for Go's standard regex package.
type Cache struct {
	// cache is a map of the cache's keys to the elements that hold the values.
	cache map[string]*list.Element

	// list is a doubly-linked list that holds the cache's elements in order of
	// most recently used to least recently used.
	list *list.List

	// capacity is the maximum number of items the cache can hold.
	capacity int

	// mu is a mutex that protects access to the cache.
	mu sync.RWMutex
}

// Compile-time check to ensure Cache implements the recache.Cache interface.
var _ recache.Cache = (*Cache)(nil)

// New returns a new LRU cache with the given capacity.
//
// If capacity is less than 1, [recache.DefaultCapacity] is used instead.
//
// [recache.DefaultCapacity]: https://godocs.io/git.sr.ht/~jamesponddotco/recache-go#DefaultCapacity
func New(capacity int) *Cache {
	if capacity < 1 {
		capacity = recache.DefaultCapacity
	}

	return &Cache{
		cache:    make(map[string]*list.Element, capacity),
		list:     list.New().Init(),
		capacity: capacity,
	}
}

// Get returns a compiled regular expression from the cache given a pattern and
// an optional flag.
//
// If the regular expression is not in the cache, it is compiled and added to
// it.
func (c *Cache) Get(_ context.Context, pattern string, flag recache.Flag) (*regexp.Regexp, error) {
	key := recache.Key(pattern, flag)

	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.cache[key]; ok {
		entry, ok := elem.Value.(*recache.Entry)
		if !ok {
			return nil, fmt.Errorf("%w", recache.ErrUnexpectedType)
		}

		c.list.MoveToFront(elem)

		re, _, err := entry.Load()
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		return re, nil
	}

	regex, err := recache.Compile(pattern, flag)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	newEntry := recache.NewEntry(key, pattern, regex)

	c.cache[key] = c.list.PushFront(newEntry)

	if c.list.Len() > c.capacity {
		elem := c.list.Back()

		entry, ok := elem.Value.(*recache.Entry)
		if !ok {
			return nil, fmt.Errorf("%w", recache.ErrUnexpectedType)
		}

		c.list.Remove(elem)

		key := recache.Key(entry.Pattern(), recache.Flag(0))

		delete(c.cache, key)
	}

	return regex, nil
}

// SetCapacity sets the maximum number of regular expressions that can be
// stored in the cache.
func (c *Cache) SetCapacity(capacity int) error {
	if capacity < 1 {
		return recache.ErrInvalidCapacity
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	for c.list.Len() > capacity {
		elem := c.list.Back()

		entry, ok := elem.Value.(*recache.Entry)
		if !ok {
			return fmt.Errorf("%w", recache.ErrUnexpectedType)
		}

		c.list.Remove(elem)

		key := recache.Key(entry.Pattern(), recache.Flag(0))

		delete(c.cache, key)
	}

	c.capacity = capacity

	return nil
}

// Capacity returns the maximum number of regular expressions that can be
// stored in the cache.
func (c *Cache) Capacity() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.capacity
}

// Size returns the number of regular expressions currently stored in the
// cache.
func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.list.Len()
}

// Clear removes all regular expressions from the cache.
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.list.Init()
	c.cache = make(map[string]*list.Element, c.capacity)
}
