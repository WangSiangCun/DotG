package uct

import (
	"dotg/board"
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sort"
	"strconv"
	"sync"
	"time"
)

type UCT struct {
}
type UCTNode struct {
	B           *board.Board
	Children    []*UCTNode
	Parents     *UCTNode
	Visit, Now  int
	Win         int //0 无，1 ，2
	UnTriedMove []Untry
	LastMove    int64
	rwMutex     sync.RWMutex
}
type Untry struct {
	m   int64
	val float64
}

type ByX []Untry

func (self ByX) Len() int           { return len(self) }
func (self ByX) Swap(i, j int)      { self[i], self[j] = self[j], self[i] }
func (self ByX) Less(i, j int) bool { return self[i].val > self[j].val }

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
	maxDeep   int           = 0
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
func Move(b *board.Board, who int, isV bool, isHeuristic bool) []*board.Edge {
	start := time.Now()
	es := []*board.Edge{}
	//固定先手开局
	if b.Turn == 0 {
		es = append(es, &board.Edge{4, 5})
	} else {
		ees := b.GetFrontMoveByTurn()
		if ees != nil {
			es = Search(b, who, isV, isHeuristic)
		} else if ees == nil {
			es = b.GetEndMove()
		}
	}

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

func (n *UCTNode) BackUp(res int, who int) {
	if n.B.Now == who {
		n.Win += res
	} else {
		n.Win += 1 - res
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
		UnTriedMove: ByX{},
		LastMove:    int64(0),
		rwMutex:     sync.RWMutex{},
	}
}
func Simulation(b *board.Board, who int) (res int) {
	for b.Status() == 0 {
		b.RandomMoveByCheck()
	}
	e := b.Status()
	if e == who {
		return 1
	} else {
		return 0
	}

}
func GetBestChild(n *UCTNode, isV bool) *UCTNode {
	//n.rwMutex.Lock()
	//defer n.rwMutex.Unlock()
	//如果游戏已经结束
	var bestN *UCTNode
	var bestUCB float64
	bestUCB = math.MinInt32
	for i := 0; i < len(n.Children); i++ {
		cUCB := n.Children[i].GetUCB()
		if isV {
			fmt.Printf("ucb: %v   w/v: %v v:%v move:%v \n", cUCB, float64(n.Children[i].Win)/float64(n.Children[i].Visit), n.Children[i].Visit, board.MtoEdges(n.Children[i].LastMove))
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
	if len(n.UnTriedMove) == 0 && len(n.Children) == 0 {
		//Untried==0 children==0还没开始扩展，比如root
		return nil
	} else if len(n.UnTriedMove) == 0 {
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
func Expand(n *UCTNode, isHeuristic bool) *UCTNode {
	//routine1的Select1 进行选择，此时未扩展完，而routine2的select2因为同时和select1读到未扩展完，
	//而2或1的expand先扩展，另一个堵塞后再去扩展发现已经被扩展，
	//这时候就会出现问题
	n.rwMutex.Lock()
	defer n.rwMutex.Unlock()
	if n.Parents != nil {
		n.Parents.rwMutex.Lock()
		defer n.Parents.rwMutex.Unlock()
	}
	if n.Visit < 50 {
		return n
	}
	if len(n.UnTriedMove) == 0 && len(n.Children) != 0 {
		if n.B.Status() != 0 {
			return n
		}
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
		n = bestN

	}

	if len(n.UnTriedMove) == 0 && len(n.Children) == 0 {
		ees := n.B.GetMove()
		maxL := min(len(ees), MaxChild)
		n.UnTriedMove = make([]Untry, maxL)
		for i := 0; i < maxL; i++ {
			n.UnTriedMove[i].m = board.EdgesToM(ees[i]...)
		}

		if isHeuristic {
			rew := map[string]float64{}
			if n.Parents != nil {
				for i, _ := range n.Parents.Children {
					if n.Parents.Children[i].Visit > 0 {
						rew[strconv.FormatInt(n.Parents.Children[i].LastMove, 10)] = (float64(n.Parents.Children[i].Win) / float64(n.Parents.Children[i].Visit)) + 1e-10
					}
				}
			}
			for i, un := range n.UnTriedMove {
				if rew[strconv.FormatInt(un.m, 10)] > 0 {
					n.UnTriedMove[i].val = rew[strconv.FormatInt(un.m, 10)]
				} else {
					n.UnTriedMove[i].val = 0.5 + rand.Float64()*1e-8
				}
			}
		} else {
			for i, _ := range n.UnTriedMove {
				{
					n.UnTriedMove[i].val = rand.Float64()
				}
			}
		}

		sort.Sort(ByX(n.UnTriedMove))
	}

	if len(n.UnTriedMove) == 0 {
		return n
	}
	es := board.MtoEdges(n.UnTriedMove[0].m)
	//fmt.Println(n.UnTriedMove)
	nB := board.CopyBoard(n.B)
	nB.MoveAndCheckout(es...)

	//生产新结点
	nN := NewUCTNode(nB)
	nN.Parents = n
	nN.LastMove = n.UnTriedMove[0].m

	if n.Children == nil {
		n.Children = make([]*UCTNode, 0, len(n.UnTriedMove))
	}
	n.Children = append(n.Children, nN)
	if len(n.UnTriedMove) > 1 {
		n.UnTriedMove = n.UnTriedMove[1:]
	} else {
		n.UnTriedMove = nil
	}
	if nN == nil {
		fmt.Println("nN为空：n:", n.B)
	}

	return nN

}
func Search(b *board.Board, who int, isV bool, isHeuristic bool) (es []*board.Edge) {
	if b.Turn <= 1 {
		runtime.GC()
	}
	var (
		exit = make(chan int, ThreadNum)
		stop = make(chan int, ThreadNum)
	)
	maxDeep = 0
	root := NewUCTNode(b)
	start := time.Now()
	res := 0
	AdjustUCB(b)
	AdjustMaxChild(b)
	AdjustTimeLimit(b)
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
				if deep > maxDeep {
					maxDeep = deep
				}
				if nowN.B.Status() == 0 {
					//扩展
					nowN = Expand(nowN, isHeuristic)
				}
				//nB仅仅用于模拟
				nB := board.CopyBoard(nowN.B)
				res = Simulation(nB, who)

				for nowN != nil {
					nowN.BackUp(res, who)
					nowN = nowN.Parents
				}
			}
			exit <- 1
		}()
	}

	for i := 0; i < ThreadNum; i++ {
		<-exit
	}
	bestChild := GetBestChild(root, isV)
	if isV {
		fmt.Printf("Tatal:%d \n MaxDeep:%d\n SimRate:%v\n", root.Visit, maxDeep, float64(bestChild.Visit)/float64(root.Visit))
	}

	return board.MtoEdges(bestChild.LastMove)
}
func AdjustUCB(b *board.Board) {
	switch {
	case b.Turn <= 11:
		C = math.Sqrt(2.0) * 1.00
	case b.Turn <= 13:
		C = math.Sqrt(2.0) * 0.80
	case b.Turn <= 15:
		C = math.Sqrt(2.0) * 0.70
	case b.Turn <= 17:
		C = math.Sqrt(2.0) * 0.60
	case b.Turn <= 19:
		C = math.Sqrt(2.0) * 0.55
	case b.Turn <= 23:
		C = math.Sqrt(2.0) * 0.50
	case b.Turn <= 27:
		C = math.Sqrt(2.0) * 0.40
	case b.Turn <= 31:
		C = math.Sqrt(2.0) * 0.30
	default:
		C = math.Sqrt(2.0) * 0.20
	}
}
func AdjustMaxChild(b *board.Board) {
	switch {
	case b.Turn <= 13:
		MaxChild = 18
	case b.Turn <= 16:
		MaxChild = 22
	default:
		MaxChild = 25
	}
}
func AdjustTimeLimit(b *board.Board) {
	switch {
	case b.Turn <= 7:
		TimeLimit = 15
	case b.Turn <= 10:
		TimeLimit = 20
	case b.Turn <= 15:
		TimeLimit = 30
	case b.Turn <= 20:
		TimeLimit = 60
	case b.Turn <= 25:
		TimeLimit = 30
	default:
		MaxChild = 10
	}
}
func min(a, b int) int {
	if a > b {
		return b
	} else {
		return a
	}
}
