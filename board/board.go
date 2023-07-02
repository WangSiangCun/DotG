package board

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
)

type Board struct {
	State [11][11]int
	Turn  int
	Now   int
	S     [3]int //S[0]占位,方便1，2的下标
	Boxes []*Box
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
)

// 输出使用的
var (
	s1 = [10]string{"F4", "F3", "|_", "_|", "|￣", "￣|", " 二 ", "| |", "F1", "F0"}
	s2 = [5]string{"无", "一格短链", "二格短链", "长链", "环"}
)

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
		Now:  1,
		S:    [3]int{0, 0, 0},
	}
	for i := 1; i < 11; i += 2 {
		for j := 1; j < 11; j += 2 {
			b.Boxes = append(b.Boxes, &Box{i, j, 0})
		}

	}

	return b
}

// CopyBoard 拷贝棋盘
func CopyBoard(b *Board) *Board {
	nB := NewBoard()
	t := 0
	for i := 1; i < 11; i += 1 {
		for j := 1; j < 11; j += 1 {
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
	nB.S[0] = nB.S[0]
	nB.S[1] = nB.S[1]
	return nB
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
func XYZToEdge(x, y, z int) (edge *Edge, err error) {
	i, j := 0, 0
	if x == 0 {
		//0-1 0-0
		//1-3 1-2
		//2-5 2-4
		//3-7
		//4-9
		i = y * 2
		j = z*2 + 1
	} else {
		i = y*2 + 1
		j = z * 2
	}
	return &Edge{i, j}, err
}

// EdgeToXYZ 转换Edge到x y z
func EdgeToXYZ(edge *Edge) (x, y, z int, err error) {
	if edge.X&1 == 0 {
		x = 0
		y = (edge.X - 1) / 2
		z = (edge.Y) / 2
	} else {
		x = 1
		y = (edge.X) / 2
		z = (edge.Y - 1) / 2
	}
	return
}

// BoxToXY 转换boxX,boxY到x y
func BoxToXY(boxX, boxY int) (x, y int, err error) {
	//1-1 :0,0
	//1-3 :0,1
	//3-1 :1,0
	//9-9 :4,4
	if boxX&1 != 1 || boxY&1 != 1 {
		return -1, -1, fmt.Errorf("并非格子坐标")
	}
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

// Move 移动所在边但不会占领
func (b *Board) Move(edges ...*Edge) error {
	for _, edge := range edges {
		if b.State[edge.X][edge.Y] != 0 {
			s := fmt.Sprintf("%s\n", b.String())
			s += "X,Y: " + strconv.Itoa(edge.X) + " " + strconv.Itoa(edge.Y) + "\n"
			return fmt.Errorf("repeated Move\n" + s)

		}
		b.State[edge.X][edge.Y] = 1
	}
	b.Now ^= 3
	b.Turn++
	return nil
}

// CheckoutEdge 通常Move后调用，用以检查edges占领，若占领则加分,同时设置box
func (b *Board) CheckoutEdge(edges ...*Edge) error {

	for _, edge := range edges {
		flag := 2
		if edge.X&1 == 1 {
			flag = 0
		}
		//flag=2代表为横边，flag=0代表为竖边
		for i := 0; i < 2; i++ {
			boxX := edge.X + d1[i+flag][0]
			boxY := edge.Y + d1[i+flag][1]
			tempBoxX, tempBoxY, BoxToXYErr := BoxToXY(boxX, boxY)
			if BoxToXYErr != nil {
				return BoxToXYErr
			}
			if boxY < 11 && boxY >= 0 && boxX < 11 && boxX >= 0 {
				f, err := b.GetFByBI(boxX, boxY)
				if err != nil {
					return err
				}
				if f == 0 && b.State[boxX][boxY] == 0 {
					b.State[boxX][boxY] = b.Now
					b.S[b.Now]++
					if b.S[1]+b.S[2] == 24 {
						a := false
						for _, box := range b.Boxes {
							if box.Type != 9 {
								a = true
							}

						}
						if !a {
							fmt.Println(b)
						}

					}
				}
				t, getBoxTypeOf2FErr := b.GetBoxType(boxX, boxY)
				if getBoxTypeOf2FErr != nil {
					return getBoxTypeOf2FErr
				}
				b.Boxes[tempBoxX*5+tempBoxY].Type = t
			}
		}
	}

	return nil
}

// MoveAndCheckout Move并checkout
func (b *Board) MoveAndCheckout(edges ...*Edge) error {
	if err := b.Move(edges...); err != nil {
		return err
	} else if err = b.CheckoutEdge(edges...); err != nil {
		return err
	}
	return nil
}

// GetMove 获得checkout后的moves ，会首先占领死格
func (b *Board) GetMove() (ees [][]*Edge, err error) {
	flag := false
	//占领所有死格
	if dGEs, err := b.GetDGridEdges(); err != nil {
		return nil, err
	} else {
		for _, dGE := range dGEs {
			//fmt.Println(dGE)
			if err = b.MoveAndCheckout(dGE); err != nil {
				return nil, err
			}
			flag = true
		}
	}
	//获得校验后的边
	ees, err = b.GetEdgesByIdentifyingChains()
	//for _, es := range ees {
	//fmt.Println(es)
	//}
	if err != nil {
		return nil, err
	}
	if len(ees) == 0 && !flag {
		fmt.Println(b)
		for i := 1; i < 11; i += 2 {
			for j := 1; j < 11; j += 2 {
				t, _ := b.GetBoxType(i, j)
				fmt.Println(t)
			}
		}

		chains, err := b.GetChains()
		if err != nil {
			return nil, err
		}
		for _, c := range chains {
			fmt.Println(c)
		}
		return nil, fmt.Errorf("没有可移动的边")
	}
	return ees, nil
}

// GetBoxType 获取格子的类型
// s := [10]string{"F4","F3","|_", "_|", "|￣", "￣|", " 二 ", "| |", "F1", "F0"}
func (b *Board) GetBoxType(boxX, boxY int) (int, error) {

	if boxX&1 != 1 || boxY&1 != 1 {
		return -1, fmt.Errorf("坐标并非格子")
	}
	f, err := b.GetFByBI(boxX, boxY)
	if err != nil {
		return -1, err
	}
	if f == 4 {
		return 0, nil
	} else if f == 3 {
		return 1, nil
	} else if f == 2 {
		if b.State[boxX][boxY-1] == 1 && b.State[boxX+1][boxY] == 1 {
			return 2, nil
		} else if b.State[boxX][boxY+1] == 1 && b.State[boxX+1][boxY] == 1 {
			return 3, nil
		} else if b.State[boxX][boxY-1] == 1 && b.State[boxX-1][boxY] == 1 {
			return 4, nil
		} else if b.State[boxX][boxY+1] == 1 && b.State[boxX-1][boxY] == 1 {
			return 5, nil
		} else if b.State[boxX+1][boxY] == 1 && b.State[boxX-1][boxY] == 1 {
			return 6, nil
		} else if b.State[boxX][boxY+1] == 1 && b.State[boxX][boxY-1] == 1 {
			return 7, nil
		} else {
			return -1, fmt.Errorf("出错，自由度为2而监测不到类型 boxX:%d,boxY:%d", boxX, boxY)
		}
	} else if f == 1 {
		return 8, nil
	} else {
		return 9, nil
	}

}

// CheckChainType 执行此方法会设置chain类型，基本必调用
func (c *Chain) CheckChainType() error {
	if len(c.Boxes) == 1 {
		c.Type = 1
		return nil
	} else if len(c.Boxes) == 2 {
		c.Type = 2
		return nil
	} else if len(c.Boxes) > 2 {
		onePoint := c.Endpoint[0]
		twoPoint := c.Endpoint[1]
		absX, absY := math.Abs(float64(onePoint.X-twoPoint.X)), math.Abs(float64(onePoint.Y-twoPoint.Y))
		if absX+absY != 2 {
			//长链
			c.Type = 3
			return nil
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
				return fmt.Errorf("absX,absY均不等于2")
			}

		}

	}
	return nil
}

// GetAllMoves 得到所有可下边
func (b *Board) GetAllMoves() (edges []*Edge, err error) {
	for i := 0; i < 11; i++ {
		for j := 0; j < 11; j++ {
			if (i+j)&1 == 1 && b.State[i][j] != 1 {
				edges = append(edges, &Edge{i, j})
			}
		}
	}
	return
}

// Get2FEdge 获取移动后不会被捕获的边
func (b *Board) Get2FEdge() (edges []*Edge, err error) {
	/*defer func() {
		fmt.Println(edges)
	}()*/
	//获取寻常边
	for i := 0; i < 11; i++ {
		for j := 1; j < 11; j++ { //正常11*11=121次 这里25次遍历,但是操作数基本一致

			if (i+j)&1 == 1 && b.State[i][j] == 0 {
				he := Edge{i, j}
				boxesF, err := b.GetFByE(&he)
				if err != nil {
					return nil, err
				}
				// 两边格子freedom大于3的边
				if (boxesF[0] >= 3 || boxesF[0] == -1) && (boxesF[1] >= 3 || boxesF[1] == -1) {
					edges = append(edges, &he)
				}
			}

		}

	}
	return
}

// GetEdgesByIdentifyingChains 通过识别链的方法获得可下边,执行前确保局面无死格
func (b *Board) GetEdgesByIdentifyingChains() (edges [][]*Edge, err error) {
	//获取寻常边
	twoFEdges, get2FEdgeErr := b.Get2FEdge()
	if get2FEdgeErr != nil {
		return nil, get2FEdgeErr
	}
	for _, e := range twoFEdges {
		edges = append(edges, []*Edge{e})
	}
	//edges = append(edges, twoFEdges...)

	//获取链边
	chains, err := b.GetChains()
	for _, chain := range chains {
		//fmt.Print(chain)
		//fmt.Println(b)
		//获取其中的一条
		boxX, boxY := chain.Endpoint[0].X, chain.Endpoint[0].Y
		//fmt.Println(boxX, boxY)
		for i := 0; i < 4; i++ {
			edgeX, edgeY := boxX+d1[i][0], boxY+d1[i][1]
			//fmt.Println(edgeX, edgeY)
			if b.State[edgeX][edgeY] == 0 {

				edges = append(edges, []*Edge{{edgeX, edgeY}})
				break
			}
		}

	}

	//获取死树
	edge, err := b.GetDTreeEdges()
	if err != nil {
		return
	}
	edges = append(edges, edge...)
	return edges, nil
}

// GetDGridEdges 获得死格的边
func (b *Board) GetDGridEdges() (edges []*Edge, err error) {
	edgesMark := make(map[string]bool)
	for i := 1; i < 11; i += 2 {
		for j := 1; j < 11; j += 2 {
			f, getFByBIErr := b.GetFByBI(i, j)
			if err != nil {
				return nil, getFByBIErr
			}
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
						f1, getFByBIErr2 := b.GetFByBI(boxX, boxY)
						if getFByBIErr2 != nil {
							return nil, getFByBIErr2
						}
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

// GetDTreeEdges 获得死树的边，务必保证调用此方法前局面已经没有死格
func (b *Board) GetDTreeEdges() (edges [][]*Edge, err error) {

	boxesMark := map[int]bool{}
	for i := 1; i < 11; i += 2 {
		for j := 1; j < 11; j += 2 {
			f, getFByBIErr := b.GetFByBI(i, j)
			if err != nil {
				return nil, getFByBIErr
			}
			if f == 1 {
				//先判断是不是死环
				if is, err := b.IsDCircle(i, j); err != nil {
					return nil, err
				} else if is > 0 {
					//全捕获
					if es, err := b.GetDCircleEdges(i, j, is-1, false); err != nil {
						return nil, err
					} else {
						edges = append(edges, es)
					}
					//双交
					if es, err := b.GetDCircleEdges(i, j, is-4, true); err != nil {
						return nil, err
					} else {
						edges = append(edges, es)
					}
					continue
				} else if is < 0 {
					//这种情况是两头都为1自由度,但是不是死环,可是是与死环一样的策略，因为本质相同
					//全捕获
					is *= -1
					if es, err := b.GetDCircleEdges(i, j, is-1, false); err != nil {
						return nil, err
					} else {
						edges = append(edges, es)
					}
					//双交
					if es, err := b.GetDCircleEdges(i, j, is-4, true); err != nil {
						return nil, err
					} else {
						edges = append(edges, es)
					}
					continue

				}

				//死链
				for k := 0; k < 4; k++ {
					edgeX := i + d1[k][0]
					edgeY := j + d1[k][1]
					if b.State[edgeX][edgeY] == 0 {
						boxX := i + d2[k][0]
						boxY := j + d2[k][1]
						f1, getFByBIErr2 := b.GetFByBI(boxX, boxY)
						if getFByBIErr2 != nil {
							return nil, getFByBIErr2
						}

						if f1 == 2 {
							chain := NewChain()
							getChainErr := b.GetChain(boxX, boxY, boxesMark, chain, true)
							if getChainErr != nil {
								return nil, getChainErr
							}
							//全捕获
							if es, err := b.GetDChainEdges(i, j, chain, chain.Length, false); err != nil {
								return nil, err
							} else {
								edges = append(edges, es)
							}
							//双交
							if es, err := b.GetDChainEdges(i, j, chain, chain.Length-2, true); err != nil {
								return nil, err
							} else {
								edges = append(edges, es)
							}
						}

					}
				}

			}
		}
	}
	return
}

// GetChains 获得链集合
func (b *Board) GetChains() (chains []*Chain, err error) {
	var boxesMark map[int]bool
	boxesMark = map[int]bool{}
	for i := 1; i < 11; i += 2 {
		for j := 1; j < 11; j += 2 {
			x, y, boxToXYErr := BoxToXY(i, j)
			if boxToXYErr != nil {
				return nil, boxToXYErr
			}
			index := x*5 + y
			//如果访问过
			if boxesMark[index] {
				continue
			}

			f, getFByBIErr := b.GetFByBI(i, j)
			if getFByBIErr != nil {
				return nil, getFByBIErr
			}
			if f == 2 {
				chain := NewChain()
				b.Boxes[(i/2)*5+j/2].Type, err = b.GetBoxType(i, j)
				if err != nil {
					return nil, err
				}
				if getChainErr := b.GetChain(i, j, boxesMark, chain, true); getChainErr != nil {
					return nil, getChainErr
				}
				if err = chain.CheckChainType(); err != nil {
					return nil, err
				}

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
func (b *Board) GetChain(boxX, boxY int, boxesMark map[int]bool, chain *Chain, isStart bool) error {
	if boxX < 1 || boxX > 10 || boxY < 1 || boxY > 10 {
		return fmt.Errorf("非格子坐标")
	}
	x, y, _ := BoxToXY(boxX, boxY)
	//flag
	flag := false
	index := x*5 + y
	t := b.Boxes[index].Type
	//"F4","F3","|_", "_|", "|￣", "￣|", " 二 ", "| |", "F1", "F0"
	if t < 2 || t > 7 {
		return fmt.Errorf("此类型不可操作:%d", t)
	}
	if boxesMark[index] {
		return nil
	}
	//加入链
	chain.Boxes = append(chain.Boxes, b.Boxes[index])
	chain.Length++
	//标记
	boxesMark[index] = true

	//先由类型得到可走的方向
	for _, aD := range bToAs[t-2] {
		nextBoxX, nextBoxY := boxX+d2[aD][0], boxY+d2[aD][1]
		if nextBoxX > 0 && nextBoxX < 11 && nextBoxY > 0 && nextBoxY < 11 {

			nextX, nextY, err := BoxToXY(nextBoxX, nextBoxY)
			if err != nil {
				return err
			}
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
					getChainErr := b.GetChain(nextBoxX, nextBoxY, boxesMark, chain, false)
					if getChainErr != nil {
						return getChainErr
					}
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

	return nil

}

// GetFByBI 通过格子下标获得格子自由度
func (b *Board) GetFByBI(boxI, boxJ int) (int, error) {
	if boxI&1 != 1 || boxJ&1 != 1 {
		return 0, fmt.Errorf("boxIndex Error")
	} else if boxI <= 0 || boxI >= 10 || boxJ <= 0 || boxJ >= 10 {
		return -1, nil
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
		return freeDom, nil
	} else {
		return 0, nil
	}

}

// GetEdgeByBI 通过格子下标获得格子所有边
func (b *Board) GetEdgeByBI(boxI, boxJ int) (edges []*Edge, err error) {
	f, GetFByBIErr := b.GetFByBI(boxI, boxJ)
	if GetFByBIErr != nil {
		return nil, GetFByBIErr
	}
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
	return edges, nil
}

// EatAllCBox 吃完所有C型格
func (b *Board) EatAllCBox() error {

	for {
		flag := false
		for i := 1; i < 11; i = i + 2 {
			for j := 1; j < 11; j = j + 2 {
				f, err := b.GetFByBI(i, j)
				if err != nil {
					return err
				}
				//自由度为一，占领
				if f == 1 {
					flag = true
					edge, getEdgeBy1FBIErr := b.GetEdgeByBI(i, j)
					if getEdgeBy1FBIErr != nil {
						return getEdgeBy1FBIErr
					}
					if moveIJERR := b.Move(edge...); moveIJERR != nil {
						return moveIJERR
					}
					//设置占领者颜色与分数
					if edge[0].X&1 == 1 {
						//竖
						if j-2 >= 0 {
							//[]|*左边   *是当前格子
							f2, GetFByBIERR := b.GetFByBI(i, j-2)
							if GetFByBIERR != nil {
								return GetFByBIERR
							}
							if f2 == 0 && b.State[i][j-2] == 0 {
								b.State[i][j-2] = b.Now
								b.S[b.Now]++
							}
						}
						if j+2 < 11 {
							//*|[]右边   *是当前格子
							f2, GetFByBIERR := b.GetFByBI(i, j+2)
							if GetFByBIERR != nil {
								return GetFByBIERR
							}
							if f2 == 0 && b.State[i][j+2] == 0 {
								b.State[i][j+2] = b.Now
								b.S[b.Now]++
							}
						}

					} else {
						//横
						if i-2 >= 0 {
							//[]
							//——       *是当前格子
							//*
							//下边
							f2, GetFByBIERR := b.GetFByBI(i-2, j)
							if GetFByBIERR != nil {
								return GetFByBIERR
							}
							if f2 == 0 && b.State[i-2][j] == 0 {
								b.State[i-2][j] = b.Now
								b.S[b.Now]++
							}

						}
						if i+2 < 11 {
							// *
							//——       *是当前格子
							//[]
							//下边
							f2, GetFByBIERR := b.GetFByBI(i+2, j)
							if GetFByBIERR != nil {
								return GetFByBIERR
							}
							if f2 == 0 && b.State[i+2][j] == 0 {
								b.State[i+2][j] = b.Now
								b.S[b.Now]++
							}
						}
					}
					b.State[i][j] = b.Now
					b.S[b.Now]++
				}
			}
		}
		if !flag {
			break
		}
	}

	return nil
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

// RandomMove 随机移动,目前为getAllMoves,不带checkout
func (b *Board) RandomMove() (edge *Edge, err error) {
	edges, getAllMovesErr := b.GetAllMoves()
	if getAllMovesErr != nil {
		return nil, getAllMovesErr
	}
	if len(edges) == 0 {
		return nil, fmt.Errorf("没有可移动的边")
	}
	randInt := rand.Intn(len(edges))
	if err = b.Move(edges[randInt]); err != nil {
		return nil, err
	}
	return edges[randInt], nil
}

// RandomMoveByCheck 随机移动,目前为GetDGridEdges()后GetEdgesByIdentifyingChains,自带checkout
func (b *Board) RandomMoveByCheck() (edge []*Edge, err error) {
	ees, err := b.GetMove()
	if err != nil {
		return nil, err
	}
	if len(ees) == 0 {
		return
	}
	randInt := rand.Intn(len(ees))
	if err = b.MoveAndCheckout(ees[randInt]...); err != nil {
		return nil, err
	}
	return ees[randInt], nil
}

// GetFByE 返回边两边的freedom ,默认 左右，上下的顺序，若在边上则对应位置为-1
func (b *Board) GetFByE(edge *Edge) (boxesF []int, err error) {
	d := [2]int{-1, 1}
	if edge.X&1 == 1 {
		//竖边
		for _, v := range d {
			boxX := edge.X
			boxY := edge.Y + v
			if boxY < 11 && boxY >= 0 {
				f, err := b.GetFByBI(boxX, boxY)
				if err != nil {
					return nil, err
				}
				boxesF = append(boxesF, f)
			} else {
				boxesF = append(boxesF, -1)

			}

		}

	} else {
		//横边
		for _, v := range d {
			boxX := edge.X + v
			boxY := edge.Y
			if boxX < 11 && boxX >= 0 {
				f, err := b.GetFByBI(boxX, boxY)
				if err != nil {
					return nil, err
				}
				boxesF = append(boxesF, f)
			} else {
				boxesF = append(boxesF, -1)

			}

		}
	}
	return
}

// IsDCircle 格子freedom为一时才可调用
func (b *Board) IsDCircle(boxX, boxY int) (is int, err error) {
	edgesMark := map[string]bool{}
	if is, err = b.dfsIsDCircle(boxX, boxY, boxX, boxY, 1, edgesMark); err != nil {
		return 0, err
	} else {
		return is, nil
	}
}
func (b *Board) dfsIsDCircle(sBoxX, sBoxY, boxX, boxY, len int, edgesMark map[string]bool) (is int, err error) {
	if !IsBox(boxX, boxY) {
		return 0, fmt.Errorf("不是格子下标")
	}
	for i := 0; i < 4; i++ {
		nEX, nEY := boxX+d1[i][0], boxY+d1[i][1]
		nBX, nBY := boxX+d2[i][0], boxY+d2[i][1]
		edge := &Edge{nEX, nEY}
		if b.State[nEX][nEY] == 0 && !edgesMark[edge.String()] {
			edgesMark[edge.String()] = true
			if f, err := b.GetFByBI(nBX, nBY); err != nil {
				return 0, err
			} else if f == 1 {
				ans := math.Abs(float64(sBoxX-nBX)) + math.Abs(float64(sBoxY-nBY))
				if ans == 2 {
					return len + 1, nil
				} else {
					return -1 * (len + 1), nil //特殊情况，一般不会有，但是出现了就处理一下
				}
			} else if f == 2 {
				if is, err = b.dfsIsDCircle(sBoxX, sBoxY, nBX, nBY, len+1, edgesMark); err != nil {
					return 0, err
				} else {
					return is, nil
				}
			}
		}
	}
	return 0, nil
}

// GetDCircleEdges 格子freedom为一时才可调用
func (b *Board) GetDCircleEdges(boxX, boxY, len int, isDoubleCross bool) (edges []*Edge, err error) {
	edgesMark := map[string]bool{}
	if err = b.dfsGetDCircleEdges(boxX, boxY, len, edgesMark, &edges, isDoubleCross); err != nil {
		return nil, err
	}
	return edges, err

}
func (b *Board) dfsGetDCircleEdges(boxX, boxY, len int, edgesMark map[string]bool, edges *[]*Edge, isDoubleCross bool) (err error) {
	if !IsBox(boxX, boxY) {
		return fmt.Errorf("不是格子下标")
	}

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
			if err = b.dfsGetDCircleEdges(nBX, nBY, len-1, edgesMark, edges, isDoubleCross); err != nil {
				return err
			}
		}
	}
	return nil
}

// GetDChainEdges 获得死链边
func (b *Board) GetDChainEdges(box1FX, box1FY int, c *Chain, len int, isDoubleCross bool) (edges []*Edge, err error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	edgesMark := map[string]bool{}
	if err = b.dfsChainEdges(box1FX, box1FY, edgesMark, len, &edges); err != nil {
		return
	} else if isDoubleCross {
		endPointX, endPointY := -1, -1
		for k := 0; k < 4; k++ {
			edgeX, edgeY := box1FX+d1[k][0], box1FY+d1[k][1]
			nextBoxX, nextBoxY := box1FX+d2[k][0], box1FY+d2[k][1]
			if b.State[edgeX][edgeY] == 0 {
				if nextBoxX == c.Endpoint[0].X && nextBoxY == c.Endpoint[0].Y {
					endPointX, endPointY = c.Endpoint[1].X, c.Endpoint[1].Y
				} else if nextBoxX == c.Endpoint[1].X && nextBoxY == c.Endpoint[1].Y {
					endPointX, endPointY = c.Endpoint[0].X, c.Endpoint[0].Y
				} else {
					fmt.Println(b, c, nextBoxX, nextBoxY)
					chains, err := b.GetChains()
					if err != nil {
						return nil, err
					}
					for _, c := range chains {
						fmt.Println(c)
					}
					return nil, fmt.Errorf("校对端点失败")
				}
				break
			}
		}
		for k := 0; k < 4; k++ {
			edgeX, edgeY := endPointX+d1[k][0], endPointY+d1[k][1]
			nextBoxX, nextBoxY := endPointX+d2[k][0], endPointY+d2[k][1]
			f, err := b.GetFByBI(nextBoxX, nextBoxY)
			if err != nil {
				return nil, err
			}
			if b.State[edgeX][edgeY] == 0 && f != 2 && f != 1 {
				edges = append(edges, &Edge{edgeX, edgeY})
				break
			}
		}
	}
	return
}
func (b *Board) dfsChainEdges(sBoxX, sBoxY int, edgesMark map[string]bool, len int, edges *[]*Edge) (err error) {

	if len >= 0 {
		for k := 0; k < 4; k++ {
			edgeX, edgeY := sBoxX+d1[k][0], sBoxY+d1[k][1]
			nextBoxX, nextBoxY := sBoxX+d2[k][0], sBoxY+d2[k][1]
			edge := &Edge{edgeX, edgeY}
			if b.State[edgeX][edgeY] == 0 && !edgesMark[edge.String()] && nextBoxX >= 0 && nextBoxX <= 10 && nextBoxY >= 0 && nextBoxY <= 10 {
				*edges = append(*edges, edge)
				edgesMark[edge.String()] = true
				len--
				if err = b.dfsChainEdges(nextBoxX, nextBoxY, edgesMark, len, edges); err != nil {
					return
				}
			}
		}

	}
	return
}
