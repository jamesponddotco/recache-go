package recache_test

import (
	"testing"

	"git.sr.ht/~jamesponddotco/recache-go"
)

const (
	_testPattern        string = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	_TestPatternPOSIX   string = `^[[:alnum:]._%+-]+@[[:alnum:].-]+\.[[:alpha:]]{2,}$`
	_testPatternInvalid string = `[`
)

func TestCompile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		givePattern string
		giveFlag    recache.Flag
		err         bool
	}{
		{
			name:        "Default",
			givePattern: _testPattern,
			giveFlag:    recache.Flag(0),
		},
		{
			name:        "POSIX",
			givePattern: _TestPatternPOSIX,
			giveFlag:    recache.FlagPOSIX,
		},
		{
			name:        "Must",
			givePattern: _testPattern,
			giveFlag:    recache.FlagMust,
		},
		{
			name:        "POSIX and Must",
			givePattern: _TestPatternPOSIX,
			giveFlag:    recache.FlagMustPOSIX,
		},
		{
			name:        "Invalid default pattern",
			givePattern: _testPatternInvalid,
			giveFlag:    recache.Flag(0),
			err:         true,
		},
		{
			name:        "Invalid POSIX pattern",
			givePattern: _testPatternInvalid,
			giveFlag:    recache.FlagPOSIX,
			err:         true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := recache.Compile(tt.givePattern, tt.giveFlag)
			if tt.err && err == nil {
				t.Errorf("expected an error, got nil")
			} else if !tt.err && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}

	// Test the Must flag with an invalid pattern, which should cause a panic.
	t.Run("Invalid pattern with Must flag", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("MustCompile should have panicked on an invalid pattern")
			}
		}()

		recache.Compile("[", recache.FlagMust) //nolint:errcheck // no error check, this should panic
	})
}

func TestKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		givePattern string
		giveFlag    recache.Flag
		wantKey     string
	}{
		{
			name:        "Default",
			givePattern: _testPattern,
			giveFlag:    recache.Flag(0),
			wantKey:     "6598b79548be6c59",
		},
		{
			name:        "POSIX",
			givePattern: _TestPatternPOSIX,
			giveFlag:    recache.FlagPOSIX,
			wantKey:     "6d579f15e31b5927",
		},
		{
			name:        "Must",
			givePattern: _testPattern,
			giveFlag:    recache.FlagMust,
			wantKey:     "3d4b9eb26fdfe02d",
		},
		{
			name:        "POSIX and Must",
			givePattern: _TestPatternPOSIX,
			giveFlag:    recache.FlagMustPOSIX,
			wantKey:     "fae3eb212ccb70b0",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotKey := recache.Key(tt.givePattern, tt.giveFlag)
			if gotKey != tt.wantKey {
				t.Errorf("Key() = %v, want %v", gotKey, tt.wantKey)
			}
		})
	}
}
