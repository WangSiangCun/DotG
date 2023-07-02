package uct

import (
	"dotg/board"
	"fmt"
	"math"
	"math/rand"
	"time"
)

type UCT struct {
}
type UCTNode struct {
	B        *board.Board
	Children []*UCTNode
	Parents  *UCTNode
	Visit    int
	Win      int
	UCBScore int
	TriedMap map[string]bool
	LastMove []*board.Edge
}

func NewUCTNode(b *board.Board) *UCTNode {
	return &UCTNode{
		B:        b,
		Children: []*UCTNode{},
		Parents:  nil,
		Visit:    0,
		Win:      0,
		UCBScore: 0,
		TriedMap: map[string]bool{},
		LastMove: []*board.Edge{},
	}
}

const (
	ucb_C float64 = 0.4142135623730951
)

func (n *UCTNode) GetUCB() float64 {
	if n.Visit == 0 {
		return rand.Float64() + 1.0
	}
	return float64(n.Win)/float64(n.Visit) + ucb_C*math.Sqrt(math.Log(float64(n.Parents.Visit))/float64(n.Visit))
}
func (n *UCTNode) AddChild(child *UCTNode) {
	n.Children = append(n.Children, child)
}
func (n *UCTNode) GetUnTriedEdges() (edges [][]*board.Edge, err error) {

	if ees, err := n.B.GetMove(); err != nil {
		return nil, err
	} else {
		for _, es := range ees {
			s := fmt.Sprintf("%v", es)
			if !n.TriedMap[s] {
				edges = append(edges, es)
			}
		}
	}
	return
}
func UCTSearch(b *board.Board, timeout int, iter, who int) ([]*board.Edge, error) {
	root := NewUCTNode(b)
	startT := time.Now()
	for i := 0; int(time.Since(startT).Milliseconds()) < timeout || i < iter; i++ {
		nowN := root
		nowN, err := SelectBest(nowN)
		if err != nil {
			return nil, err
		}
		nB := board.CopyBoard(nowN.B)
		res, err := Simulation(nB, who)
		fmt.Println(nB)
		if err != nil {
			return nil, err
		}
		BackUp(nowN, res)
	}
	bestChild := GetBestChild(root, true)
	fmt.Println(bestChild.B)
	return bestChild.LastMove, nil
}

func GetBestChild(n *UCTNode, isEnd bool) *UCTNode {

	var bestN *UCTNode
	var bestUCB float64
	bestUCB = math.MinInt32
	for _, child := range n.Children {
		cUCB := child.GetUCB()
		if isEnd {
			fmt.Print(cUCB, "|", child.Win/child.Visit, " ")
		}
		if cUCB > bestUCB {
			bestUCB = cUCB
			bestN = child
		}
	}

	return bestN
}
func Simulation(b *board.Board, who int) (res int, err error) {
	for b.Status() == 0 {
		if _, err := b.RandomMoveByCheck(); err != nil {
			return 0, err
		}
	}
	e := b.Status()
	if e == who {
		return 1, err
	} else {
		return 0, err
	}

}
func SelectBest(n *UCTNode) (*UCTNode, error) {
	//如果游戏已经结束
	if n.B.Status() != 0 {
		return n, nil
	}
	if ees, err := n.GetUnTriedEdges(); err != nil {
		return nil, err
		//没有可以扩展的子结点,选择ucb值最大的子结点继续
	} else if len(ees) == 0 {
		n = GetBestChild(n, false)
		n, err = SelectBest(n)
	} else {
		if n, err = Expand(&ees, n); err != nil {
			return nil, err
		} else {
			return n, nil
		}
	}
	return n, nil
}
func Expand(edges *[][]*board.Edge, n *UCTNode) (*UCTNode, error) {
	//随机选一个扩展
	l := len(*edges)
	randInt := rand.Intn(l)
	es := (*edges)[randInt]
	nB := board.CopyBoard(n.B)
	if err := nB.MoveAndCheckout(es...); err != nil {
		return nil, err
	}
	//fmt.Println(n.B, nB)
	//生产新结点
	nN := NewUCTNode(nB)
	nN.Parents = n
	nN.LastMove = es
	n.TriedMap[fmt.Sprintf("%v", es)] = true
	n.Children = append(n.Children, nN)
	return nN, nil

}
func BackUp(n *UCTNode, res int) {
	for n != nil {
		n.Win += res
		n.Visit++
		n = n.Parents
	}
}
