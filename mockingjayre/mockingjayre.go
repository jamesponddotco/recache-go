// Package mockingjayre implements a thread-safe cache for [Go's standard regex
// package] that complies with the [recache.Cache] interface. It uses
// [Mockingjay] as its cache replacement policy.
//
// [Go's standard regex package]: https://godocs.io/regexp
// [recache.Cache]: https://godocs.io/git.sr.ht/~jamesponddotco/recache-go#Cache
// [Mockingjay]: https://en.wikipedia.org/wiki/Cache_replacement_policies#Mockingjay
package mockingjayre

import (
	"context"
	"fmt"
	"regexp"
	"sync"

	"git.sr.ht/~jamesponddotco/recache-go"
)

// Cache is a thread-safe regex cache using the Mockingjay policy.
type Cache struct {
	cache    map[string]*recache.Entry
	capacity int
	mu       sync.RWMutex
}

// New returns a new Mockingjay cache with the given capacity.
//
// If capacity is less than 1, [recache.DefaultCapacity] is used instead.
//
// [recache.DefaultCapacity]: https://godocs.io/git.sr.ht/~jamesponddotco/recache-go#DefaultCapacity
func New(capacity int) *Cache {
	if capacity < 1 {
		capacity = recache.DefaultCapacity
	}

	return &Cache{
		cache:    make(map[string]*recache.Entry, capacity),
		capacity: capacity,
	}
}

// Get returns a compiled regular expression from the cache given a pattern and
// an optional flag. If the regular expression is not in the cache, it is
// compiled and added to it.
func (c *Cache) Get(_ context.Context, pattern string, flag recache.Flag) (*regexp.Regexp, error) {
	key := recache.Key(pattern, flag)

	c.mu.Lock()
	defer c.mu.Unlock()

	if entry, ok := c.cache[key]; ok {
		regex, _, err := entry.Load()
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		return regex, nil
	}

	regex, err := recache.Compile(pattern, flag)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if len(c.cache) >= c.capacity {
		c.evict()
	}

	newEntry := recache.NewEntry(key, pattern, regex)

	c.cache[key] = newEntry

	return regex, nil
}

// SetCapacity sets the maximum number of regular expressions that can be
// stored in the cache.
func (c *Cache) SetCapacity(capacity int) error {
	if capacity < 1 {
		return fmt.Errorf("%w", recache.ErrInvalidCapacity)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

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

	return len(c.cache)
}

// Clear removes all regular expressions from the cache.
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache = make(map[string]*recache.Entry, c.capacity)
}

// evict removes the least frequently used entry from the cache.
func (c *Cache) evict() {
	var (
		leastFrequentlyUsedKey   string
		leastFrequentlyUsedEntry *recache.Entry
	)

	for key, entry := range c.cache {
		if leastFrequentlyUsedEntry == nil || entry.Frequency() < leastFrequentlyUsedEntry.Frequency() {
			leastFrequentlyUsedKey = key
			leastFrequentlyUsedEntry = entry
		}
	}

	delete(c.cache, leastFrequentlyUsedKey)
}
