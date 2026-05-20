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

var store [1000000]uint8

store[0] = First;
store[1] = Second;
// ... and so on
```

But you only need 2 bits to store each value! Per element, the above array `store` won't use 6 out of each 8 bits. That's a waste of 75%.

Or you can do this:

- Create an enumvec-based store: `store := enumvec.New(Third)`. The argument says what the max of each value can be. (If you have already a hint to the required capacity, you can use: `enumvec.NewWithCapacity(Third, 1000000)`. That saves some reallocations.)
- Store values: `store.Set(Second, 500000)` which stores at index 500000. Never mind the required storage to reach that half million index, it's automatically adjusted.
- Similarly, you can `store.Get(500000)` to retrieve a value.
- If you are interested how many bytes are actually used in the storage, ask `store.Size()`.

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
- You have to know in advance the maximum value to store. If you create a storage for say up to 6 as in `enumvec.New(6)`, then you can store anything from 0 to 6 inclusive. Trying to store value 7 won't work, but it will return an error.

As to the speed/memory trade-off, for maximum values up to 15 the `enumvec` is slower than a regular array. But as the max value grows, `enumvec` gets faster and faster in relation to the plain array. At values of 256 (so 8 bits per value) there is no difference anymore, as it uses 8 bits per value as well. In case that you use `enumvec` with a power of two size, it uses bit-shifting and masking for index calculations which is slightly faster.

> Note: The above is true in case that the values are stored densely. If you store scattered values, the plain array is faster in both speed and memory.	

## Benchmarks

The following is the output of benchmarking the included file `bench/enumvec_bench.go`.

```text
==================================================
    EnumVec vs. Raw []uint8 Benchmark
    Elements: 1,000,000 | Values: 1, 2, 3, 4, 5
==================================================
Implementation       Time Taken      Memory Used     Checksum (Sum) 
------------------------------------------------------------------
Raw []uint8          931.70µs        976.6 KB        3,000,000        
EnumVec (3-bit)      6.81ms          372.0 KB        3,000,000        
------------------------------------------------------------------
Memory Saved: 604.5 KB (61.90% savings)
Speed Overhead: EnumVec is 7.31x slower than raw slice (due to bitwise operations)
==================================================
```

So yes, `enumvec` is slower than plain array, but it saves a lot of memory in comparison.
