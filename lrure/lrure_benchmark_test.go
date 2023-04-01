package lrure_test

import (
	"regexp"
	"testing"

	"git.sr.ht/~jamesponddotco/recache-go"
	"git.sr.ht/~jamesponddotco/recache-go/lrure"
)

const (
	testPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
)

func BenchmarkLRUCache(b *testing.B) {
	cache := lrure.New(recache.DefaultCapacity)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cache.Get(nil, testPattern, recache.Flag(0))
	}
}

func BenchmarkRegexpCompile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = regexp.Compile(testPattern)
	}
}
