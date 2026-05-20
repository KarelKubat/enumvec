// This is just a test- and demo program.
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
