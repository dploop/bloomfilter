# Bloom filter

A [Bloom filter](https://en.wikipedia.org/wiki/Bloom_filter) is a space-efficient
probabilistic data structure, conceived by Burton Howard Bloom in 1970, that is
used to test whether an element is a member of a set. False positive matches are
possible, but false negatives are not - in other words, a query returns either
"possibly in set" or "definitely not in set". Elements can be added to the set, but
not removed (though this can be addressed with the counting Bloom filter variant);
the more elements that are added to the set, the larger the probability of false positives.

## Installation

```bash
go get -u github.com/dploop/bloomfilter
```

## How to use

```golang
package main

import (
	"log"

	"github.com/dploop/bloomfilter"
)

func main() {
	bf, err := bloomfilter.NewWithEstimate(1000000, 0.03)
	if err != nil {
		log.Fatalf("failed to new bloom filter: %v", err)
	}
	bf.Add([]byte("foo"))
	if bf.Contains([]byte("bar")) {
		log.Printf("bar is possibly in the set")
	} else {
		log.Printf("bar is definitely not in the set")
	}
}
```
