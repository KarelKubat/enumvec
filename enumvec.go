// Package enumvec stores small uint-like values (typically enums) in a memory-optimized way.
package enumvec

import (
	"fmt"
	"math/bits"
)

// EnumVec stores small integer values packed into a uint64 slice.
type EnumVec struct {
	vec         []uint64
	max         uint64
	bitsPerVal  int
	valsPerWord int
	mask        uint64
}

// New initializes and returns a *EnumVec.
// max defines the maximum value that can be stored (0 to max).
func New(max uint64) *EnumVec {
	if max == 0 {
		// Even if max is 0, we need at least 1 bit to represent the value 0.
		// If max is 0, we only store 0, which takes 1 bit (or technically 0, but 1 is safer for logic).
		return &EnumVec{
			max:         0,
			bitsPerVal:  1,
			valsPerWord: 64,
			mask:        0,
		}
	}

	// bits.Len64(max) gives the number of bits needed to represent max.
	// e.g., max=1 (binary 1) -> Len=1
	// max=3 (binary 11) -> Len=2
	// max=15 (binary 1111) -> Len=4
	bpv := bits.Len64(max)
	if bpv > 64 {
		bpv = 64
	}

	valsPerWord := 64 / bpv
	mask := uint64((1 << bpv) - 1)
	if bpv == 64 {
		mask = ^uint64(0)
	}

	return &EnumVec{
		max:         max,
		bitsPerVal:  bpv,
		valsPerWord: valsPerWord,
		mask:        mask,
	}
}

// Set stores value at the given index.
// Returns an error if value exceeds the max value provided in New().
// Dynamically resizes the internal slice if needed.
func (ev *EnumVec) Set(value uint64, index uint64) error {
	if value > ev.max {
		return fmt.Errorf("value %d exceeds max %d", value, ev.max)
	}

	wordIdx := int(index / uint64(ev.valsPerWord))
	valOffset := uint(index%uint64(ev.valsPerWord)) * uint(ev.bitsPerVal)

	// Ensure slice is large enough
	if wordIdx >= len(ev.vec) {
		newVec := make([]uint64, wordIdx+1)
		copy(newVec, ev.vec)
		ev.vec = newVec
	}

	// Clear the existing bits at that position
	ev.vec[wordIdx] &= ^(ev.mask << valOffset)
	// Set the new bits
	ev.vec[wordIdx] |= (value & ev.mask) << valOffset

	return nil
}

// Get returns the value stored at the given index.
// If index is out of bounds, it returns 0 (default value).
func (ev *EnumVec) Get(index uint64) uint64 {
	wordIdx := int(index / uint64(ev.valsPerWord))
	if wordIdx >= len(ev.vec) {
		return 0
	}

	valOffset := uint(index%uint64(ev.valsPerWord)) * uint(ev.bitsPerVal)
	return (ev.vec[wordIdx] >> valOffset) & ev.mask
}

// Size returns the number of bytes used by the internal vector.
func (ev *EnumVec) Size() uint64 {
	return uint64(len(ev.vec)) * 8
}

