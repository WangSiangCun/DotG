package board

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	//"sync"
)

type Board struct {
	State [11][11]int
	Turn  int
	Now   int
	S     [3]int //S[0]占位,方便1，2的下标
	Boxes []*Box
	//M     [2]uint64 //[0]为前64位0-63 [1]是剩下的64-128
	//Edges []*Edge
}

type Chain struct {
	Boxes    []*Box
	Length   int
	Type     int //0 nil, 1一格短链,2 二格短链,3 长链,4 环
	Endpoint []*Box
}
type Edge struct {
	X, Y int
}
type Box struct {
	X, Y int //"F4","F3","|_", "_|", "|￣", "￣|", " 二 ", "| |", "F1", "F0"
	Type int //0:四自由度 1: 三自由度   2:|_  3:_|  4:|￣   5: ￣| 6: 二   7. | | 8. 一自由度 9.  0
}
type DTree struct {
	X, Y  int
	Chain *Chain //如果是链，则有该属性
	Len   int
	Type  int //0 死环，1死链
}

// BoxLocation 仅用于优化
type BoxLocation struct {
	X, Y int
}

// 状态机
var (
	//A类检查 对应方向必须的类型
	aToBs = [4][3]int{{3, 5, 6}, {2, 4, 6}, {2, 3, 7}, {4, 5, 7}}
	//b类检查 类型可走的方向
	bToAs = [6][2]int{{3, 0}, {3, 1}, {2, 0}, {1, 2}, {0, 1}, {2, 3}}
	//右，左，下，上 0,1,2,3,跳度为1
	d1 = [4][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
	//右，左，下，上 0,1,2,3,跳度为2
	d2 = [4][2]int{{0, 2}, {0, -2}, {2, 0}, {-2, 0}}
	//上下可组合的
	d3 = [4][2]int{{2, 5}, {3, 5}, {0, 5}, {1, 5}}
	//左右可组合的
	d4 = [4][2]int{{1, 4}, {0, 4}, {3, 4}, {2, 4}}
	d5 = [2]int{-1, 1}
)

// 输出使用的
var (
	s1 = [10]string{"F4", "F3", "|_", "_|", "|￣", "￣|", " 二 ", "| |", "F1", "F0"}
	s2 = [5]string{"无", "一格短链", "二格短链", "长链", "环"}
)

// 前几回合，只走三自由度
// 中间几回合只走三和短链+任意四自由度=规定的数字
// 后面几回合全走
const (
	TurnMark1 int = 7
	TurnMark2 int = 16
)

// CopyBoard 拷贝棋盘
func CopyBoard(b *Board) *Board {
	//fmt.Println("front3", b)
	//rw.RLock()         // 加锁
	//defer rw.RUnlock() // 确保最终释放锁

	nB := NewBoard()
	t := 0
	for i := 0; i < 11; i += 1 {
		for j := 0; j < 11; j += 1 {
			nB.State[i][j] = b.State[i][j]
			if i&1 == 1 && j&1 == 1 {
				nB.Boxes[t].X = b.Boxes[t].X
				nB.Boxes[t].Y = b.Boxes[t].Y
				nB.Boxes[t].Type = b.Boxes[t].Type
				t++
			}

		}

	}

	nB.Now = b.Now
	nB.Turn = b.Turn
	nB.S[1] = b.S[1]
	nB.S[2] = b.S[2]
	//	nB.M[0] = b.M[0]
	//	nB.M[1] = b.M[1]
	//fmt.Println("front4", nB)
	return nB
}

// NewBoard 获得一个新棋盘
func NewBoard() *Board {
	b := &Board{
		State: [11][11]int{
			{-1, 0, -1, 0, -1, 0, -1, 0, -1, 0, -1},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{-1, 0, -1, 0, -1, 0, -1, 0, -1, 0, -1},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{-1, 0, -1, 0, -1, 0, -1, 0, -1, 0, -1},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{-1, 0, -1, 0, -1, 0, -1, 0, -1, 0, -1},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{-1, 0, -1, 0, -1, 0, -1, 0, -1, 0, -1},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{-1, 0, -1, 0, -1, 0, -1, 0, -1, 0, -1}},
		Turn: 0,
		Now:  2,
		S:    [3]int{0, 0, 0},
	}
	t := 0
	b.Boxes = make([]*Box, 25)
	for i := 1; i < 11; i += 2 {
		for j := 1; j < 11; j += 2 {
			b.Boxes[t] = &Box{i, j, 0}
			t++
		}

	}

	return b
}

// NewChain 获得新链
func NewChain() *Chain {
	return &Chain{
		Boxes:    []*Box{},
		Length:   0,
		Type:     0,
		Endpoint: []*Box{},
	}
}

// XYZToEdge 移动x,y,z所在边但不会占领
func XYZToEdge(x, y, z int) (edge *Edge) {
	i, j := 0, 0
	if x == 0 {
		i = y * 2
		j = z*2 + 1
	} else {
		i = y*2 + 1
		j = z * 2
	}
	return &Edge{i, j}
}

// EdgeToXYZ 转换Edge到x y z
func EdgeToXYZ(edge *Edge) (x, y, z int) {
	if edge.X&1 == 0 {
		x = 0
		y = (edge.X) / 2
		z = (edge.Y - 1) / 2
	} else {
		x = 1
		y = (edge.X - 1) / 2
		z = (edge.Y) / 2
	}
	return
}
func EdgesToHV(edges ...*Edge) (H, V int) {
	for _, edge := range edges {
		x, y, z := EdgeToXYZ(edge)
		if x == 0 {
			//横边
			//	fmt.Printf("%b %b\n", H, V)
			H |= 1 << (y*5 + z)
			//	fmt.Printf("%b %b\n", H, V)
		} else {
			//	fmt.Printf("%b %b\n", H, V)
			V |= 1 << (z*5 + y)
			//		fmt.Printf("%b %b\n", H, V)
		}
	}
	return H, V
}
func EdgesToM(edges ...*Edge) (M int64) {
	for _, edge := range edges {
		x, y, z := EdgeToXYZ(edge)
		//fmt.Println(x, y, z)
		if x == 0 {
			//横边
			M |= 1 << (y*5 + z)
		} else {
			M |= 1 << (z*5 + y + 30)
		}
	}
	return M

}
func MtoEdges(M int64) (es []*Edge) {
	i := 0
	for M > 0 {
		if M&1 == 1 {
			if i < 30 {
				//fmt.Println(0, i/5, i%5)
				es = append(es, XYZToEdge(0, i/5, i%5))
			} else {
				//fmt.Println(1, (i-30)%5, (i-30)/5)
				es = append(es, XYZToEdge(1, ((i-30)%5), (i-30)/5))
			}
		}
		M >>= 1
		i++
	}
	return es
}

// BoxToXY 转换boxX,boxY到x y
func BoxToXY(boxX, boxY int) (x, y int) {
	//1-1 :0,0
	//1-3 :0,1
	//3-1 :1,0
	//9-9 :4,4
	x, y = boxX/2, boxY/2
	return
}

// String 打印边
func (e *Edge) String() string {
	return strconv.Itoa(e.X) + "-" + strconv.Itoa(e.Y)
}

// String 打印格子
func (b *Box) String() string {
	return strconv.Itoa(b.X/2) + "-" + strconv.Itoa(b.Y/2) + ":" + s1[b.Type] + "   "
}

// String 打印链
func (c *Chain) String() string {
	//0 nil, 1一格短链,2 二格短链,3 长链,4 环
	var s string
	s += fmt.Sprintf("Type:%s,length:%d", s2[c.Type], c.Length)
	s += "\nchain:"
	for _, box := range c.Boxes {
		s += box.String()
	}
	s += "\nEndPoint:"
	for _, point := range c.Endpoint {
		s += point.String()
	}
	s += "\n"
	return s

}

// String 打印棋盘
func (b *Board) String() string {
	//rw.RLock()
	//defer rw.RUnlock()
	builder := strings.Builder{}
	builder.WriteString("\\ ")
	for i := 0; i < 11; i++ {
		if i != 10 {
			builder.WriteString(strconv.Itoa(i) + " ")
		} else {
			builder.WriteString("0 ")

		}

	}
	builder.WriteString("\n")
	for i := 0; i < 11; i++ {
		if i != 10 {
			builder.WriteString(strconv.Itoa(i) + " ")
		} else {
			builder.WriteString("0 ")
		}
		for j := 0; j < 11; j++ {
			if (i+j)&1 == 1 {
				if i&1 == 1 {
					//竖
					if b.State[i][j] == 1 {
						builder.WriteString("#")
					} else {
						builder.WriteString(" ")
					}

				} else {
					//横
					if b.State[i][j] == 1 {
						builder.WriteString("===")
					} else {
						builder.WriteString("   ")
					}
				}
			} else if i&1 == 0 && j&1 == 0 {
				//点
				builder.WriteString("0")
			} else {
				//占领
				if b.State[i][j] == 0 {
					builder.WriteString("   ")
				} else {
					builder.WriteString(" " + strconv.Itoa(b.State[i][j]) + " ")
				}
			}
		}
		builder.WriteString("\n")
	}
	builder.WriteString(fmt.Sprintf("Turn:%d Now:%d S[1]:%d S[2]=%d\n", b.Turn, b.Now, b.S[1], b.S[2]))
	return builder.String()
}

// IsBox 判断是否为Box
func IsBox(boxX, boxY int) bool {
	if boxX > 0 && boxX < 10 && boxY > 0 && boxY < 10 {
		return true
	} else {
		return false
	}
}

func (b *Board) GetPlayerMove() {
	n := 0
	x, y := 0, 0
	num := []*Edge{}
	for {
		fmt.Println("x")
		fmt.Scan(&x)
		fmt.Println("y")
		fmt.Scan(&y)
		num = append(num, &Edge{x, y})
		fmt.Println("------1:继续输入，2：结束--------")
		fmt.Scan(&n)
		if n == 2 {
			break
		}
	}
	b.MoveAndCheckout(&Edge{x, y})
	fmt.Println(b)
}

// Status 获得游戏状态
func (b *Board) Status() int {
	if b.S[1]+b.S[2] < 25 {
		return 0
	}
	if b.S[1] > b.S[2] {
		return 1
	} else {
		return 2
	}
}

// Move 移动所在边但不会占领
func (b *Board) Move(edges ...*Edge) {
	for _, edge := range edges {
		b.State[edge.X][edge.Y] = 1
	}
	b.Now ^= 3
	b.Turn++
	return
}

// CheckoutEdge 通常Move后调用，用以检查edges占领，若占领则加分,同时设置box
func (b *Board) CheckoutEdge(edges ...*Edge) {

	for _, edge := range edges {
		flag := 2
		if edge.X&1 == 1 {
			flag = 0
		}
		//flag=2代表为横边，flag=0代表为竖边
		for i := 0; i < 2; i++ {
			boxX := edge.X + d1[i+flag][0]
			boxY := edge.Y + d1[i+flag][1]
			tempBoxX, tempBoxY := BoxToXY(boxX, boxY)
			if boxY < 11 && boxY >= 0 && boxX < 11 && boxX >= 0 {
				f := b.GetFByBI(boxX, boxY)
				if f == 0 && b.State[boxX][boxY] == 0 {
					//fmt.Printf("%v %b%b\n", b, b.M[1], b.M[0])
					//	b.BitMove(boxX*11 + boxY)
					//fmt.Printf("%v %b%b\n", b, b.M[1], b.M[0])

					b.State[boxX][boxY] = b.Now
					b.S[b.Now]++

				}
				t := b.GetBoxType(boxX, boxY)
				b.Boxes[tempBoxX*5+tempBoxY].Type = t

			}
		}
	}

	return
}

// MoveAndCheckout Move并checkout
func (b *Board) MoveAndCheckout(edges ...*Edge) {
	b.Move(edges...)
	b.CheckoutEdge(edges...)
	return
}

// RandomMoveByCheck 随机移动,目前为GetDGridEdges()后GetEdgesByIdentifyingChains,自带checkout
func (b *Board) RandomMoveByCheck() (edge [][]*Edge) {
	ees := b.GetMove()
	if ees == nil {
		return nil
	}
	randInt := rand.Intn(len(ees))
	b.MoveAndCheckout(ees[randInt]...)
	return ees
}

// GetBoxType 获取格子的类型
// s := [10]string{"F4","F3","|_", "_|", "|￣", "￣|", " 二 ", "| |", "F1", "F0"}
func (b *Board) GetBoxType(boxX, boxY int) int {
	f := b.GetFByBI(boxX, boxY)
	if f == 4 {
		return 0
	} else if f == 3 {
		return 1
	} else if f == 2 {
		if b.State[boxX][boxY-1] == 1 && b.State[boxX+1][boxY] == 1 {
			return 2
		} else if b.State[boxX][boxY+1] == 1 && b.State[boxX+1][boxY] == 1 {
			return 3
		} else if b.State[boxX][boxY-1] == 1 && b.State[boxX-1][boxY] == 1 {
			return 4
		} else if b.State[boxX][boxY+1] == 1 && b.State[boxX-1][boxY] == 1 {
			return 5
		} else if b.State[boxX+1][boxY] == 1 && b.State[boxX-1][boxY] == 1 {
			return 6
		} else if b.State[boxX][boxY+1] == 1 && b.State[boxX][boxY-1] == 1 {
			return 7
		} else {
			return -1
		}
	} else if f == 1 {
		return 8
	} else {
		return 9
	}

}

// GetOneEdgeOfMinChain 获取最短的链的一条边
func (b *Board) GetOneEdgeOfMinChain() *Edge {
	//没有死树，只能走链
	//获取链边
	minL := 26
	var minChain *Chain
	chains := b.GetChains()
	for _, chain := range chains {
		if chain.Length < minL {
			minL = chain.Length
			minChain = chain
			if chain.Length == 1 {
				break
			}
		}

	}

	//死格
	if minChain == nil {
		return nil
	}
	//如果是二格短链则有两种方式,一种对手能双交，一种不能
	if minL == 2 {
		//获取中间的那一条
		boxX, boxY := minChain.Endpoint[0].X, minChain.Endpoint[0].Y
		for i := 0; i < 4; i++ {
			edgeX, edgeY := boxX+d1[i][0], boxY+d1[i][1]
			nextBX, nextBY := boxX+d2[i][0], boxY+d2[i][1]
			f := b.GetFByBI(nextBX, nextBY)
			if f == 2 && b.State[edgeX][edgeY] == 0 {
				return &Edge{edgeX, edgeY}
			}
		}
	}
	//如果是长链,或者一格短链
	return b.GetOneEdgeByBI(minChain.Endpoint[0].X, minChain.Endpoint[0].Y)

}

// GetOneEdgeOfChains 获取链们的一条边(含组合链)
func (b *Board) GetOneEdgeOfChains() *Edge {
	//没有死树，只能走链
	//获取链边
	minL := 26
	isHaveCircle := false
	var minChain, circleChain *Chain
	chains := b.GetChains()
	for i, chain := range chains {
		if chain.Type == 4 {
			isHaveCircle = true
			circleChain = chains[i]
		}
		if chain.Length < minL {
			minL = chain.Length
			minChain = chain
			if chain.Length == 1 {
				break
			}
		}

	}
	//死格 不可删除
	if minL == 26 {
		return nil
	}
	//如果是二格短链则有两种方式,一种对手能双交，一种不能
	if minL == 1 {
		return b.GetOneEdgeByBI(minChain.Endpoint[0].X, minChain.Endpoint[0].Y)
	} else if minL == 2 {
		//获取中间的那一条
		boxX, boxY := minChain.Endpoint[0].X, minChain.Endpoint[0].Y
		for i := 0; i < 4; i++ {
			edgeX, edgeY := boxX+d1[i][0], boxY+d1[i][1]
			nextBX, nextBY := boxX+d2[i][0], boxY+d2[i][1]
			f := b.GetFByBI(nextBX, nextBY)
			if f == 2 && b.State[edgeX][edgeY] == 0 {
				return &Edge{edgeX, edgeY}
			}
		}
	} else if isHaveCircle {
		return b.GetOneEdgeByBI(circleChain.Endpoint[0].X, circleChain.Endpoint[0].Y)
	} else if minL == 3 || minL == 6 || minL == 8 {
		//注意组合链中的环可能不只为三，但是最小了话基本是这样的
		//获取中间的那一条
		box3fX1, box3fY1, box3fX2, box3fY2 := 0, 0, 0, 0
		eX1, eY1, eX2, eY2 := 0, 0, 0, 0
		boxX, boxY := minChain.Endpoint[0].X, minChain.Endpoint[0].Y
		for i := 0; i < 4; i++ {
			edgeX, edgeY := boxX+d1[i][0], boxY+d1[i][1]
			nextBX, nextBY := boxX+d2[i][0], boxY+d2[i][1]
			f := b.GetFByBI(nextBX, nextBY)
			if f == 3 && b.State[edgeX][edgeY] == 0 {
				box3fX1, box3fY1 = nextBX, nextBY
				eX1, eY1 = edgeX, edgeY
			}
		}
		if box3fX1 == 0 && box3fY1 == 0 {
			//没找到，返回
			return b.GetOneEdgeByBI(minChain.Endpoint[0].X, minChain.Endpoint[0].Y)
		}

		boxX, boxY = minChain.Endpoint[1].X, minChain.Endpoint[1].Y
		for i := 0; i < 4; i++ {
			edgeX, edgeY := boxX+d1[i][0], boxY+d1[i][1]
			nextBX, nextBY := boxX+d2[i][0], boxY+d2[i][1]
			f := b.GetFByBI(nextBX, nextBY)
			if f == 3 && b.State[edgeX][edgeY] == 0 {
				box3fX2, box3fY2 = nextBX, nextBY
				eX2, eY2 = edgeX, edgeY
			}
		}
		if box3fX2 == 0 && box3fY2 == 0 {
			//没找到，返回
			return b.GetOneEdgeByBI(minChain.Endpoint[0].X, minChain.Endpoint[0].Y)
		}
		//找到自由度为3的那个格子
		if box3fX1 == box3fX2 && box3fY1 == box3fY2 {
			for i := 0; i < 4; i++ {
				edgeX, edgeY := box3fX1+d1[i][0], box3fY1+d1[i][1]
				//	fmt.Println(box3fX1, box3fY1, edgeX, edgeY)
				if b.State[edgeX][edgeY] == 0 {
					//与那两个边相等则继续找
					if (edgeX == eX1 && edgeY == eY1) || (edgeX == eX2 && edgeY == eY2) {
						continue
					}
					//找到了以后，返回这个边
					//	fmt.Println(b)
					return &Edge{edgeX, edgeY}
				}
			}
		}
	}
	return b.GetOneEdgeByBI(minChain.Endpoint[0].X, minChain.Endpoint[0].Y)

}

// GetEdgesOfAllChain 获取所有的链的一条边
func (b *Board) GetEdgesOfAllChain() (es []*Edge) {
	//没有死树，只能走链
	//获取链边
	chains := b.GetChains()

	for _, chain := range chains {
		boxX, boxY := chain.Endpoint[0].X, chain.Endpoint[0].Y
		if chain.Length == 2 {
			//中间的那条
			for i := 0; i < 4; i++ {
				edgeX, edgeY := boxX+d1[i][0], boxY+d1[i][1]
				nextBX, nextBY := boxX+d2[i][0], boxY+d2[i][1]
				f := b.GetFByBI(nextBX, nextBY)
				if f == 2 && b.State[edgeX][edgeY] == 0 {
					es = append(es, &Edge{edgeX, edgeY})
					break
				}
			}
		} else {
			edge := b.GetOneEdgeByBI(boxX, boxY)
			es = append(es, edge)

		}
	}
	return es

}

// CheckChainType 执行此方法会设置chain类型，基本必调用
func (c *Chain) CheckChainType() {
	if len(c.Boxes) == 1 {
		c.Type = 1
		return
	} else if len(c.Boxes) == 2 {
		c.Type = 2
		return
	} else if len(c.Boxes) > 2 {
		onePoint := c.Endpoint[0]
		twoPoint := c.Endpoint[1]
		absX, absY := math.Abs(float64(onePoint.X-twoPoint.X)), math.Abs(float64(onePoint.Y-twoPoint.Y))
		if absX+absY != 2 {
			//长链
			c.Type = 3
			return
		} else {
			//	{"F4","F3","|_", "_|", "|￣", "￣|", " 二 ", "| |", "F1", "F0"}
			if absX == 2 { //上下
				//"|_", "_|", "|￣", "￣|"," 二 ", "| |" 0:2,5  1:3,5 2:0,5 3:1,5  5:0,1,2,3
				// 0     1     2      3     4      5
				oT := onePoint.Type - 2
				tT := twoPoint.Type - 2
				if oT == 4 {
					c.Type = 3
				} else if oT == 5 {
					if tT == 0 || tT == 1 || tT == 2 || tT == 3 {
						c.Type = 4
					}
				} else {
					if tT == d3[oT][0] || tT == d3[oT][1] {
						c.Type = 4
					}
				}

			} else if absY == 2 { //左右
				//"|_", "_|", "|￣", "￣|"," 二 ", "| |"    0:1,4  1:0,4 2:3,4 3:2,4 4:0,1,2,3
				// 0     1     2      3     4       5
				oT := onePoint.Type - 2
				tT := twoPoint.Type - 2
				if oT == 5 {
					c.Type = 3
				} else if oT == 4 {
					if tT == 0 || tT == 1 || tT == 2 || tT == 3 {
						c.Type = 4
					}
				} else {
					if tT == d4[oT][0] || tT == d4[oT][1] {
						c.Type = 4
					}
				}

			} else {
				return
			}

		}

	}
	return
}

// GetSafeNo4Edge 获取除了四自由度之外的安全边
func (b *Board) GetSafeNo4Edge() (edges []*Edge) {
	/*defer func() {
		fmt.Println(edges)
	}()*/
	tempEdges := []*Edge{}
	//获取寻常边
	for i := 0; i < 11; i++ {
		for j := 0; j < 11; j++ { //正常11*11=121次 这里25次遍历,但是操作数基本一致

			if (i+j)&1 == 1 && b.State[i][j] == 0 {
				he := Edge{i, j}
				boxesF := b.GetFByE(&he)
				// 两边格子freedom不为四的边
				if (boxesF[0] >= 3 && boxesF[1] >= 3 && (boxesF[0] != 4 || boxesF[1] != 4)) || ((boxesF[0] == -1 && boxesF[1] == 3) || (boxesF[1] == -1 && boxesF[0] == 3)) {
					edges = append(edges, &he)
				} else if (boxesF[0] >= 3 || boxesF[0] == -1) && (boxesF[1] >= 3 || boxesF[1] == -1) {
					tempEdges = append(tempEdges, &he)
				}
			}

		}

	}
	if len(edges) == 0 {
		return tempEdges
	}
	return
}

// GetSafeAndChain12Edge 获取移动后不会被捕获的边和一格短链二格短链
func (b *Board) GetSafeAndChain12Edge() (edges []*Edge) {
	boxesMark := map[int]bool{}
	chains := []*Chain{}
	for i := 0; i < 11; i++ {
		for j := 0; j < 11; j++ { //正常11*11=121次 这里25次遍历,但是操作数基本一致

			if (i+j)&1 == 1 && b.State[i][j] == 0 {
				he := Edge{i, j}
				boxesF := b.GetFByE(&he)
				// 两边格子freedom大于3的边
				if (boxesF[0] >= 3 || boxesF[0] == -1) && (boxesF[1] >= 3 || boxesF[1] == -1) {
					edges = append(edges, &he)
				}
			} else if i&1 == 1 && j&1 == 1 {
				x, y := BoxToXY(i, j)
				index := x*5 + y
				//如果访问过
				if boxesMark[index] {
					continue
				}
				f := b.GetFByBI(i, j)
				if f == 2 {
					chain := NewChain()
					b.Boxes[index].Type = b.GetBoxType(i, j)
					b.GetChain(i, j, boxesMark, chain, true)
					chains = append(chains, chain)
				}

			}

		}
	}
	//一格短链,二格短链的边也可以尝试

	for _, chain := range chains {
		boxX, boxY := chain.Endpoint[0].X, chain.Endpoint[0].Y
		if chain.Length == 1 {
			edge := b.GetOneEdgeByBI(boxX, boxY)
			edges = append(edges, edge)

		} else if chain.Length == 2 {
			//中间的那条
			for i := 0; i < 4; i++ {
				edgeX, edgeY := boxX+d1[i][0], boxY+d1[i][1]
				nextBX, nextBY := boxX+d2[i][0], boxY+d2[i][1]
				f := b.GetFByBI(nextBX, nextBY)
				if f == 2 && b.State[edgeX][edgeY] == 0 {
					edges = append(edges, &Edge{edgeX, edgeY})
					break
				}
			}
		}
	}

	return edges
}

// GetEdgeBy12LChain 获得一二长度的链的可下边
func (b *Board) GetEdgeBy12LChain() (es []*Edge) {
	//一格短链,二格短链的边也可以尝试
	chains := b.GetChains()
	for _, chain := range chains {
		boxX, boxY := chain.Endpoint[0].X, chain.Endpoint[0].Y
		if chain.Length == 1 {
			edge := b.GetOneEdgeByBI(boxX, boxY)
			es = append(es, edge)

		} else if chain.Length == 2 {
			//中间的那条
			for i := 0; i < 4; i++ {
				edgeX, edgeY := boxX+d1[i][0], boxY+d1[i][1]
				nextBX, nextBY := boxX+d2[i][0], boxY+d2[i][1]
				f := b.GetFByBI(nextBX, nextBY)
				if f == 2 && b.State[edgeX][edgeY] == 0 {
					es = append(es, &Edge{edgeX, edgeY})
					break
				}
			}
		}
	}

	return es
}

// GetSafeAnd123ChainEdge 获取移动后不会被捕获的边和所有链的边
func (b *Board) GetSafeAndAllChainEdge() (edges []*Edge) {
	boxesMark := map[int]bool{}
	chains := []*Chain{}
	for i := 0; i < 11; i++ {
		for j := 0; j < 11; j++ { //正常11*11=121次 这里25次遍历,但是操作数基本一致
			if (i+j)&1 == 1 && b.State[i][j] == 0 {
				he := Edge{i, j}
				boxesF := b.GetFByE(&he)
				// 两边格子freedom大于3的边
				if (boxesF[0] >= 3 || boxesF[0] == -1) && (boxesF[1] >= 3 || boxesF[1] == -1) {
					edges = append(edges, &he)
				}
			} else if i&1 == 1 && j&1 == 1 {
				x, y := BoxToXY(i, j)
				index := x*5 + y
				//如果访问过
				if boxesMark[index] {
					continue
				}
				f := b.GetFByBI(i, j)
				if f == 2 {
					chain := NewChain()
					b.Boxes[index].Type = b.GetBoxType(i, j)
					b.GetChain(i, j, boxesMark, chain, true)
					chains = append(chains, chain)
				}

			}

		}
	}

	for _, chain := range chains {
		boxX, boxY := chain.Endpoint[0].X, chain.Endpoint[0].Y
		if chain.Length == 2 {
			//中间的那条
			for i := 0; i < 4; i++ {
				edgeX, edgeY := boxX+d1[i][0], boxY+d1[i][1]
				nextBX, nextBY := boxX+d2[i][0], boxY+d2[i][1]
				f := b.GetFByBI(nextBX, nextBY)
				if f == 2 && b.State[edgeX][edgeY] == 0 {
					edges = append(edges, &Edge{edgeX, edgeY})
					break
				}
			}
		} else if chain.Length <= 3 {
			edge := b.GetOneEdgeByBI(boxX, boxY)
			edges = append(edges, edge)

		}
	}

	return edges
}

// GetChains 获得链集合
func (b *Board) GetChains() (chains []*Chain) {

	boxesMark := map[int]bool{}
	for i := 1; i < 11; i += 2 {
		for j := 1; j < 11; j += 2 {
			x, y := BoxToXY(i, j)
			index := x*5 + y
			//如果访问过
			if boxesMark[index] {
				continue
			}

			f := b.GetFByBI(i, j)
			if f == 2 {
				chain := NewChain()
				b.Boxes[index].Type = b.GetBoxType(i, j)
				b.GetChain(i, j, boxesMark, chain, true)
				chain.CheckChainType()
				chains = append(chains, chain)
			}
		}
	}
	return

}

// GetChain 只往右下递归，并标记,boxX>=1&&<=10,
/*	A类检查 对应方向必须的类型
//	3->4,5,7 0->3,5,6 1->2,4,6 2->2,3,7
//	2 "|_", 3 "_|",4  "|￣",5 "￣|",6 " 二 ",7 "| |"
//	b类检查 类型可走的方向
//	2:3 0, 3:3 1, 4:2 0, 5:1 2, 6:0 1, 7:2 3.
//	如果这个格子已经被标记过了
*/
func (b *Board) GetChain(boxX, boxY int, boxesMark map[int]bool, chain *Chain, isStart bool) {
	x, y := BoxToXY(boxX, boxY)
	//flag
	flag := false
	index := x*5 + y
	t := b.Boxes[index].Type
	//"F4","F3","|_", "_|", "|￣", "￣|", " 二 ", "| |", "F1", "F0"
	//加入链
	chain.Boxes = append(chain.Boxes, b.Boxes[index])
	chain.Length++
	//标记
	boxesMark[index] = true

	//先由类型得到可走的方向
	for _, aD := range bToAs[t-2] {
		nextBoxX, nextBoxY := boxX+d2[aD][0], boxY+d2[aD][1]
		if nextBoxX > 0 && nextBoxX < 11 && nextBoxY > 0 && nextBoxY < 11 {

			nextX, nextY := BoxToXY(nextBoxX, nextBoxY)

			//允许的方向的下一个格子的类型
			nextIndex := nextX*5 + nextY
			nextT := b.Boxes[nextIndex].Type
			//该方向格子已经被访问
			if boxesMark[nextIndex] {
				continue
			}
			//再有方向得到下一格必须的类型
			for _, bT := range aToBs[aD] {
				//对应方向的下一个格子若有对应的类型
				if nextT == bT {
					flag = true
					b.GetChain(nextBoxX, nextBoxY, boxesMark, chain, false)

					//这个方向找到对应类型了则不可能还有了直接break去另一个方向
					break
				}
			}

		}

	}
	//当这个结点没有找到下一格的时候， 可以认为是端点
	if !flag {
		//加入端点
		chain.Endpoint = append(chain.Endpoint, b.Boxes[index])
	}
	// 如果是第一个监测
	if isStart {
		if len(chain.Endpoint) == 1 {
			chain.Endpoint = append(chain.Endpoint, b.Boxes[index])
		}

	}

	return

}

// GetFByBI 通过格子下标获得格子自由度
func (b *Board) GetFByBI(boxI, boxJ int) int {
	if boxI <= 0 || boxI >= 10 || boxJ <= 0 || boxJ >= 10 {
		return -1
	}
	freeDom := 4
	if b.State[boxI][boxJ] == 0 {
		//上
		if b.State[boxI-1][boxJ] == 1 {
			freeDom--
		}
		//下
		if b.State[boxI+1][boxJ] == 1 {
			freeDom--
		}
		//左
		if b.State[boxI][boxJ-1] == 1 {
			freeDom--
		}
		//右
		if b.State[boxI][boxJ+1] == 1 {
			freeDom--
		}
		return freeDom
	} else {
		return 0
	}

}

// GetFByE 返回边两边的freedom ,默认 左右，上下的顺序，若在边上则对应位置为-1
func (b *Board) GetFByE(edge *Edge) (boxesF [2]int) {
	i := 0
	if edge.X&1 == 1 {
		//竖边
		for _, v := range d5 {
			boxX := edge.X
			boxY := edge.Y + v
			if boxY < 11 && boxY >= 0 {
				f := b.GetFByBI(boxX, boxY)
				boxesF[i] = f
			} else {
				boxesF[i] = -1
			}
			i++

		}

	} else {
		//横边
		for _, v := range d5 {
			boxX := edge.X + v
			boxY := edge.Y
			if boxX < 11 && boxX >= 0 {
				f := b.GetFByBI(boxX, boxY)
				boxesF[i] = f
			} else {
				boxesF[i] = -1
			}
			i++

		}
	}
	return
}

// GetEdgeByBI 通过格子下标获得格子所有边
func (b *Board) GetEdgeByBI(boxI, boxJ int) (edges []*Edge) {
	f := b.GetFByBI(boxI, boxJ)
	if f != 0 {
		//上
		if b.State[boxI-1][boxJ] == 0 {
			edges = append(edges, &Edge{boxI - 1, boxJ})
		}
		//下
		if b.State[boxI+1][boxJ] == 0 {
			edges = append(edges, &Edge{boxI + 1, boxJ})
		}
		//左
		if b.State[boxI][boxJ-1] == 0 {
			edges = append(edges, &Edge{boxI, boxJ - 1})
		}
		//右
		if b.State[boxI][boxJ+1] == 0 {
			edges = append(edges, &Edge{boxI, boxJ + 1})

		}
	}
	return edges
}

// GetOneEdgeByBI 通过格子下标获得格子所有边
func (b *Board) GetOneEdgeByBI(boxI, boxJ int) (edges *Edge) {
	f := b.GetFByBI(boxI, boxJ)
	if f != 0 {
		//上
		if b.State[boxI-1][boxJ] == 0 {
			return &Edge{boxI - 1, boxJ}
		}
		//下
		if b.State[boxI+1][boxJ] == 0 {
			return &Edge{boxI + 1, boxJ}
		}
		//左
		if b.State[boxI][boxJ-1] == 0 {
			return &Edge{boxI, boxJ - 1}

		}
		//右
		if b.State[boxI][boxJ+1] == 0 {
			return &Edge{boxI, boxJ + 1}
		}
	}
	return edges
}

// IsDCircle 格子freedom为一时才可调用
func (b *Board) IsDCircle(boxX, boxY int, boxesMark map[int]bool) (is int) {
	boxesMark[(boxX/2)*5+(boxY/2)] = true
	edgesMark := map[string]bool{}
	is = b.dfsIsDCircle(boxX, boxY, boxX, boxY, 1, edgesMark, boxesMark)
	return is
}
func (b *Board) dfsIsDCircle(sBoxX, sBoxY, boxX, boxY, len int, edgesMark map[string]bool, boxesMark map[int]bool) (is int) {
	for i := 0; i < 4; i++ {
		nEX, nEY := boxX+d1[i][0], boxY+d1[i][1]
		nBX, nBY := boxX+d2[i][0], boxY+d2[i][1]
		edge := &Edge{nEX, nEY}
		if b.State[nEX][nEY] == 0 && !edgesMark[edge.String()] {
			edgesMark[edge.String()] = true
			f := b.GetFByBI(nBX, nBY)
			if f == 1 {
				ans := math.Abs(float64(sBoxX-nBX)) + math.Abs(float64(sBoxY-nBY))
				boxesMark[(nBX/2)*5+(nBY/2)] = true
				if ans == 2 {
					return len + 1
				} else {
					return -1 * (len + 1) //特殊情况，一般不会有，但是出现了就处理一下
				}
			} else if f == 2 {
				is = b.dfsIsDCircle(sBoxX, sBoxY, nBX, nBY, len+1, edgesMark, boxesMark)
				return is
			}
		}
	}
	return 0
}

// GetDCircleEdges 格子freedom为一时才可调用
func (b *Board) GetDCircleEdges(boxX, boxY, len int, isDoubleCross bool) (edges []*Edge) {
	edgesMark := map[string]bool{}
	b.dfsGetDCircleEdges(boxX, boxY, len, edgesMark, &edges, isDoubleCross)
	return edges
}
func (b *Board) dfsGetDCircleEdges(boxX, boxY, len int, edgesMark map[string]bool, edges *[]*Edge, isDoubleCross bool) {
	for i := 0; i < 4; i++ {
		nEX, nEY := boxX+d1[i][0], boxY+d1[i][1]
		nBX, nBY := boxX+d2[i][0], boxY+d2[i][1]
		edge := &Edge{nEX, nEY}
		if b.State[nEX][nEY] == 0 && !edgesMark[edge.String()] {
			edgesMark[edge.String()] = true
			if isDoubleCross && len == -1 {
				*edges = append(*edges, edge)
				break
			} else if len > 0 {
				*edges = append(*edges, edge)
			}
			b.dfsGetDCircleEdges(nBX, nBY, len-1, edgesMark, edges, isDoubleCross)
		}
	}
	return
}

// GetDChainEdges 获得死链边
func (b *Board) GetDChainEdges(box1FX, box1FY int, c *Chain, len int, isDoubleCross bool) (edges []*Edge) {

	edgesMark := map[string]bool{}
	b.dfsChainEdges(box1FX, box1FY, edgesMark, len, &edges)
	if isDoubleCross {
		endPointX, endPointY := -1, -1
		for k := 0; k < 4; k++ {
			edgeX, edgeY := box1FX+d1[k][0], box1FY+d1[k][1]
			nextBoxX, nextBoxY := box1FX+d2[k][0], box1FY+d2[k][1]
			if b.State[edgeX][edgeY] == 0 {
				//获取端点相对关系
				if nextBoxX == c.Endpoint[0].X && nextBoxY == c.Endpoint[0].Y {
					endPointX, endPointY = c.Endpoint[1].X, c.Endpoint[1].Y
				} else if nextBoxX == c.Endpoint[1].X && nextBoxY == c.Endpoint[1].Y {
					endPointX, endPointY = c.Endpoint[0].X, c.Endpoint[0].Y
				} else {
					//fmt.Println(b, c, nextBoxX, nextBoxY)
					chains := b.GetChains()
					for _, c := range chains {
						fmt.Println(c)
					}
					return nil
				}
				break
			}
		}
		for k := 0; k < 4; k++ {
			edgeX, edgeY := endPointX+d1[k][0], endPointY+d1[k][1]
			nextBoxX, nextBoxY := endPointX+d2[k][0], endPointY+d2[k][1]
			f := b.GetFByBI(nextBoxX, nextBoxY)
			if b.State[edgeX][edgeY] == 0 && f != 2 && f != 1 {
				edges = append(edges, &Edge{edgeX, edgeY})
				break
			}
		}
	}
	return
}
func (b *Board) dfsChainEdges(sBoxX, sBoxY int, edgesMark map[string]bool, len int, edges *[]*Edge) {

	if len > 0 {
		for k := 0; k < 4; k++ {
			edgeX, edgeY := sBoxX+d1[k][0], sBoxY+d1[k][1]
			nextBoxX, nextBoxY := sBoxX+d2[k][0], sBoxY+d2[k][1]
			edge := &Edge{edgeX, edgeY}
			if b.State[edgeX][edgeY] == 0 && !edgesMark[edge.String()] {
				*edges = append(*edges, edge)
				edgesMark[edge.String()] = true
				len--
				if nextBoxX >= 0 && nextBoxX <= 10 && nextBoxY >= 0 && nextBoxY <= 10 {
					b.dfsChainEdges(nextBoxX, nextBoxY, edgesMark, len, edges)
				}

			}
		}

	}
	return
}

// GetControlValue 保证此时没有短链,获取控制值
func (b *Board) GetControlValue() (cv int) {
	circlesNum, longChainsNum, threeLenChainsNum, circles4Num := 0, 0, 0, 0
	cs := b.GetChains()
	for _, c := range cs {
		c.CheckChainType()
		if c.Type == 3 {
			longChainsNum++
			if c.Length == 3 {
				threeLenChainsNum++
			}
		} else if c.Type == 4 {
			circlesNum++
			if c.Length == 4 {
				circles4Num++
			}
		}

	}

	tbValue := 0
	if circlesNum != 0 {
		if longChainsNum == threeLenChainsNum {
			if longChainsNum == 0 {
				tbValue = 8
			} else {
				tbValue = 6

			}
		} else {
			tbValue = 4
		}
	} else {
		tbValue = 4
	}
	controlValue := (25 - b.S[1] - b.S[2]) - 4*longChainsNum - 8*circlesNum + tbValue
	//fmt.Println(longChainsNum, threeLenChainsNum, circlesNum, tbValue)
	return controlValue

}

// GetDGridEdges 获得死格的边
func (b *Board) GetDGridEdges() (edges []*Edge) {
	edgesMark := make(map[string]bool)
	for i := 1; i < 11; i += 2 {
		for j := 1; j < 11; j += 2 {
			f := b.GetFByBI(i, j)
			if f == 1 {
				for k := 0; k < 4; k++ {
					edgeX := i + d1[k][0]
					edgeY := j + d1[k][1]
					tE := &Edge{edgeX, edgeY}
					if b.State[edgeX][edgeY] == 0 {
						//只有一个边，若这个边还已经加入了，则直接跳出此格寻边循环
						if edgesMark[tE.String()] {
							break
						}
						boxX := i + d2[k][0]
						boxY := j + d2[k][1]
						//说明是棋盘边上的,直接加入
						if boxX <= 0 || boxX >= 10 || boxY <= 0 || boxY >= 10 {
							edges = append(edges, tE)
							edgesMark[tE.String()] = true
							break

						}
						f1 := b.GetFByBI(boxX, boxY)
						//不为二就是死格
						if f1 != 2 {
							edges = append(edges, tE)
							edgesMark[tE.String()] = true
							//只有一个边，找到就退出此格循环
							break
						}

					}
				}

			}
		}
	}
	return
}

// GetDTreeEdges 获得死树的边，务必保证调用此方法前局面已经没有死格 t: 1 死环，2 死链
func (b *Board) GetDTreeEdges() (doubleCrossEdges, allEdges []*Edge, t int) {
	dCs, dLs := []*DTree{}, []*DTree{}
	boxesMark := map[int]bool{}

	//获取信息
	for i := 1; i < 11; i += 2 {
		for j := 1; j < 11; j += 2 {
			f := b.GetFByBI(i, j)
			if f == 1 && !boxesMark[(i/2*5+(j/2))] {
				//先判断是不是死环
				//如果有两头1，则是死环这一类的，如果没有，则为死链，不用担心已经访问过的会再次访问
				is := b.IsDCircle(i, j, boxesMark)
				if is > 0 || is < 0 {
					dCs = append(dCs, &DTree{i, j, nil, is, 0})
					continue
				} else {
					//is==0不是死环,进入死链
					//死链
					for k := 0; k < 4; k++ {
						edgeX := i + d1[k][0]
						edgeY := j + d1[k][1]
						if b.State[edgeX][edgeY] == 0 {
							boxX := i + d2[k][0]
							boxY := j + d2[k][1]
							f1 := b.GetFByBI(boxX, boxY)

							if f1 == 2 {
								chain := NewChain()
								b.GetChain(boxX, boxY, boxesMark, chain, true)
								dLs = append(dLs, &DTree{i, j, chain, chain.Length + 1, 1})
							}

						}
					}

				}

			}
		}
	}
	if len(dCs) != 0 {
		//只有死环，剩一个来双交或全吃
		edges := []*Edge{}
		for i := 0; i < len(dCs)-1; i++ {
			//全吃
			es := b.GetDCircleEdges(dCs[i].X, dCs[i].Y, dCs[i].Len-1, false)
			edges = append(edges, es...)

		}
		//全捕获
		edgesTemp := []*Edge{}
		es := b.GetDCircleEdges(dCs[len(dCs)-1].X, dCs[len(dCs)-1].Y, dCs[len(dCs)-1].Len-1, false)
		edgesTemp = append(edgesTemp, es...)
		edgesTemp = append(edgesTemp, edges...)
		allEdges = append(allEdges, edgesTemp...)

		//双交
		edgesTemp = []*Edge{}
		es = b.GetDCircleEdges(dCs[len(dCs)-1].X, dCs[len(dCs)-1].Y, dCs[len(dCs)-1].Len-4, true)

		edgesTemp = append(edgesTemp, es...)
		edgesTemp = append(edgesTemp, edges...)

		doubleCrossEdges = append(doubleCrossEdges, edgesTemp...)

		t = 1
	} else if len(dLs) != 0 {
		edges := []*Edge{}
		boxesMark := map[int]bool{}
		//只有死链，剩一个来双交或全吃
		i, j := 0, 0
		for l := 0; l < len(dLs)-1; l++ {
			//死链
			i = dLs[l].X
			j = dLs[l].Y
			for k := 0; k < 4; k++ {
				edgeX := i + d1[k][0]
				edgeY := j + d1[k][1]
				if b.State[edgeX][edgeY] == 0 {
					boxX := i + d2[k][0]
					boxY := j + d2[k][1]
					f1 := b.GetFByBI(boxX, boxY)
					if f1 == 2 {
						chain := dLs[l].Chain
						b.GetChain(boxX, boxY, boxesMark, chain, true)

						//全捕获
						es := b.GetDChainEdges(i, j, chain, dLs[l].Len, false)
						edges = append(edges, es...)

					}

				}
			}

		}
		//	fmt.Println(b)
		i = dLs[len(dLs)-1].X
		j = dLs[len(dLs)-1].Y
		for k := 0; k < 4; k++ {
			edgeX := i + d1[k][0]
			edgeY := j + d1[k][1]
			if b.State[edgeX][edgeY] == 0 {
				boxX := i + d2[k][0]
				boxY := j + d2[k][1]
				f1 := b.GetFByBI(boxX, boxY)
				if f1 == 2 {
					chain := dLs[len(dLs)-1].Chain
					//全捕获
					es := b.GetDChainEdges(i, j, chain, dLs[len(dLs)-1].Len, false)

					//fmt.Println("全捕获:", es, dLs[len(dLs)-1].Len)
					edgesTemp := []*Edge{}
					edgesTemp = append(edgesTemp, es...)
					edgesTemp = append(edgesTemp, edges...) //edges里是之前的全捕获
					allEdges = append(allEdges, edgesTemp...)

					//双交
					es = b.GetDChainEdges(i, j, chain, dLs[len(dLs)-1].Len-2, true)

					edgesTemp = []*Edge{}
					//fmt.Println("双交:", es)
					edgesTemp = append(edgesTemp, es...)
					edgesTemp = append(edgesTemp, edges...)
					doubleCrossEdges = append(doubleCrossEdges, edgesTemp...)

				}

			}
		}
		t = 2

	} else {
		//没有死环或死链
		return
	}
	return

}

// Get2FEdge 获取移动后不会被捕获的边
func (b *Board) Get2FEdge() (edges []*Edge) {
	/*defer func() {
		fmt.Println(edges)
	}()*/
	//获取寻常边
	for i := 0; i < 11; i++ {
		for j := 0; j < 11; j++ { //正常11*11=121次 这里25次遍历,但是操作数基本一致

			if (i+j)&1 == 1 && b.State[i][j] == 0 {
				he := Edge{i, j}
				boxesF := b.GetFByE(&he)
				// 两边格子freedom大于3的边
				if (boxesF[0] >= 3 || boxesF[0] == -1) && (boxesF[1] >= 3 || boxesF[1] == -1) {

					//fmt.Println(b, b.State[i][j], i, j, boxesF)
					edges = append(edges, &he)
				}
			}

		}

	}
	//fmt.Println("2F:", edges)
	return
}

// GetFrontMoveByTurn 存在安全边时的走法 获取前期走法边
func (b *Board) GetFrontMoveByTurn() (ees [][]*Edge) {
	nB := CopyBoard(b)
	//存在安全边
	edges2f := nB.Get2FEdge()
	if len(edges2f) > 0 {
		preEdges := []*Edge{}
		//存在安全边
		//获取死格
		dGEdges := nB.GetDGridEdges()
		if len(dGEdges) > 0 {
			//模拟 局面不可有死格
			nB.MoveAndCheckout(dGEdges...)
			preEdges = append(preEdges, dGEdges...)
		}

		//获取死树的全吃走法
		doubleCrossEdges, allEdges, _ := nB.GetDTreeEdges()
		if len(allEdges) > 0 {
			nB.MoveAndCheckout(allEdges...)
			preEdges = append(preEdges, allEdges...)
		}
		switch {
		case b.Turn == 0:
			//0肯定是先手
			ees = append(ees, []*Edge{&Edge{4, 5}})
		case b.Turn < TurnMark1:
			//前几回合，只走三自由度
			es := nB.GetSafeNo4Edge()
			for _, e := range es {
				tempEdges := []*Edge{}
				tempEdges = append(tempEdges, preEdges...)
				tempEdges = append(tempEdges, e)
				ees = append(ees, tempEdges)
			}
			if doubleCrossEdges != nil {
				tempEdges := []*Edge{}
				tempEdges = append(tempEdges, dGEdges...)
				tempEdges = append(tempEdges, doubleCrossEdges...)
				ees = append(ees, tempEdges)
			}
		case b.Turn < TurnMark2:
			es := nB.GetSafeAndChain12Edge()
			for _, e := range es {
				tempEdges := []*Edge{}
				tempEdges = append(tempEdges, preEdges...)
				tempEdges = append(tempEdges, e)
				ees = append(ees, tempEdges)
			}
			if doubleCrossEdges != nil {
				tempEdges := []*Edge{}
				tempEdges = append(tempEdges, dGEdges...)
				tempEdges = append(tempEdges, doubleCrossEdges...)
				ees = append(ees, tempEdges)
			}
		default:
			es := nB.GetSafeAndAllChainEdge()
			for _, e := range es {
				tempEdges := []*Edge{}
				tempEdges = append(tempEdges, preEdges...)
				tempEdges = append(tempEdges, e)
				ees = append(ees, tempEdges)
			}
			if doubleCrossEdges != nil {
				tempEdges := []*Edge{}
				tempEdges = append(tempEdges, dGEdges...)
				tempEdges = append(tempEdges, doubleCrossEdges...)
				ees = append(ees, tempEdges)
			}
		}

	}
	return
}

// Get2FEdgeAndMessage 获取2FE和一些信息
func (b *Board) Get2FEdgeAndMessage() (edges []*Edge, f1Boxes []*BoxLocation) {
	/*defer func() {
		fmt.Println(edges)
	}()*/
	//获取寻常边
	for i := 0; i < 11; i++ {
		for j := 0; j < 11; j++ { //正常11*11=121次 这里25次遍历,但是操作数基本一致

			if (i+j)&1 == 1 && b.State[i][j] == 0 {
				he := Edge{i, j}
				boxesF := b.GetFByE(&he)
				// 两边格子freedom大于3的边
				if (boxesF[0] >= 3 || boxesF[0] == -1) && (boxesF[1] >= 3 || boxesF[1] == -1) {

					//fmt.Println(b, b.State[i][j], i, j, boxesF)
					edges = append(edges, &he)
				}
			} else if i&1 == 1 && j&1 == 1 {

				boxesF := b.GetFByBI(i, j)
				if boxesF == 1 {
					f1Boxes = append(f1Boxes, &BoxLocation{i, j})
				}
			}

		}

	}
	//fmt.Println("2F:", edges)
	return
}

// GetDGridEdgesByMessage 获得死格的边
func (b *Board) GetDGridEdgesByMessage(f1Boxes []*BoxLocation) (edges []*Edge) {

	edgesMark := make(map[string]bool)
	i, j, k := 0, 0, 0

	for index := 0; index < len(f1Boxes); index++ {

		i, j = f1Boxes[index].X, f1Boxes[index].Y
		for k = 0; k < 4; k++ {
			edgeX := i + d1[k][0]
			edgeY := j + d1[k][1]
			tE := &Edge{edgeX, edgeY}
			if b.State[edgeX][edgeY] == 0 {
				//只有一个边，若这个边还已经加入了，则直接跳出此格寻边循环,防止一箭双雕型死格
				if edgesMark[tE.String()] {
					break
				}
				boxX := i + d2[k][0]
				boxY := j + d2[k][1]
				//说明是棋盘边上的,直接加入
				if boxX <= 0 || boxX >= 10 || boxY <= 0 || boxY >= 10 {
					edges = append(edges, tE)
					edgesMark[tE.String()] = true
					break

				}
				f1 := b.GetFByBI(boxX, boxY)
				//不为二就是死格
				if f1 != 2 {
					edges = append(edges, tE)
					edgesMark[tE.String()] = true

					//移除此f1Box

					f1Boxes[index], f1Boxes[len(f1Boxes)-1] = f1Boxes[len(f1Boxes)-1], f1Boxes[index]
					f1Boxes = f1Boxes[:len(f1Boxes)-1]
					index--

					//只有一个边，找到就退出此格循环
					break
				}

			}
		}

	}
	return
}

// GetDTreeEdgesByMessage 获得死树的边，务必保证调用此方法前局面已经没有死格 t: 1 死环，2 死链
func (b *Board) GetDTreeEdgesByMessage(f1Boxes []*BoxLocation) (doubleCrossEdges, allEdges []*Edge, t int) {
	dCs, dLs := []*DTree{}, []*DTree{}
	boxesMark := map[int]bool{}
	i, j := 0, 0
	for _, f1Box := range f1Boxes {
		i, j = f1Box.X, f1Box.Y
		if !boxesMark[(i/2*5 + (j / 2))] {
			//先判断是不是死环
			//如果有两头1，则是死环这一类的，如果没有，则为死链，不用担心已经访问过的会再次访问
			is := b.IsDCircle(i, j, boxesMark)
			if is > 0 || is < 0 {
				dCs = append(dCs, &DTree{i, j, nil, is, 0})
				continue
			} else {
				//is==0不是死环,进入死链
				//死链
				for k := 0; k < 4; k++ {
					edgeX := i + d1[k][0]
					edgeY := j + d1[k][1]
					if b.State[edgeX][edgeY] == 0 {
						boxX := i + d2[k][0]
						boxY := j + d2[k][1]
						f1 := b.GetFByBI(boxX, boxY)

						if f1 == 2 {
							chain := NewChain()
							b.GetChain(boxX, boxY, boxesMark, chain, true)
							dLs = append(dLs, &DTree{i, j, chain, chain.Length + 1, 1})
						}

					}
				}

			}

		}
	}
	//获取信息
	if len(dCs) != 0 {
		//只有死环，剩一个来双交或全吃
		edges := []*Edge{}
		for i := 0; i < len(dCs)-1; i++ {
			//全吃
			es := b.GetDCircleEdges(dCs[i].X, dCs[i].Y, dCs[i].Len-1, false)
			edges = append(edges, es...)

		}
		//全捕获
		edgesTemp := []*Edge{}
		es := b.GetDCircleEdges(dCs[len(dCs)-1].X, dCs[len(dCs)-1].Y, dCs[len(dCs)-1].Len-1, false)
		edgesTemp = append(edgesTemp, es...)
		edgesTemp = append(edgesTemp, edges...)
		allEdges = append(allEdges, edgesTemp...)

		//双交
		edgesTemp = []*Edge{}
		es = b.GetDCircleEdges(dCs[len(dCs)-1].X, dCs[len(dCs)-1].Y, dCs[len(dCs)-1].Len-4, true)

		edgesTemp = append(edgesTemp, es...)
		edgesTemp = append(edgesTemp, edges...)

		doubleCrossEdges = append(doubleCrossEdges, edgesTemp...)

		t = 1
	} else if len(dLs) != 0 {
		edges := []*Edge{}
		boxesMark := map[int]bool{}
		//只有死链，剩一个来双交或全吃
		i, j := 0, 0
		for l := 0; l < len(dLs)-1; l++ {
			//死链
			i = dLs[l].X
			j = dLs[l].Y
			for k := 0; k < 4; k++ {
				edgeX := i + d1[k][0]
				edgeY := j + d1[k][1]
				if b.State[edgeX][edgeY] == 0 {
					boxX := i + d2[k][0]
					boxY := j + d2[k][1]
					f1 := b.GetFByBI(boxX, boxY)
					if f1 == 2 {
						chain := dLs[l].Chain
						b.GetChain(boxX, boxY, boxesMark, chain, true)

						//全捕获
						es := b.GetDChainEdges(i, j, chain, dLs[l].Len, false)
						edges = append(edges, es...)

					}

				}
			}

		}
		//	fmt.Println(b)
		i = dLs[len(dLs)-1].X
		j = dLs[len(dLs)-1].Y
		for k := 0; k < 4; k++ {
			edgeX := i + d1[k][0]
			edgeY := j + d1[k][1]
			if b.State[edgeX][edgeY] == 0 {
				boxX := i + d2[k][0]
				boxY := j + d2[k][1]
				f1 := b.GetFByBI(boxX, boxY)
				if f1 == 2 {
					chain := dLs[len(dLs)-1].Chain
					//全捕获
					es := b.GetDChainEdges(i, j, chain, dLs[len(dLs)-1].Len, false)

					//fmt.Println("全捕获:", es, dLs[len(dLs)-1].Len)
					edgesTemp := []*Edge{}
					edgesTemp = append(edgesTemp, es...)
					edgesTemp = append(edgesTemp, edges...) //edges里是之前的全捕获
					allEdges = append(allEdges, edgesTemp...)

					//双交
					es = b.GetDChainEdges(i, j, chain, dLs[len(dLs)-1].Len-2, true)

					edgesTemp = []*Edge{}
					//fmt.Println("双交:", es)
					edgesTemp = append(edgesTemp, es...)
					edgesTemp = append(edgesTemp, edges...)
					doubleCrossEdges = append(doubleCrossEdges, edgesTemp...)

				}

			}
		}
		t = 2

	} else {
		//没有死环或死链
		return
	}
	return

}

// GetFrontMoveByTurnAndMessage 存在安全边时的走法 获取前期走法边
func (b *Board) GetFrontMoveByTurnAndMessage() (ees [][]*Edge) {
	nB := CopyBoard(b)
	//存在安全边
	edges2f, f1Boxes := nB.Get2FEdgeAndMessage()
	if len(edges2f) > 0 {
		preEdges := []*Edge{}
		//存在安全边
		//获取死格
		dGEdges := nB.GetDGridEdgesByMessage(f1Boxes)
		if len(dGEdges) > 0 {
			//模拟 局面不可有死格
			nB.MoveAndCheckout(dGEdges...)
			preEdges = append(preEdges, dGEdges...)
		}

		//获取死树的全吃走法
		doubleCrossEdges, allEdges, _ := nB.GetDTreeEdgesByMessage(f1Boxes)
		if len(allEdges) > 0 {
			nB.MoveAndCheckout(allEdges...)
			preEdges = append(preEdges, allEdges...)
		}
		switch {
		case b.Turn == 0:
			//0肯定是先手
			ees = append(ees, []*Edge{&Edge{4, 5}})
		case b.Turn < TurnMark1:
			//前几回合，只走三自由度
			es := nB.GetSafeNo4Edge()
			for _, e := range es {
				tempEdges := []*Edge{}
				tempEdges = append(tempEdges, preEdges...)
				tempEdges = append(tempEdges, e)
				ees = append(ees, tempEdges)
			}
			if doubleCrossEdges != nil {
				tempEdges := []*Edge{}
				tempEdges = append(tempEdges, dGEdges...)
				tempEdges = append(tempEdges, doubleCrossEdges...)
				ees = append(ees, tempEdges)
			}
		case b.Turn < TurnMark2:
			es := nB.GetSafeAndChain12Edge()
			for _, e := range es {
				tempEdges := []*Edge{}
				tempEdges = append(tempEdges, preEdges...)
				tempEdges = append(tempEdges, e)
				ees = append(ees, tempEdges)
			}
			if doubleCrossEdges != nil {
				tempEdges := []*Edge{}
				tempEdges = append(tempEdges, dGEdges...)
				tempEdges = append(tempEdges, doubleCrossEdges...)
				ees = append(ees, tempEdges)
			}
		default:
			es := nB.GetSafeAndAllChainEdge()
			for _, e := range es {
				tempEdges := []*Edge{}
				tempEdges = append(tempEdges, preEdges...)
				tempEdges = append(tempEdges, e)
				ees = append(ees, tempEdges)
			}
			if doubleCrossEdges != nil {
				tempEdges := []*Edge{}
				tempEdges = append(tempEdges, dGEdges...)
				tempEdges = append(tempEdges, doubleCrossEdges...)
				ees = append(ees, tempEdges)
			}
		}

	}
	return
}

// GetEndMove 不存在安全边时的走法
func (b *Board) GetEndMove() (es []*Edge) {
	nB := CopyBoard(b)
	//不存在安全边

	preEdges := []*Edge{}
	//获取死格
	dGEdges := nB.GetDGridEdges()
	if len(dGEdges) > 0 {
		//模拟 局面不可有死格
		nB.MoveAndCheckout(dGEdges...)
		preEdges = append(preEdges, dGEdges...)
	}
	//获取死树的全吃和双交走法
	doubleCrossEdges, allEdges, t := nB.GetDTreeEdges()
	if len(doubleCrossEdges) == 0 && len(allEdges) == 0 {
		//没有死树，只能走链
		//获取链边
		edge := nB.GetOneEdgeOfChains()
		if edge == nil {
			//没有链看，游戏也没结束，也就是只有死格
			es = append(es, preEdges...)
			return es
		} else {
			//有链
			es = append(es, preEdges...)
			es = append(es, edge)

		}
	} else {
		//有死树
		//全吃后走链
		//模拟全吃死树,能结束游戏就选择，否则双交或全吃
		nB.MoveAndCheckout(allEdges...)
		if nB.Status() != 0 {
			es = append(es, preEdges...)
			es = append(es, allEdges...)
		} else {

			//是个环
			criticalValue := 0
			if t == 1 {
				criticalValue = 4
			} else if t == 2 {
				criticalValue = 2
			} else {
				return nil
			}
			cV := nB.GetControlValue()

			//全吃，吃完还要走个链
			edge := nB.GetOneEdgeOfChains()
			if edge != nil {
				if cV > criticalValue {
					//保持控制的走法
					//双交
					tempEdges := []*Edge{}
					tempEdges = append(tempEdges, preEdges...)
					tempEdges = append(tempEdges, doubleCrossEdges...)
					es = append(es, preEdges...)
					es = append(es, doubleCrossEdges...)
				} else {

					es = append(es, preEdges...)
					es = append(es, allEdges...)
					es = append(es, edge)

				}

			}

		}

	}

	return es

}

// GetMove 获取移动 若没到EndGame则会返回前期的移动，否则返回nil
func (b *Board) GetMove() (ees [][]*Edge) {
	//获取前期走法边
	ees = b.GetFrontMoveByTurnAndMessage()
	if len(ees) > 0 {
		return ees
	} else {
		//不存在安全边
		for b.Status() == 0 {
			endMoves := b.GetEndMove()
			b.MoveAndCheckout(endMoves...)
		}
		return nil

	}
}
