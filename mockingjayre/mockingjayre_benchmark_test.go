package mockingjayre_test

import (
	"context"
	"regexp"
	"testing"

	"git.sr.ht/~jamesponddotco/recache-go"
	"git.sr.ht/~jamesponddotco/recache-go/mockingjayre"
)

const (
	testPattern = `p([a-z]+)ch`
)

func BenchmarkMockingjayCache(b *testing.B) {
	cache := mockingjayre.New(recache.DefaultCapacity)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		re, err := cache.Get(context.Background(), testPattern, recache.FlagMust)
		if err != nil {
			b.Fatal(err)
		}

		re.MatchString("peach")

		reAgain, err := cache.Get(context.Background(), testPattern, recache.FlagMust)
		if err != nil {
			b.Fatal(err)
		}

		reAgain.MatchString("peach")
	}
}

func BenchmarkRegexpCompile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		re := regexp.MustCompile(testPattern)
		reAgain := regexp.MustCompile(testPattern)

		re.MatchString("peach")
		reAgain.MatchString("peach")
	}
}
