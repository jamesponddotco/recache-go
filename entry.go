package recache

import (
	"regexp"
	"sync/atomic"

	"git.sr.ht/~jamesponddotco/xstd-go/xerrors"
)

// ErrInvalidEntry is returned when the entry in the cache is invalid.
const ErrInvalidEntry xerrors.Error = "invalid entry"

// Entry represents an item in the cache.
type Entry struct {
	regex     *regexp.Regexp
	pattern   string
	key       string
	frequency atomic.Uint64
}

// NewEntry creates a new entry in the cache.
func NewEntry(key, pattern string, regex *regexp.Regexp) *Entry {
	if pattern == "" || regex == nil {
		return nil
	}

	return &Entry{
		regex:     regex,
		pattern:   pattern,
		key:       key,
		frequency: atomic.Uint64{},
	}
}

// Load returns the compiled regex and pattern, and increments the frequency of
// the entry by one.
func (e *Entry) Load() (*regexp.Regexp, string, error) {
	e.frequency.Add(1)

	return e.regex, e.pattern, nil
}

// Pattern returns the entry's pattern.
func (e *Entry) Pattern() string {
	return e.pattern
}

// Key returns the entry's cache key.
func (e *Entry) Key() string {
	return e.key
}

// Frequency returns the number of times the entry has been loaded.
func (e *Entry) Frequency() uint64 {
	return e.frequency.Load()
}
