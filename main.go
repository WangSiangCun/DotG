package dotg

import (
	"dotg/board"
	"fmt"
	"time"
)

func main() {
	now := time.Now()

	b := board.NewBoard()
	for b.Status() == 0 {

		edge, _ := b.RandomMove()

		b.CheckoutEdge(edge)

		fmt.Println(b, b.Boxes)

	}
	fmt.Println(time.Since(now))
}
