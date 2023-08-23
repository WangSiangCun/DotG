package uct

import (
	"dotg/board"
	"testing"
)

func BenchmarkSearch(b *testing.B) {
	bb := board.NewBoard()
	bb.Turn = 7
	for i := 0; i < b.N; i++ {
		Search(bb, 1, true, true)
	}
	//BenchmarkSearch
	//Tatal:100881  Tatal:133308

}
