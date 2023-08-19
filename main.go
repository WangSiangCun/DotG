package main

import (
	"dotg/algorithm/uct"
	"dotg/board"
	"dotg/record"
	"fmt"
	"time"
)

func AIToAI() {
	b := board.NewBoard()
	//for i := 0; i < 2; i++ {
	//	b.RandomMoveByCheck()
	//}
	for b.Status() == 0 {
		//uct.Move(b, 100, 20000, 1)
		uct.Move(b, 2, true, true)

		if b.Status() != 0 {
			break
		}
		//-----------------------------------------1
		//uct.Move(b, 100, 10000, 1)
		uct.Move(b, 2, true, false)

	}
	record.SetR("RRRR")
	record.SetB("BBBB")
	record.PrintContentStart(b.S[1], b.S[2], time.Now().String())
	record.PrintContentBack()
	record.WriteToFile(b)
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

		uct.Move(b, 1, true, false)

		if b.Status() != 0 {
			break
		}

		b.GetPlayerMove()
	}
}
func Test() {
	score1, score2 := 0, 0
	isH := true
	for {
		b := board.NewBoard()
		for b.Status() == 0 {
			uct.Move(b, 2, false, isH)
			if b.Status() != 0 {
				break
			}
			uct.Move(b, 2, false, !isH)
		}

		if b.Status() == 1 {
			if isH {
				score1++
			} else {
				score2++
			}

		} else {
			if !isH {
				score1++
			} else {
				score2++
			}
		}

		fmt.Printf("score1:score2  %v:%v isH :%v\n", score1, score2, isH)
		isH = !isH
	}
}
func FastAI() {
	b := board.NewBoard()
	for b.Status() == 0 {
		//uct.Move(b, 100, 20000, 1)
		uct.Move(b, 3, true, false)

		if b.Status() != 0 {
			break
		}
		//-----------------------------------------1
		//uct.Move(b, 100, 10000, 1)
		uct.Move(b, 3, true, false)

	}
	record.SetR("RRRR")
	record.SetB("BBBB")
	record.PrintContentStart(b.S[1], b.S[2], time.Now().String())
	record.PrintContentBack()
	record.WriteToFile(b)
}
func main() {
	mode := 0

	fmt.Printf("输入游戏模式：1机机，2人机,3测试,4快速机机")
	fmt.Scan(&mode)
	if mode == 1 {
		AIToAI()
	} else if mode == 2 {
		PToAI()
	} else if mode == 3 {
		Test()
	} else if mode == 4 {
		FastAI()
	}

}
