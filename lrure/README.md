# lrure

Package `lrure` is a thread-safe in-memory LRU cache for [Go's standard regex
package](https://godocs.io/regexp) that complies with the [recache.Cache
interface](https://godocs.io/git.sr.ht/~jamesponddotco/recache-go#Cache).

## Installation

To install `lrure`, run:

```sh
go get git.sr.ht/~jamesponddotco/recache-go/lrure
```

## Usage

```go
package main

import (
	"context"
	"log"

	"git.sr.ht/~jamesponddotco/recache-go"
	"git.sr.ht/~jamesponddotco/recache-go/lrure"
)

// Global Cache instance with the default cache capacity.
var _reCache = lrure.New(recache.DefaultCapacity)

func main() {
	// Add the regular expression to the cache for the first time, which will
	// cause it to be compiled.
	regex, err := _reCache.Get(context.Background(), "p([a-z]+)ch", recache.DefaultFlag)
	if err != nil {
		log.Fatal(err)
	}

	// Match the string against the regular expression.
	log.Println(regex.MatchString("peach"))

	// Get the regular expression, which by now has been compiled and returns
	// super fast, without the need for recompilation.
	sameRegex, err := _reCache.Get(context.Background(), "p([a-z]+)ch", recache.DefaultFlag)
	if err != nil {
		log.Fatal(err)
	}

	// Match the string against the regular expression.
	log.Println(sameRegex.MatchString("peach"))
}
```

Refer to [the API documentation](https://godocs.io/git.sr.ht/~jamesponddotco/recache-go) for more information.
