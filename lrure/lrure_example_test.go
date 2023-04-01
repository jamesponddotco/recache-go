package lrure_test

import (
	"context"
	"fmt"
	"log"

	"git.sr.ht/~jamesponddotco/recache-go"
	"git.sr.ht/~jamesponddotco/recache-go/lrure"
)

func ExampleCache_Get() {
	// Create a new Cache instance with the default cache capacity.
	cache := lrure.New(recache.DefaultCapacity)

	// Add the regular expression to the cache for the first time, which will
	// cause it to be compiled.
	regex, err := cache.Get(context.Background(), "p([a-z]+)ch", recache.DefaultFlag)
	if err != nil {
		log.Fatal(err)
	}

	// Match the string against the regular expression.
	fmt.Println(regex.MatchString("peach"))

	// Get the regular expression, which by now has been compiled and returns
	// super fast, without the need for recompilation.
	sameRegex, err := cache.Get(context.Background(), "p([a-z]+)ch", recache.DefaultFlag)
	if err != nil {
		log.Fatal(err)
	}

	// Match the string against the regular expression.
	fmt.Println(sameRegex.MatchString("peach"))

	// Output:
	// true
	// true
}
