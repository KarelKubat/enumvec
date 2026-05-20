package enumvec

import (
	"testing"
)

func TestEnumVec_Basic(t *testing.T) {
	// Example: max=15 (4 bits per value)
	ev := New(15)
	if err := ev.Set(10, 0); err != nil {
		t.Errorf("Set failed: %v", err)
	}
	if v := ev.Get(0); v != 10 {
		t.Errorf("Expected 10, got %d", v)
	}
	if err := ev.Set(16, 0); err == nil {
		t.Error("Expected error for value > 15, got nil")
	}
}

func TestEnumVec_Bits(t *testing.T) {
	// Example: max=1 (1 bit per value)
	ev := New(1)
	for i := uint64(0); i < 100; i++ {
		val := i % 2
		if err := ev.Set(val, i); err != nil {
			t.Fatalf("Set(%d, %d) failed: %v", val, i, err)
		}
	}
	for i := uint64(0); i < 100; i++ {
		expected := i % 2
		if v := ev.Get(i); v != expected {
			t.Errorf("Get(%d) expected %d, got %d", i, expected, v)
		}
	}
}

func TestEnumVec_MultiBit(t *testing.T) {
	// Example: max=3 (2 bits per value, 32 values per uint64)
	ev := New(3)
	if err := ev.Set(2, 35); err != nil {
		t.Fatalf("Set(2, 35) failed: %v", err)
	}

	if v := ev.Get(35); v != 2 {
		t.Errorf("Get(35) expected 2, got %d", v)
	}

	// Verify it's in the second word
	if len(ev.vec) < 2 {
		t.Errorf("Expected at least 2 words in vec, got %d", len(ev.vec))
	}
}

func TestEnumVec_Max64(t *testing.T) {
	// Edge case: max=^uint64(0) (64 bits per value)
	ev := New(^uint64(0))
	val := uint64(0xDEADBEEFCAFEBABE)
	if err := ev.Set(val, 5); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	if v := ev.Get(5); v != val {
		t.Errorf("Expected %x, got %x", val, v)
	}
}

func TestEnumVec_Overwrite(t *testing.T) {
	ev := New(7) // 3 bits
	ev.Set(5, 10)
	if v := ev.Get(10); v != 5 {
		t.Errorf("Expected 5, got %d", v)
	}
	ev.Set(2, 10)
	if v := ev.Get(10); v != 2 {
		t.Errorf("Expected 2, got %d", v)
	}
}

func TestEnumVec_LargeIndex(t *testing.T) {
	// Use a large index to ensure uint64 works correctly.
	ev := New(15)
	index := uint64(1000000)
	value := uint64(13)

	if err := ev.Set(value, index); err != nil {
		t.Fatalf("Set(%d, %d) failed: %v", value, index, err)
	}

	if v := ev.Get(index); v != value {
		t.Errorf("Get(%d) expected %d, got %d", index, value, v)
	}

	// Check another large index
	index2 := uint64(2000000)
	value2 := uint64(7)
	if err := ev.Set(value2, index2); err != nil {
		t.Fatalf("Set(%d, %d) failed: %v", value2, index2, err)
	}
	if v := ev.Get(index2); v != value2 {
		t.Errorf("Get(%d) expected %d, got %d", index2, value2, v)
	}

	// Ensure first one is still there
	if v := ev.Get(index); v != value {
		t.Errorf("Get(%d) after second set expected %d, got %d", index, value, v)
	}
}

func TestEnumVec_Size(t *testing.T) {
	ev := New(1) // 1 bit per value, 64 values per word
	if s := ev.Size(); s != 0 {
		t.Errorf("Expected initial size 0, got %d", s)
	}

	ev.Set(1, 0) // Should allocate 1 word (8 bytes)
	if s := ev.Size(); s != 8 {
		t.Errorf("Expected size 8 after first set, got %d", s)
	}

	ev.Set(1, 64) // Should allocate 2 words (16 bytes)
	if s := ev.Size(); s != 16 {
		t.Errorf("Expected size 16 after setting index 64, got %d", s)
	}
}
