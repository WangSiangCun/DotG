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
		Move(b, 30, 2000000, turn, true, true)
		turn ^= 3
		if b.Status() != 0 {
			break
		}
		Move(b, 30, 2000000, turn, true, false)
	}

}
func BenchmarkSearch(b *testing.B) {
	bb := board.NewBoard()
	for i := 0; i < b.N; i++ {
		Search(bb, 10000, 100000, 1, false, true)
	}
	//BenchmarkSearch-12             1        9628555300 ns/op  固定下法还扩展
	//BenchmarkSearch-12             1        6946457300 ns/op  固定下发不扩展
	//BenchmarkSearch-12             1        6771940600 ns/op  优化掉err 固定下发不扩展

}
func TestMutex(t *testing.T) {
	b := board.NewBoard()
	ThreadNum = 11
	var (
		exit = make(chan int, ThreadNum)
		stop = make(chan int, ThreadNum)
	)
	maxDeep = 0
	root := NewUCTNode(b)
	res := 0
	for i := 0; i < ThreadNum; i++ {
		go func() {

			for len(stop) == 0 {

				if root.Visit > 100000 {
					stop <- 1
				}
				nowN := root
				deep := 0
				for next := SelectBest(nowN); next != nil; {
					nowN = next
					next = SelectBest(nowN)
					deep++
				}
				if nowN == nil {
					fmt.Println("select得到了NULL")
				}

				if deep > maxDeep {
					maxDeep = deep
				}

				if nowN.B.Status() == 0 {
					nowN = Expand(nowN, true)
				}
				//nB仅仅用于模拟
				nB := board.CopyBoard(nowN.B)

				res = Simulation(nB, 1)

				for nowN != nil {
					nowN.BackUp(res, 1)
					nowN = nowN.Parents
				}
			}
			exit <- 1
		}()
	}
	for i := 0; i < ThreadNum; i++ {
		<-exit
	}
	GetBestChild(root, true)

}
