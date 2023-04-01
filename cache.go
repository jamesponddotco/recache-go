package recache

import (
	"context"
	"regexp"
)

const (
	// DefaultFlag is the default flag used when compiling regular expressions.
	DefaultFlag Flag = 0

	// FlagPOSIX specifies that the regular expression should be restricted to
	// POSIX syntax.
	FlagPOSIX Flag = 1 << iota

	// FlagMust specifies that the regular expression should be compiled using
	// MustCompile.
	FlagMust

	// FlagMustPOSIX is like Must but restricts the regular expression to POSIX
	// syntax.
	FlagMustPOSIX = FlagMust | FlagPOSIX
)

// Flag controls the behavior of the Get method when compiling regular
// expressions.
type Flag int

// String returns a string representation of the flag.
func (f Flag) String() string {
	switch f {
	case FlagPOSIX:
		return "POSIX"
	case FlagMust:
		return "Must"
	case FlagMustPOSIX:
		return "MustPOSIX"
	default:
		return "Default"
	}
}

// Cache is a storage mechanism used to store and retrieve compiled regular
// expressions for improved performance.
type Cache interface {
	// Get returns a compiled regular expression from the cache given a pattern
	// and an optional flag.
	Get(ctx context.Context, pattern string, flag Flag) (*regexp.Regexp, error)

	// SetCapacity sets the maximum number of regular expressions that can be
	// stored in the cache.
	SetCapacity(capacity int) error

	// Capacity returns the maximum number of regular expressions that can be
	// stored in the cache.
	Capacity() int

	// Size returns the number of regular expressions currently stored in the
	// cache.
	Size() int

	// Clear removes all regular expressions from the cache.
	Clear()
}
