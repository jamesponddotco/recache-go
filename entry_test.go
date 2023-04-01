package recache_test

import (
	"regexp"
	"testing"

	"git.sr.ht/~jamesponddotco/recache-go"
)

func TestNewEntry(t *testing.T) {
	t.Parallel()

	regex := regexp.MustCompile(_testPattern)

	tests := []struct {
		name        string
		giveKey     string
		givePattern string
		giveRegex   *regexp.Regexp
		wantNil     bool
	}{
		{
			name:        "Valid Entry",
			giveKey:     "test_key",
			givePattern: _testPattern,
			giveRegex:   regex,
		},
		{
			name:        "Empty Pattern",
			giveKey:     "test_key",
			givePattern: "",
			giveRegex:   regex,
			wantNil:     true,
		},
		{
			name:        "Nil Regexp",
			giveKey:     "test_key",
			givePattern: _testPattern,
			giveRegex:   nil,
			wantNil:     true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := recache.NewEntry(tt.giveKey, tt.givePattern, tt.giveRegex)
			if tt.wantNil && got != nil {
				t.Errorf("NewEntry() should return nil, got %v", got)
			} else if !tt.wantNil && got == nil {
				t.Errorf("NewEntry() should not return nil")
			}
		})
	}
}

func TestEntryMethods(t *testing.T) {
	t.Parallel()

	var (
		regex = regexp.MustCompile(_testPattern)
		entry = recache.NewEntry("test_key", _testPattern, regex)
	)

	t.Run("Load", func(t *testing.T) {
		t.Parallel()

		loadedRegex, loadedPattern, err := entry.Load()
		if err != nil {
			t.Errorf("Load() returned an error: %v", err)
		}

		if loadedRegex != regex {
			t.Errorf("Load() returned incorrect regex: expected %v, got %v", regex, loadedRegex)
		}

		if loadedPattern != _testPattern {
			t.Errorf("Load() returned incorrect pattern: expected %v, got %v", _testPattern, loadedPattern)
		}
	})

	t.Run("Key", func(t *testing.T) {
		t.Parallel()

		key := entry.Key()
		if key != "test_key" {
			t.Errorf("Key() returned incorrect key: expected %v, got %v", "test_key", key)
		}
	})

	t.Run("Pattern", func(t *testing.T) {
		t.Parallel()

		pattern := entry.Pattern()
		if pattern != _testPattern {
			t.Errorf("Pattern() returned incorrect pattern: expected %v, got %v", _testPattern, pattern)
		}
	})

	t.Run("Frequency", func(t *testing.T) {
		t.Parallel()

		for i := 0; i < 3; i++ {
			_, _, err := entry.Load()
			if err != nil {
				t.Fatalf("Load() returned an error: %v", err)
			}
		}

		frequency := entry.Frequency()
		if frequency != 4 { // Including the first Load() call in the "Load" subtest
			t.Errorf("Frequency() returned incorrect value: expected %v, got %v", 4, frequency)
		}
	})
}
