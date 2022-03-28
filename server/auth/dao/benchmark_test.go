package dao

import "testing"
func Benchmark_Add(b *testing.B) {
    var n int
    for i := 0; i < b.N; i++ {
        n++
    }
}