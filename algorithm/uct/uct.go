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

const (
	ucb_C float64 = 0.4142135623730951

	MaxChild  int     = 25
	hashBlock uint    = 23 // prime
	hashSize  uint    = 13000000 / hashBlock
	INF       float64 = 1e100
)

var (
	ThreadNum int = 4
	maxDeep   int = 0
	rw        sync.RWMutex
	mutex     sync.Mutex
	hashTable [hashBlock]map[HashKey]*HashValue
	rwMutex   [hashBlock]sync.RWMutex
)

func (n *UCTNode) GetUCB() float64 {
	n.rwMutex.RLock()
	defer n.rwMutex.RUnlock()
	if n.Visit == 0 {
		return rand.Float64() + 1.0
	}
	return float64(n.Win)/float64(n.Visit) + ucb_C*math.Sqrt(math.Log(float64(n.Parents.Visit))/float64(n.Visit))
}
func Move(b *board.Board, timeout int, iter, who int, isV bool) []*board.Edge {
	start := time.Now()
	es := []*board.Edge{}
	if edges2F, err := b.Get2FEdge(); err != nil {
		fmt.Println(err)
		return nil
	} else if len(edges2F) != 0 {
		es, err = Search(b, timeout, iter, who, isV)
		if err != nil {
			fmt.Println(err)
			return nil
		}
	} else {
		ess, err := b.GetMove()
		if err != nil {
			fmt.Println(err)
			return nil
		}
		es = ess[0]
	}
	b.MoveAndCheckout(es...)
	fmt.Println(es)
	fmt.Println(b, time.Since(start))
	fmt.Println("-------------------------")
	return es
}
func init() {
	rand.Seed(time.Now().Unix())
	ThreadNum = runtime.NumCPU()
	for i := uint(0); i < hashBlock; i++ {
		hashTable[i] = make(map[HashKey]*HashValue, hashSize)
	}
}
func (n *UCTNode) BackUp(res int, who int) {
	n.rwMutex.Lock()
	defer n.rwMutex.Unlock()
	if n.Parents != nil {
		n.Parents.rwMutex.Lock()
		defer n.Parents.rwMutex.Unlock()
	}
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
func GetBestChild(n *UCTNode, isV bool) (*UCTNode, error) {
	n.rwMutex.Lock()
	defer n.rwMutex.Unlock()
	if n.Parents != nil {
		n.Parents.rwMutex.Lock()
		defer n.Parents.rwMutex.Unlock()
	}
	//如果游戏已经结束
	if n.B.Status() != 0 {
		return nil, fmt.Errorf(" GetBestChild:游戏已经结束")
	}
	var bestN *UCTNode
	var bestUCB float64
	bestUCB = math.MinInt32
	if len(n.Children) == 0 {
		return nil, fmt.Errorf("GetBestChild:错误，没有孩子结点")
	}
	for _, child := range n.Children {
		cUCB := child.GetUCB()
		if isV {
			fmt.Print("move:", child.LastMove, "ucb:", cUCB, "  w/v:", float64(child.Win)/float64(child.Visit), "  v:", child.Visit, "\n ")
		}
		if cUCB > bestUCB {
			bestUCB = cUCB
			bestN = child
		}
	}
	if bestN == nil {
		return nil, fmt.Errorf("未找到最好孩子结点")
	}
	if isV {
		fmt.Printf("Select:\n UCB:%.4f  w/v: %.4f Child: %d\n", bestUCB, float64(bestN.Win)/float64(bestN.Visit), len(n.Children))
	}
	return bestN, nil
}
func SelectBest(n *UCTNode) (next *UCTNode) {
	n.rwMutex.Lock()
	defer n.rwMutex.Unlock()
	if n.Parents != nil {
		n.Parents.rwMutex.Lock()
		defer n.Parents.rwMutex.Unlock()
	}
	for {
		if n.B.Status() != 0 {
			return nil
		}
		//select三种情况，一是untryMove
		//获取还没尝试的边
		if len(n.UnTriedMove) == 0 && len(n.Children) == 0 {
			return nil
		} else if len(n.UnTriedMove) == 0 {
			//Untried==0 而childreng不为0
			if n, err := GetBestChild(n, false); err != nil {
				fmt.Println(err)
				return nil
			} else {
				return n
			}
		} else {
			//还有可扩展
			return nil
		}

	}

}
func Expand(n *UCTNode) (*UCTNode, error) {
	n.rwMutex.Lock()
	defer n.rwMutex.Unlock()
	if n.Parents != nil {
		n.Parents.rwMutex.Lock()
		defer n.Parents.rwMutex.Unlock()
	}
	if n.B.Status() != 0 {
		return n, nil
	}
	if len(n.UnTriedMove) == 0 && len(n.Children) == 0 {
		if ees, err := n.B.GetMove(); err != nil {
			return nil, err
		} else {
			n.UnTriedMove = make([]Untry, len(ees))
			for i, es := range ees {
				board.MtoEdges(n.UnTriedMove[i].m)
				n.UnTriedMove[i].m = board.EdgesToM(es...)
			}
		}
	}

	if len(n.UnTriedMove) == 0 {
		return nil, fmt.Errorf("错误，没有扩展边")
	}

	rew := map[string]float64{}
	if n.Parents != nil {
		for _, c := range n.Parents.Children {
			if c.Visit > 0 {
				//fmt.Println(c.LastMove, strconv.FormatInt(c.LastMove, 10))
				rew[strconv.FormatInt(c.LastMove, 10)] = float64(c.Win)/float64(c.Visit) + 1e-10
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
	sort.Sort(ByX(n.UnTriedMove))
	//随机选一个扩展

	es := board.MtoEdges(n.UnTriedMove[0].m)
	nB := board.CopyBoard(n.B)
	if err := nB.MoveAndCheckout(es...); err != nil {
		return nil, err
	}

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
	return nN, nil

}
func Search(b *board.Board, timeoutSeconds int, iter, who int, isV bool) (es []*board.Edge, err error) {
	var (
		exit = make(chan int, ThreadNum)
		stop = make(chan int, ThreadNum)
	)
	maxDeep = 0
	root := NewUCTNode(b)
	start := time.Now()
	res := 0
	for i := 0; i < ThreadNum; i++ {
		go func() {

			for len(stop) == 0 {

				if root.Visit > iter || int(time.Since(start).Seconds()) > timeoutSeconds {
					stop <- 1
				}
				nowN := root
				deep := 0
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
					if nowN, err = Expand(nowN); err != nil {
						fmt.Println(err)
						return
					}
				}

				//nB仅仅用于模拟
				nB := board.CopyBoard(nowN.B)

				if res, err = Simulation(nB, who); err != nil {
					fmt.Println(err)
					return
				}

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
	bestChild, err := GetBestChild(root, isV)
	if err != nil {
		return nil, err
	}
	if isV {
		fmt.Printf("Tatal:%d \n MaxDeep:%d\n", root.Visit, maxDeep)
	}

	return board.MtoEdges(bestChild.LastMove), nil
}
