// Package enumvec stores small uint-like values (typically enums) in a memory-optimized way.
package enumvec

import (
	"fmt"
	"math/bits"
)

// EnumVec stores small integer values packed into a uint64 slice.
type EnumVec struct {
	vec              []uint64
	max              uint64
	bitsPerVal       int
	valsPerWord      int
	mask             uint64
	isPowerOfTwo     bool
	valsPerWordShift uint
	valsPerWordMask  uint64
}

// New initializes and returns a *EnumVec.
// max defines the maximum value that can be stored (0 to max).
func New(max uint64) *EnumVec {
	return NewWithCapacity(max, 0)
}

// NewWithCapacity initializes and returns an *EnumVec with pre-allocated storage for initialSize values.
// max defines the maximum value that can be stored (0 to max).
func NewWithCapacity(max uint64, initialSize uint64) *EnumVec {
	if max == 0 {
		// Even if max is 0, we need at least 1 bit to represent the value 0.
		// If max is 0, we only store 0, which takes 1 bit (or technically 0, but 1 is safer for logic).
		return &EnumVec{
			max:              0,
			bitsPerVal:       1,
			valsPerWord:      64,
			mask:             0,
			isPowerOfTwo:     true,
			valsPerWordShift: 6,
			valsPerWordMask:  63,
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

	isPowerOfTwo := (valsPerWord & (valsPerWord - 1)) == 0
	var shift uint
	var vpwMask uint64
	if isPowerOfTwo {
		shift = uint(bits.TrailingZeros64(uint64(valsPerWord)))
		vpwMask = uint64(valsPerWord - 1)
	}

	var vec []uint64
	if initialSize > 0 {
		wordCount := (initialSize + uint64(valsPerWord) - 1) / uint64(valsPerWord)
		vec = make([]uint64, wordCount)
	}

	return &EnumVec{
		vec:              vec,
		max:              max,
		bitsPerVal:       bpv,
		valsPerWord:      valsPerWord,
		mask:             mask,
		isPowerOfTwo:     isPowerOfTwo,
		valsPerWordShift: shift,
		valsPerWordMask:  vpwMask,
	}
}

// Set stores value at the given index.
// Returns an error if value exceeds the max value provided in New().
// Dynamically resizes the internal slice if needed.
func (ev *EnumVec) Set(value uint64, index uint64) error {
	if value > ev.max {
		return fmt.Errorf("value %d exceeds max %d", value, ev.max)
	}
	if ev.max == 0 {
		// When max is 0, any valid value is 0. 0 is default and takes no memory, so short-circuit.
		return nil
	}

	var wordIdx int
	var valOffset uint
	if ev.isPowerOfTwo {
		wordIdx = int(index >> ev.valsPerWordShift)
		valOffset = uint(index & ev.valsPerWordMask) * uint(ev.bitsPerVal)
	} else {
		wordIdx = int(index / uint64(ev.valsPerWord))
		valOffset = uint(index % uint64(ev.valsPerWord)) * uint(ev.bitsPerVal)
	}

	// Ensure slice is large enough
	if wordIdx >= len(ev.vec) {
		newLen := wordIdx + 1
		if newLen <= cap(ev.vec) {
			ev.vec = ev.vec[:newLen]
		} else {
			newCap := cap(ev.vec) * 2
			if newCap < newLen {
				newCap = newLen
			}
			newVec := make([]uint64, newLen, newCap)
			copy(newVec, ev.vec)
			ev.vec = newVec
		}
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
	if ev.max == 0 {
		return 0
	}

	var wordIdx int
	var valOffset uint
	if ev.isPowerOfTwo {
		wordIdx = int(index >> ev.valsPerWordShift)
		valOffset = uint(index & ev.valsPerWordMask) * uint(ev.bitsPerVal)
	} else {
		wordIdx = int(index / uint64(ev.valsPerWord))
		valOffset = uint(index % uint64(ev.valsPerWord)) * uint(ev.bitsPerVal)
	}

	if wordIdx >= len(ev.vec) {
		return 0
	}

	return (ev.vec[wordIdx] >> valOffset) & ev.mask
}

// Size returns the number of bytes used by the internal vector.
func (ev *EnumVec) Size() uint64 {
	return uint64(len(ev.vec)) * 8
}
