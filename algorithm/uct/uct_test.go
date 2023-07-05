package uct

import (
	"dotg/board"
	"fmt"
	"testing"
)

func TestSimulation(t *testing.T) {
	b := board.NewBoard()
	for i := 0; i <= 100000; i++ {
		Simulation(b, 1)
	}

}

func TestSearch(t *testing.T) {
	b := board.NewBoard()
	turn := 1
	for b.Status() == 0 {
		Move(b, 10, 2000000, turn, false)
		turn ^= 3
	}

}

func BenchmarkSearch(b *testing.B) {
	bb := board.NewBoard()
	for i := 0; i < b.N; i++ {
		Search(bb, 10000, 1, 1, false)
	}
	//BenchmarkSearch-12           100          79125197 ns/op
	//BenchmarkSearch-12           100          64363073 ns/op
	//BenchmarkSearch-12           100          69319926 ns/op
}
func TestSetHashValue(t *testing.T) {

	SetHashValue(&HashKey{M: [2]uint64{1123123, 2312312312}, Now: 1}, &HashValue{
		Visit: 0,
		Win:   0,
		Turn:  0,
	})
	v, ok := GetHashValue(&HashKey{M: [2]uint64{1123123, 2312312312}, Now: 1})
	fmt.Println(v, ok)

}
