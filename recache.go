// Package recache provides a simple caching interface for Go's regular
// expressions package, [regexp].
//
// [regexp]: https://godocs.io/regexp
package recache

import (
	"fmt"
	"regexp"

	"git.sr.ht/~jamesponddotco/xstd-go/xerrors"
	"git.sr.ht/~jamesponddotco/xstd-go/xhash/xfnv"
)

const (
	_keySeparator string = ":"
	_keyPattern   string = "pattern:"
	_KeyFlag      string = "flag:"
)

const (
	// ErrInvalidCapacity is returned when the provided capacity is less than 1.
	//
	// This error is not used by the package itself, but is exported for use by
	// packages implementing the Cache interface.
	ErrInvalidCapacity xerrors.Error = "invalid capacity"

	// ErrUnexpectedType is returned when the cache encounters an unexpected
	// type in the list.
	//
	// This error is not used by the package itself, but is exported for use by
	// packages implementing the Cache interface.
	ErrUnexpectedType xerrors.Error = "unexpected type found in the list: expected *recache.Entry"

	// ErrNotFound is returned when a regular expression is not found in the
	// cache.
	//
	// This error is not used by the package itself, but is exported for use by
	// packages implementing the Cache interface.
	ErrNotFound xerrors.Error = "not found in the cache"
)

// DefaultCapacity is the default maximum number of regular expressions that
// can be stored in the cache.
//
// This constant is not used by the package itself, but is exported for use by
// packages implementing the Cache interface.
const DefaultCapacity int = 25

// Compile compiles the provided regular expression pattern taking the provided
// control flag into account.
//
// This function is not used by the package itself, but is exported for use by
// packages implementing the Cache interface.
func Compile(pattern string, flag Flag) (*regexp.Regexp, error) {
	switch flag {
	case FlagPOSIX:
		re, err := regexp.CompilePOSIX(pattern)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		return re, nil
	case FlagMust:
		return regexp.MustCompile(pattern), nil
	case FlagMustPOSIX:
		return regexp.MustCompilePOSIX(pattern), nil
	default:
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		return re, nil
	}
}

// Key generates a cache key for the provided regular expression pattern and
// control flag by concatenating the two and hashing the result.
//
// The generated key is of the form "pattern:PATTERN:flag:FLAG".
//
// This function is not used by the package itself, but is exported for use by
// packages implementing the Cache interface.
func Key(pattern string, flag Flag) string {
	return xfnv.String(_keyPattern + pattern + _keySeparator + _KeyFlag + flag.String())
}
