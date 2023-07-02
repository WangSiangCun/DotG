package uct

import (
	"dotg/board"
	"fmt"
	"testing"
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
func TestUCTSearch(t *testing.T) {
	b := board.NewBoard()
	for b.Status() == 0 {
		nb := board.CopyBoard(b)
		es, err := UCTSearch(nb, 2000, 10000, 1)
		if err != nil {
			t.Fatal(err)
		}

		b.MoveAndCheckout(es...)
		nb1 := board.CopyBoard(b)
		fmt.Println(es, b)

		es, err = UCTSearch(nb1, 1000, 10000, 2)
		if err != nil {
			t.Fatal(err)
		}
		b.MoveAndCheckout(es...)
		fmt.Println(es, b)
	}

}
