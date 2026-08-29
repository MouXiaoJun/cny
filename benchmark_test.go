package cny

import (
	"math"
	"testing"
)

// BenchmarkConvert 覆盖小数值、整数值与大数值。
func BenchmarkConvert(b *testing.B) {
	values := []float64{0.05, 100, 12345678.9, 123456789012, 1.005}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := Convert(values[i%len(values)]); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkConvertCents 整数分路径。
func BenchmarkConvertCents(b *testing.B) {
	values := []int64{5, 100, 1234567890, math.MaxInt64}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := ConvertCents(values[i%len(values)]); err != nil {
			b.Fatal(err)
		}
	}
}
