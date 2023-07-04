package main

import (
	"dotg/algorithm/uct"
	"dotg/board"
	"fmt"
)

func AIToAI() {
	b := board.NewBoard()
	for i := 0; i < 10; i++ {
		b.RandomMoveByCheck()
	}
	for b.Status() == 0 {
		//uct.Move(b, 100, 20000, 1)
		uct.Move(b, 10000, 300000, 1)

		if b.Status() != 0 {
			break
		}
		//-----------------------------------------1
		//uct.Move(b, 100, 10000, 1)
		uct.Move(b, 10000, 300000, 1)

	}
}
func PToAI() {
	b := board.NewBoard()
	fmt.Println("输入您的先后手：1先手，2后手")
	playerTurn := 0
	fmt.Scan(&playerTurn)
	n := 0
	x, y := 0, 0
	if playerTurn == 1 {
		fmt.Println("输入几条边：")
		fmt.Scan(&n)
		for i := 0; i < n; i++ {
			fmt.Println("x")
			fmt.Scan(&x)
			fmt.Println("y")
			fmt.Scan(&y)
			b.MoveAndCheckout(&board.Edge{x, y})
			fmt.Println(b)
		}
	}
	playerTurn ^= 3
	for b.Status() == 0 {

		uct.Move(b, 30, 30000000, playerTurn)

		if b.Status() != 0 {
			break
		}

		fmt.Println("输入几条边：")
		fmt.Scan(&n)
		num := []*board.Edge{}
		for i := 0; i < n; i++ {
			fmt.Println("x")
			fmt.Scan(&x)
			fmt.Println("y")
			fmt.Scan(&y)
			num = append(num, &board.Edge{x, y})
		}
		b.MoveAndCheckout(&board.Edge{x, y})
		fmt.Println(b)
	}
}
func main() {
	mode := 0
	fmt.Println("输入游戏模式：1机机，2人机,3测试")
	fmt.Scan(&mode)
	if mode == 1 {
		AIToAI()
	} else if mode == 2 {
		PToAI()
	}

}
