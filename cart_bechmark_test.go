package main

import "testing"

func BenchmarkCalculateTotalItem(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = CalculateTotalItem(cartData)
	}
}

// ======================
// BENCHMARK GOROUTINE
// ======================
func BenchmarkCalculateTotalItemConcurrent(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = CalculateTotalItemConcurrent(cartData)
	}
}
