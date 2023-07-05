package uct

import (
	"dotg/board"
	"fmt"
	"testing"
	"time"
)

func TestUCTNode_GetUnTriedEdges(t *testing.T) {
	b := board.NewBoard()
	ms, _ := b.RandomMove()
	fmt.Println(ms)
	n := NewUCTNode(b)
	tri, err := n.GetUnTriedEdges()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(tri, len(tri))

}
func TestSimulation(t *testing.T) {
	b := board.NewBoard()
	for i := 0; i <= 100000; i++ {
		Simulation(b, 1)
	}

}
func TestSelectBest(t *testing.T) {
	b := board.NewBoard()
	for i := 0; i < 11; i++ {
		for j := 0; j < 11; j++ {
			if (i+j)&1 == 1 {
				b.State[i][j] = 1
				if i == 3 && j == 0 {
					b.State[i][j] = 0
				}

			} else if i&1 == 1 && j&1 == 1 {
				b.State[i][j] = 1
				if i == 3 && j == 1 {
					b.State[i][j] = 0
				}
			}
		}
	}
	fmt.Println(b)
	b.S[1] = 24
	es, err := Search(b, 100, 100, 1, false)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(es)
}
func TestSearch(t *testing.T) {
	b := board.NewBoard()
	for i := 0; i < 10; i++ {
		b.RandomMoveByCheck()
	}
	for b.Status() == 0 {

		start := time.Now()
		es, err := Search(b, 1000000, 20000, 1, false)
		if err != nil {
			t.Fatal(err)
		}
		b.MoveAndCheckout(es...)
		fmt.Println(es, b, time.Since(start))
		fmt.Println("-------------------------")
		if b.Status() != 0 {
			break
		}

		start = time.Now()
		es, err = Search(b, 10000000, 10000, 2, false)
		if err != nil {
			t.Fatal(err)
		}
		b.MoveAndCheckout(es...)
		fmt.Println(es, b, time.Since(start))
		fmt.Println("-------------------------")
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
