package lrure_test

import (
	"context"
	"testing"

	"git.sr.ht/~jamesponddotco/recache-go"
	"git.sr.ht/~jamesponddotco/recache-go/lrure"
)

// TestNew tests the New function from the lrure package.
func TestNew(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		give int
		want int
	}{
		{
			name: "Default capacity",
			give: recache.DefaultCapacity,
			want: recache.DefaultCapacity,
		},
		{
			name: "Custom capacity",
			give: 100,
			want: 100,
		},
		{
			name: "Negative capacity",
			give: -1,
			want: recache.DefaultCapacity,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			lru := lrure.New(tt.give)

			if lru.Capacity() != tt.want {
				t.Errorf("Capacity() = %d, want %d", lru.Capacity(), tt.want)
			}
		})
	}
}

func TestCache(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name             string
		capacity         int
		patterns         []string
		flags            []recache.Flag
		expectedSize     int
		expectedCapacity int
	}{
		{
			name:             "basic test",
			capacity:         3,
			patterns:         []string{`^hello`, `^world`, `^\d+$`},
			flags:            []recache.Flag{recache.DefaultFlag, recache.DefaultFlag, recache.DefaultFlag},
			expectedSize:     3,
			expectedCapacity: 3,
		},
		{
			name:             "LRU eviction",
			capacity:         2,
			patterns:         []string{`^hello`, `^world`, `^\d+$`},
			flags:            []recache.Flag{recache.DefaultFlag, recache.DefaultFlag, recache.DefaultFlag},
			expectedSize:     2,
			expectedCapacity: 2,
		},
		{
			name:             "with flags",
			capacity:         3,
			patterns:         []string{`^hello`, `^world`},
			flags:            []recache.Flag{recache.FlagPOSIX, recache.FlagMust},
			expectedSize:     2,
			expectedCapacity: 3,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cache := lrure.New(tt.capacity)

			for i, pattern := range tt.patterns {
				flag := tt.flags[i%len(tt.flags)]
				_, err := cache.Get(ctx, pattern, flag)
				if err != nil {
					t.Errorf("Cache.Get() error = %v, wantErr = false", err)
				}
			}

			size := cache.Size()
			if size != tt.expectedSize {
				t.Errorf("Cache.Size() = %v, want = %v", size, tt.expectedSize)
			}

			capacity := cache.Capacity()
			if capacity != tt.expectedCapacity {
				t.Errorf("Cache.Capacity() = %v, want = %v", capacity, tt.expectedCapacity)
			}

			// Test clearing the cache
			cache.Clear()
			size = cache.Size()
			if size != 0 {
				t.Errorf("Cache.Size() after Clear() = %v, want = 0", size)
			}
		})
	}
}

func TestSetCapacity(t *testing.T) {
	tests := []struct {
		name             string
		initialCapacity  int
		newCapacity      int
		expectedError    error
		expectedSize     int
		expectedCapacity int
	}{
		{
			name:             "set capacity to a higher value",
			initialCapacity:  2,
			newCapacity:      4,
			expectedError:    nil,
			expectedSize:     2,
			expectedCapacity: 4,
		},
		{
			name:             "set capacity to a lower value",
			initialCapacity:  3,
			newCapacity:      2,
			expectedError:    nil,
			expectedSize:     2,
			expectedCapacity: 2,
		},
		{
			name:             "set capacity to zero",
			initialCapacity:  3,
			newCapacity:      0,
			expectedError:    recache.ErrInvalidCapacity,
			expectedSize:     2,
			expectedCapacity: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := lrure.New(tt.initialCapacity)
			_, err := cache.Get(context.Background(), "a", recache.DefaultFlag)
			if err != nil {
				t.Fatalf("Cache.Get() error = %v, wantErr = false", err)
			}
			_, err = cache.Get(context.Background(), "b", recache.DefaultFlag)
			if err != nil {
				t.Fatalf("Cache.Get() error = %v, wantErr = false", err)
			}

			err = cache.SetCapacity(tt.newCapacity)

			if !(err == nil && tt.expectedError == nil) && !(err != nil && tt.expectedError != nil && err.Error() == tt.expectedError.Error()) {
				t.Errorf("SetCapacity() error = %v, want = %v", err, tt.expectedError)
			}

			if cache.Size() != tt.expectedSize {
				t.Errorf("Cache.Size() after SetCapacity() = %d, want = %d", cache.Size(), tt.expectedSize)
			}

			if cache.Capacity() != tt.expectedCapacity {
				t.Errorf("Cache.Capacity() after SetCapacity() = %d, want = %d", cache.Capacity(), tt.expectedCapacity)
			}
		})
	}
}
