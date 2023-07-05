package uct

import (
	"dotg/board"
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sync"
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
	TriedMap map[string]bool
	LastMove []*board.Edge
}
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

	MaxChild  int  = 25
	hashBlock uint = 23 // prime
	hashSize  uint = 13000000 / hashBlock
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
	if n.Visit == 0 {
		return rand.Float64() + 1.0
	}
	return float64(n.Win)/float64(n.Visit) + ucb_C*math.Sqrt(math.Log(float64(n.Parents.Visit))/float64(n.Visit))
}
func (n *UCTNode) AddChild(child *UCTNode) {
	n.Children = append(n.Children, child)
}
func (n *UCTNode) GetUnTriedEdges() (edges [][]*board.Edge, err error) {
	if n.B.Status() != 0 {
		return nil, fmt.Errorf("游戏已经结束，没有可拓展边")
	}

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
	if len(n.Children) >= MaxChild {
		return nil, err
	}
	return
}
func NewUCTNode(b *board.Board) *UCTNode {
	return &UCTNode{
		B:        b,
		Children: []*UCTNode{},
		Parents:  nil,
		Visit:    0,
		Win:      0,
		TriedMap: map[string]bool{},
		LastMove: []*board.Edge{},
	}
}
func Search(b *board.Board, timeoutSeconds int, iter, who int, isV bool) (es []*board.Edge, err error) {
	//defer func() {
	//	count, cleanCount := CleanHashTable(b.Turn)
	//	fmt.Printf("清理hash表:%d个结点 共计:%d个结点\n", cleanCount, count)
	//}()
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
				//1:5 2:3 3:1 4:2 5:1 6:1
				if root.Visit > iter || int(time.Since(start).Seconds()) > timeoutSeconds {
					stop <- 1
				}
				nowN := root
				ees := [][]*board.Edge{}
				mutex.Lock()
				if ees, nowN, err = SelectBest(nowN); err != nil {
					fmt.Println(err)
					return
				}
				if nowN.B.Status() == 0 {
					//扩展
					if nowN, err = Expand(&ees, nowN); err != nil {
						fmt.Println(err)
						return
					}
				}
				mutex.Unlock()

				//nB仅仅用于模拟
				nB := board.CopyBoard(nowN.B)
				if res, err = Simulation(nB, who); err != nil {
					fmt.Println(err)
					return
				}

				mutex.Lock()
				/*
					//hashKey := NewHashKey(nowN.B)
					//if hashV, ok := GetHashValue(hashKey); !ok {
					//	SetHashValue(hashKey, &HashValue{
					//		Visit: 1,
					//		Win:   res,
					//		Turn:  nowN.B.Turn,
					//	})
					//} else {
					//	hashV.Win += res
					//	hashV.Visit++
					//}
				*/
				BackUp(nowN, res, who)
				mutex.Unlock()

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

	return bestChild.LastMove, nil
}
func GetBestChild(n *UCTNode, isV bool) (*UCTNode, error) {
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
func SelectBest(n *UCTNode) ([][]*board.Edge, *UCTNode, error) {
	for {
		if n.B.Status() != 0 {
			return nil, n, nil
		}
		//获取还没尝试的边
		if ees, err := n.GetUnTriedEdges(); err != nil {
			fmt.Println(err)
			return nil, nil, err
		} else if len(ees) == 0 {
			//获取最好的孩子结点
			if n, err = GetBestChild(n, false); err != nil {
				fmt.Println(err)
				return nil, nil, err
			}
		} else {
			return ees, n, nil
		}
	}

}
func Expand(edges *[][]*board.Edge, n *UCTNode) (*UCTNode, error) {
	if n.B.Status() != 0 {
		return n, nil
	}
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
	//获取hash表中的值
	//hashKey := &HashKey{M: nN.B.M, Now: nN.B.Now}
	//if hashV, ok := GetHashValue(hashKey); ok {
	//	BackUpHash(nN, hashV)
	//	}
	return nN, nil

}

func BackUp(n *UCTNode, res int, who int) {
	deep := 0
	for n != nil {
		deep++
		if deep > maxDeep {
			maxDeep = deep
		}
		if n.B.Now == who {
			n.Win += res
		} else {
			n.Win += 1 - res
		}

		n.Visit++
		n = n.Parents
	}
}
func Move(b *board.Board, timeout int, iter, who int, isV bool) {
	start := time.Now()
	es := []*board.Edge{}
	if edges2F, err := b.Get2FEdge(); err != nil {
		fmt.Println(err)
		return
	} else if len(edges2F) != 0 {
		es, err = Search(b, timeout, iter, who, isV)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		ess, err := b.GetMove()
		if err != nil {
			fmt.Println(err)
			return
		}
		es = ess[0]
	}
	b.MoveAndCheckout(es...)
	fmt.Println(es)
	fmt.Println(b, time.Since(start))
	fmt.Println("-------------------------")
}
func BackUpHash(n *UCTNode, value *HashValue) {
	for n != nil {
		fmt.Println(n)
		n.Visit += value.Visit
		n.Win += value.Win
		fmt.Println(n, value.Visit, value.Win)
		n = n.Parents

	}
}

//	func NewHashKey(b *board.Board) *HashKey {
//		return &HashKey{b.M, b.Turn}
//	}
func GetHashValue(k *HashKey) (value *HashValue, ok bool) {
	idx := uint8((k.M[0] ^ k.M[1]) % uint64(hashBlock))
	rwMutex[idx].RLock()
	defer rwMutex[idx].RUnlock()
	if v, ok := hashTable[idx][*k]; !ok {
		//	fmt.Println(*k, v, ok)
		return v, false
	} else {
		//	fmt.Println(*k, v, ok)
		return v, true
	}

}
func SetHashValue(k *HashKey, value *HashValue) {
	idx := uint8((k.M[0] ^ k.M[1]) % uint64(hashBlock))
	rwMutex[idx].Lock()
	defer rwMutex[idx].Unlock()

	hashTable[idx][*k] = value
}
func CleanHashTable(turn int) (count int, cleanCount int) {
	for i := uint(0); i < hashBlock; i++ {
		for k, v := range hashTable[i] {
			if v.Turn <= turn+1 {
				delete(hashTable[i], k)
				cleanCount++
			}
			count++

		}
	}

	return count, cleanCount
}

func init() {
	rand.Seed(time.Now().Unix())
	ThreadNum = runtime.NumCPU()
	for i := uint(0); i < hashBlock; i++ {
		hashTable[i] = make(map[HashKey]*HashValue, hashSize)
	}
}
