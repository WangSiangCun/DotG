package main

import (
	"dotg/algorithm/uct"
	"dotg/board"
	"fmt"
)

func AIToAI() {
	b := board.NewBoard()
	//for i := 0; i < 2; i++ {
	//	b.RandomMoveByCheck()
	//}
	for b.Status() == 0 {
		//uct.Move(b, 100, 20000, 1)
		uct.Move(b, 1000, 300000, 1, true, false)

		if b.Status() != 0 {
			break
		}
		//-----------------------------------------1
		//uct.Move(b, 100, 10000, 1)
		uct.Move(b, 2000, 100000, 2, true, false)

	}
}
func PToAI() {
	b := board.NewBoard()
	fmt.Printf("\033[1;40;40m%s\033[0m\n", "输入您的先后手：1先手，2后手")
	playerTurn := 0
	fmt.Scan(&playerTurn)

	if playerTurn == 1 {
		b.GetPlayerMove()
	}
	playerTurn ^= 3
	for b.Status() == 0 {

		uct.Move(b, 30, 30000000, playerTurn, true, false)

		if b.Status() != 0 {
			break
		}

		b.GetPlayerMove()
	}
}
func main() {
	mode := 0

	fmt.Printf("输入游戏模式：1机机，2人机,3测试")
	fmt.Scan(&mode)
	if mode == 1 {
		AIToAI()
	} else if mode == 2 {
		PToAI()
	} else if mode == 3 {
		score1, score2 := 0, 0
		for {
			b := board.NewBoard()
			turn := 1
			for b.Status() == 0 {
				uct.Move(b, 10, 2000000, turn, false, true)
				turn ^= 3
				if b.Status() != 0 {
					break
				}
				uct.Move(b, 10, 2000000, turn, false, false)
			}
			if b.Status() == 1 {
				score1++
			} else {
				score2++
			}
			fmt.Printf("score1:score2   %v:%v \n", score1, score2)
		}

	}

}
