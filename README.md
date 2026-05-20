# enumvec

`enumvec` is a package for Go to store small `uint`-like values in a memory-optimized way. Such values are typically `enum`s, hence the name. The idea is that small values don't require the standard number of bits for storage (8, 16, and so on). To store small values, one can use less bits.

## Example

Imagine that you have `uint`-like values 0, 1 and 2:

```go
const (
	First uint8 = iota
	Second
	Third
)

var store uint8[1000000];

store[0] = First;
store[1] = Second;
// ... and so on
```

But you only need 2 bits to store each value! Per element, the above array `store` won't leaves 6 out of each 8 bits unused. 

Or you can do this:

```go
package main

import (
	"fmt"

	"github.com/KarelKubat/enumvec"
)

const (
	First  uint64 = iota // enumvec works with uint64
	Second               // You can use any type you want, but then typecast below.
	Third
)

func main() {
	store := enumvec.New(Third)

	var i uint64
	for i = 0; i < 1000000; i += 3 {
		store.Set(First, i)
		store.Set(Second, i+1)
		store.Set(Third, i+2)
	}
	for i = 0; i < 20; i++ {
		fmt.Println("At index", i, "is value", store.Get(i))
	}

	fmt.Println("We used", store.Size(), "bytes to store 1.000.000 values")
}
```

## Limitations

- `enumvec` relies on bit operations (masking, shifting) to set and get values. This is inherently slower than plain indexing. As written above, `enumvec` is optimized to be memory-saving.
- All `enumvec`'s types are `uint64`: values, indexes, sizes. Typecast to whichever type you wish accordingly.
