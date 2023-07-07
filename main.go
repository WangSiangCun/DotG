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
		uct.Move(b, 10, 30000, 1, true)

		if b.Status() != 0 {
			break
		}
		//-----------------------------------------1
		//uct.Move(b, 100, 10000, 1)
		uct.Move(b, 20, 10000, 2, true)

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

		uct.Move(b, 30, 30000000, playerTurn, true)

		if b.Status() != 0 {
			break
		}

		b.GetPlayerMove()
	}
}
func main() {
	mode := 0

	/*输出后第一行是红字黑底，第二行红底白字，可以看到这里的输出与正常输出多了点东西，下面一个一个来看看。

	\033 这是标记变换颜色的起始标记，这之后的[1;31;40m则是定义颜色，1表示代码的意义或者说是显示方式，31表示前景颜色，40则是背景颜色。在这定义之后终端就会显示你设定的样式，如果只是要改变一行的样式则在结尾加入\033[0m表示恢复终端默认样式。

	颜色和配置的取值范围：

	 前景 背景 颜色
	 ---------------------------------------
	 30  40  黑色
	 31  41  红色
	 32  42  绿色
	 33  43  黄色
	 34  44  蓝色
	 35  45  紫红色
	 36  46  青蓝色
	 37  47  白色

	 代码 意义
	// -------------------------
	//  0  终端默认设置
	//  1  高亮显示
	//  4  使用下划线
	//  5  闪烁
	//  7  反白显示
	//  8  不可见*/

	fmt.Printf("\033[1;40;40m%s\033[0m\n", "输入游戏模式：1机机，2人机,3测试")
	fmt.Scan(&mode)
	if mode == 1 {
		AIToAI()
	} else if mode == 2 {
		PToAI()
	}

}
