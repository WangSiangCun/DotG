package uct

import (
	"dotg/board"
	"testing"
)

func BenchmarkSearch(b *testing.B) {
	bb := board.NewBoard()
	for i := 0; i < b.N; i++ {
		Search(bb, 1, false, true)
	}
	//BenchmarkSearch
	//BenchmarkSearch-12             1        15001395100 ns/op
}
