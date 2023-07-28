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
	//BenchmarkSearch-12             1        9628555300 ns/op  固定下法还扩展
	//BenchmarkSearch-12             1        6946457300 ns/op  固定下发不扩展
	//BenchmarkSearch-12             1        6771940600 ns/op  优化掉err 固定下发不扩展

}
