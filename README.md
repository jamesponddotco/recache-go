# recache

Package `recache` is a lightweight caching library for [Go's standard regular
expression package](https://godocs.io/regexp) that offers improved performance
by avoiding recompilation of global regular expression variables and by caching
regular expressions.

## Features

- Stable cache interface.
- Simple and easy-to-use API.
- Thread-safe caching of compiled regular expressions.
- Lazy compilation of regular expressions.
- Minimal memory allocations.


### `recache.Cache` implementations

The `recache` package itself only provides a cache interface and some utility
functions for users who wish to implement that interface. You can either use an
implementation created by someone else or write your own.

**Implementations**

- [`lrure`](https://git.sr.ht/~jamesponddotco/recache-go/tree/trunk/item/lrure)
  provides a thread-safe in-memory cache using the [least recently used
  (LRU)](https://en.wikipedia.org/wiki/Cache_replacement_policies#Least_recently_used_(LRU))
  cache replacement policy.

If wrote a `recache.Cache` implementation and wish it to be linked here,
[please send a patch](https://git.sr.ht/~jamesponddotco/recache-go#resources).

## Installation

To install `recache` alone, run:

```sh
go get git.sr.ht/~jamesponddotco/recache-go
```

## Contributing

Anyone can help make recache better. Check out [the contribution
guidelines](https://git.sr.ht/~jamesponddotco/recache-go/tree/master/item/CONTRIBUTING.md)
for more information.

## Resources

The following resources are available:

- [Package documentation](https://godocs.io/git.sr.ht/~jamesponddotco/recache-go).
- [Support and general discussions](https://lists.sr.ht/~jamesponddotco/recache-discuss).
- [Patches and development related questions](https://lists.sr.ht/~jamesponddotco/recache-devel).
- [Instructions on how to prepare patches](https://git-send-email.io/).
- [Feature requests and bug reports](https://todo.sr.ht/~jamesponddotco/recache).

---

Released under the [MIT License](LICENSE.md).
