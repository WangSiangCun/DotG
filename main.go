package main

import (
	"dotg/algorithm/uct"
	"dotg/board"
	"fmt"
	"time"
)

func main() {
	mode := 0
	fmt.Println("输入游戏模式：1机机，2人机,3测试")
	for {
		fmt.Scan(&mode)
		if mode == 1 {
			b := board.NewBoard()
			for i := 0; i < 10; i++ {
				b.RandomMoveByCheck()
			}
			for b.Status() == 0 {

				start := time.Now()
				es, err := uct.Search(b, 100, 20000, 1)
				if err != nil {
					fmt.Println(err)
					return
				}
				b.MoveAndCheckout(es...)
				fmt.Println(es, b, time.Since(start))
				fmt.Println("-------------------------")
				if b.Status() != 0 {
					break
				}

				start = time.Now()
				es, err = uct.Search(b, 0, 10000, 2)
				if err != nil {
					fmt.Println(err)
					return
				}
				b.MoveAndCheckout(es...)
				fmt.Println(es, b, time.Since(start))
				fmt.Println("-------------------------")
			}
		} else if mode == 2 {
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
			for b.Status() == 0 {
				start := time.Now()
				es, err := uct.Search(b, 20000, 1, 1)
				if err != nil {
					fmt.Println(err)
					return
				}
				b.MoveAndCheckout(es...)
				fmt.Println(es)
				fmt.Println(b, time.Since(start))
				fmt.Println("-------------------------")
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

		} else if mode == 3 {
			oneSCore, twoScore := 0, 0
			for i := 0; i < 1000; i++ {
				b := board.NewBoard()
				for b.Status() == 0 {
					start := time.Now()
					es, err := uct.Search(b, 100, 200000, 1)
					if err != nil {
						fmt.Println(err)
						return
					}
					b.MoveAndCheckout(es...)
					fmt.Println(es, b, time.Since(start))
					fmt.Println("-------------------------")
					if b.Status() != 0 {
						break
					}

					start = time.Now()
					es, err = uct.Search(b, 0, 100000, 2)
					if err != nil {
						fmt.Println(err)
						return
					}
					b.MoveAndCheckout(es...)
					fmt.Println(es, b, time.Since(start))
					fmt.Println("-------------------------")
				}
				if b.Status() == 1 {
					oneSCore++
				} else {
					twoScore++
				}
				fmt.Printf("S:%d,%d\n", oneSCore, twoScore)

			}

		}
	}

}
