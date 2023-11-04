package uct

import (
	"dotg/board"
	"dotg/record"
	"fmt"
	"math"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"
)

type UCT struct {
}
type UCTNode struct {
	B           *board.Board
	Children    []*UCTNode
	Parents     *UCTNode
	Visit       int
	Win         int //0 无，1 ，2
	UnTriedMove [][]*board.Edge
	LastMove    []*board.Edge
	rwMutex     sync.RWMutex
}

type HashKey struct {
	M   [2]uint64
	Now int
}

type HashValue struct {
	Visit, Win int
	Turn       int
}

var (
	C         float64       = 1.4142135623730951
	ThreadNum int           = 4
	MaxDeep   int           = 0
	MaxChild  int           = 16
	TimeLimit int           = 10
	sumTime   time.Duration = 0
)

func init() {
	rand.Seed(time.Now().Unix())
	ThreadNum = runtime.NumCPU()
}
func (n *UCTNode) GetUCB() float64 {
	if n.Visit == 0 {
		return rand.Float64() + 1.0
	}
	return float64(n.Win)/float64(n.Visit) + C*math.Sqrt(math.Log(float64(n.Parents.Visit))/float64(n.Visit))
}
func Move(b *board.Board, mode int, isV bool) []*board.Edge {
	start := time.Now()
	es := []*board.Edge{}
	//固定先手开局
	if b.Turn == 0 {
		es = append(es, &board.Edge{4, 5})
	} else {
		ees := b.GetFrontMoveByTurn()
		if ees != nil {
			es = Search(b, mode, isV)
		} else if ees == nil {
			es = b.GetEndMove()
		}
	}
	//先记录
	record.PrintContentMiddle(b, es)

	b.MoveAndCheckout(es...)
	if isV {
		fmt.Println(es)
		theTime := time.Since(start)
		sumTime += theTime
		fmt.Printf("%v\n本次花费时间：%v,总耗时:%v\n", b, theTime, sumTime)
		fmt.Println("-------------------------")
	}

	return es
}

func (n *UCTNode) BackUp(res int) {
	if n.B.Now == res {
		n.Win += 1
	} else {
		n.Win += 0
	}
	n.Visit++
}
func NewUCTNode(b *board.Board) *UCTNode {
	return &UCTNode{
		B:           b,
		Children:    []*UCTNode{},
		Parents:     nil,
		Visit:       0,
		Win:         0,
		UnTriedMove: [][]*board.Edge{},
		LastMove:    []*board.Edge{},
		rwMutex:     sync.RWMutex{},
	}
}
func Simulation(b *board.Board) (res int) {
	//nB仅仅用于模拟
	nB := board.CopyBoard(b)
	for nB.Status() == 0 {
		nB.RandomMoveByCheck()
	}
	//fmt.Println(nB)
	return nB.Status()

}
func GetBestChild(n *UCTNode, isV bool) *UCTNode {
	//如果游戏已经结束
	var bestN *UCTNode
	var bestUCB float64
	bestUCB = math.MinInt32
	for i := 0; i < len(n.Children); i++ {
		cUCB := n.Children[i].GetUCB()
		if isV {
			fmt.Printf("ucb: %v   w/v: %v v:%v move:%v \n", cUCB, float64(n.Children[i].Win)/float64(n.Children[i].Visit), n.Children[i].Visit, n.Children[i].LastMove)
		}
		if cUCB > bestUCB {
			bestUCB = cUCB
			bestN = n.Children[i]
		}
	}
	if isV {
		fmt.Printf("Select:\n UCB:%.4f  w/v: %.4f Child: %d C: %v\n", bestUCB, float64(bestN.Win)/float64(bestN.Visit), len(n.Children), C)
	}
	return bestN
}
func GetBestChildByMV(n *UCTNode, isV bool) *UCTNode {
	//如果游戏已经结束
	var bestN *UCTNode
	var bestWV float64
	bestWV = math.MinInt32
	for i := 0; i < len(n.Children); i++ {
		wv := float64(n.Children[i].Win) / float64(n.Children[i].Visit)
		if isV {
			fmt.Printf("  w/v: %v v:%v move:%v \n", float64(n.Children[i].Win)/float64(n.Children[i].Visit), n.Children[i].Visit, n.Children[i].LastMove)
		}
		if wv > bestWV {
			bestWV = wv
			bestN = n.Children[i]
		}
	}
	if isV {
		fmt.Printf("Select:\n  w/v: %.4f Child: %d C: %v\n", float64(bestN.Win)/float64(bestN.Visit), len(n.Children), C)
	}
	return bestN
}
func SelectBest(n *UCTNode) (next *UCTNode) {
	n.rwMutex.RLock()
	defer n.rwMutex.RUnlock()
	if n.Parents != nil {
		n.Parents.rwMutex.RLock()
		defer n.Parents.rwMutex.RUnlock()
	}
	//游戏结束
	if n.B.Status() != 0 {
		return nil
	}
	//获取还没尝试的边
	if len(n.UnTriedMove) == 0 {
		//Untried==0 children!=0 属于扩展完全
		var bestN *UCTNode
		var bestUCB float64
		bestUCB = math.MinInt32
		for i := 0; i < len(n.Children); i++ {
			cUCB := n.Children[i].GetUCB()
			if cUCB > bestUCB {
				bestUCB = cUCB
				bestN = n.Children[i]
			}
		}
		return bestN
	} else {
		//Untried!=0 children!=0 属于扩展不完全
		//还有可扩展
		return nil
	}

}
func Expand(n *UCTNode) *UCTNode {
	//routine1的Select1 进行选择，此时未扩展完，而routine2的select2因为同时和select1读到未扩展完，
	//而2或1的expand先扩展，另一个堵塞后再去扩展发现已经被扩展，
	//这时候就会出现问题
	n.rwMutex.Lock()
	defer n.rwMutex.Unlock()
	if n.Parents != nil {
		n.Parents.rwMutex.Lock()
		defer n.Parents.rwMutex.Unlock()
	}

	if len(n.UnTriedMove) != 0 {
		///已扩展，未扩展完毕

		//生产新结点
		es := n.UnTriedMove[0]
		nB := board.CopyBoard(n.B)
		nB.MoveAndCheckout(es...)

		nN := NewUCTNode(nB)
		nN.Parents = n
		nN.LastMove = es
		//初始化后，将nN加入n的子节点
		n.Children = append(n.Children, nN)
		//fmt.Println(len(n.Children))
		n.UnTriedMove = n.UnTriedMove[1:]

		//fmt.Println(n.UnTriedMove)
		return nN

	} else if len(n.UnTriedMove) == 0 && len(n.Children) == 0 {
		//未扩展
		ees := n.B.GetMove()
		if len(ees) == 0 {
			return n
		}

		//	maxL := min(len(ees), MaxChild)
		//打乱
		//Shuffle(ees)
		//只要前maxL个
		n.UnTriedMove = ees
		//fmt.Println(ees)
		//fmt.Println(n.UnTriedMove)
		//生产新结点
		es := n.UnTriedMove[0]
		nB := board.CopyBoard(n.B)
		nB.MoveAndCheckout(es...)

		nN := NewUCTNode(nB)
		nN.Parents = n
		nN.LastMove = n.UnTriedMove[0]
		//初始化后，将nN加入n的子节点
		n.Children = append(n.Children, nN)
		n.UnTriedMove = n.UnTriedMove[1:]
		return nN
	} else {
		//避免并发问题
		//fmt.Println(n.UnTriedMove, n.B.Status())
		return n
	}
}
func Search(b *board.Board, mode int, isV bool) (es []*board.Edge) {
	if b.Turn <= 1 {
		runtime.GC()
	}
	var (
		exit = make(chan int, ThreadNum)
		stop = make(chan int, ThreadNum)
	)
	MaxDeep = 0
	root := NewUCTNode(b)

	start := time.Now()
	res := 0
	AdjustUCB(b)
	AdjustMaxChild(b)
	AdjustTimeLimit(b, mode)

	for i := 0; i < ThreadNum; i++ {
		go func() {

			for len(stop) == 0 {

				if int(time.Since(start).Seconds()) >= TimeLimit {
					stop <- 1
				}

				nowN := root
				deep := 0

				//选择节点，如果该节点没有扩展完全或者游戏结束则返回nil，否则继续选择
				for next := SelectBest(nowN); next != nil; {
					nowN = next
					next = SelectBest(nowN)
					deep++
				}
				if deep > MaxDeep {
					MaxDeep = deep
				}
				if nowN.B.Status() == 0 {
					//扩展
					nowN = Expand(nowN)

					res = Simulation(nowN.B)
				} else {
					res = nowN.B.Status()
				}

				for nowN != nil {
					nowN.BackUp(res)
					nowN = nowN.Parents
				}
			}
			exit <- 1
		}()
	}

	for i := 0; i < ThreadNum; i++ {
		<-exit
	}
	bestChild := GetBestChildByMV(root, isV)
	if isV {
		fmt.Printf("Tatal:%d \n MaxDeep:%d\n SimRate:%v\n", root.Visit, MaxDeep, float64(bestChild.Visit)/float64(root.Visit))
		//file, err := os.OpenFile("uctNodeTree.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		//if err != nil {
		//	fmt.Println("Error opening file:", err)
		//	return
		//}

		//fmt.Fprintf(file, "\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n")
		//printTree(root, 0, file)
		//file.Close()
	}
	return bestChild.LastMove
}
func AdjustUCB(b *board.Board) {

	C = math.Sqrt(2.0) * 0.60

}
func AdjustMaxChild(b *board.Board) {
	switch {
	case b.Turn <= 11:
		MaxChild = 16
	case b.Turn <= 13:
		MaxChild = 18
	case b.Turn <= 15:
		MaxChild = 20
	default:
		MaxChild = 22
	}
}
func AdjustTimeLimit(b *board.Board, mode int) {
	if mode == 0 {
		switch {
		case b.Turn <= 7:
			TimeLimit = 27
		case b.Turn <= 10:
			TimeLimit = 35
		case b.Turn <= 15:
			TimeLimit = 45
		case b.Turn <= 20:
			TimeLimit = 60
		case b.Turn <= 25:
			TimeLimit = 30
		default:
			TimeLimit = 10
		}
	} else if mode == 1 {
		switch {
		case b.Turn <= 7:
			TimeLimit = 15
		case b.Turn <= 10:
			TimeLimit = 20
		case b.Turn <= 15:
			TimeLimit = 30
		case b.Turn <= 20:
			TimeLimit = 40
		case b.Turn <= 25:
			TimeLimit = 30
		default:
			TimeLimit = 10
		}
	} else if mode == 2 {
		TimeLimit = 20
	} else if mode == 3 {
		TimeLimit = 2
	} else if mode == 4 {
		TimeLimit = 10
	}

}
func min(a, b int) int {
	if a > b {
		return b
	} else {
		return a
	}
}

func printTree(node *UCTNode, depth int, writer *os.File) {
	fmt.Fprintf(writer, "%v %v %v:%v/%v es: %v\n", depth, getIndent(depth), node.B.Now, node.Win, node.Visit, node.LastMove)

	for _, child := range node.Children {
		printTree(child, depth+1, writer)
	}
}

func getIndent(depth int) string {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "\t"
	}
	return indent
}
func Shuffle(arr [][]*board.Edge) [][]*board.Edge {
	for i, j := range rand.Perm(len(arr)) {
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}
