// Package main is a benchmarking tool to compare a raw uint8 slice with EnumVec.
package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/KarelKubat/enumvec"
)

func main() {
	const N = 1000000

	fmt.Println("==================================================")
	fmt.Println("    EnumVec vs. Raw []uint8 Benchmark")
	fmt.Println("    Elements: 1,000,000 | Values: 1, 2, 3, 4, 5")
	fmt.Println("==================================================")

	// Force garbage collection to ensure a clean start
	runtime.GC()

	// --- 1. Raw []uint8 Slice Benchmark ---
	startU8 := time.Now()
	u8Slice := make([]uint8, N)
	for i := 0; i < N; i++ {
		val := uint8((i % 5) + 1) // Store values 1, 2, 3, 4, 5
		u8Slice[i] = val
	}
	var sumU8 uint64
	for i := 0; i < N; i++ {
		sumU8 += uint64(u8Slice[i])
	}
	durationU8 := time.Since(startU8)
	sizeU8 := N // 1,000,000 bytes

	// Force garbage collection
	runtime.GC()

	// --- 2. EnumVec Benchmark ---
	startEV := time.Now()
	// NewWithCapacity pre-allocates to avoid dynamic reallocation overhead
	ev := enumvec.NewWithCapacity(5, N)
	for i := 0; i < N; i++ {
		val := uint64((i % 5) + 1) // Store values 1, 2, 3, 4, 5
		_ = ev.Set(val, uint64(i))
	}
	var sumEV uint64
	for i := 0; i < N; i++ {
		sumEV += ev.Get(uint64(i))
	}
	durationEV := time.Since(startEV)
	sizeEV := ev.Size()

	// --- 3. Results Output ---
	fmt.Printf("%-20s %-15s %-15s %-15s\n", "Implementation", "Time Taken", "Memory Used", "Checksum (Sum)")
	fmt.Println("------------------------------------------------------------------")
	fmt.Printf("%-20s %-15v %-15s %-15d\n", "Raw []uint8", durationU8, formatBytes(uint64(sizeU8)), sumU8)
	fmt.Printf("%-20s %-15v %-15s %-15d\n", "EnumVec (3-bit)", durationEV, formatBytes(sizeEV), sumEV)
	fmt.Println("------------------------------------------------------------------")

	savingsBytes := int64(sizeU8) - int64(sizeEV)
	savingsPct := (float64(savingsBytes) / float64(sizeU8)) * 100
	fmt.Printf("Memory Saved: %s (%.2f%% savings)\n", formatBytes(uint64(savingsBytes)), savingsPct)
	
	slowdownRatio := float64(durationEV.Nanoseconds()) / float64(durationU8.Nanoseconds())
	fmt.Printf("Speed Overhead: EnumVec is %.2fx slower than raw slice (due to bitwise operations)\n", slowdownRatio)
	fmt.Println("==================================================")
}

func formatBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
