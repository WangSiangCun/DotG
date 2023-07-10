package board

import (
	"fmt"
	"testing"
)

func TestNewBoard(t *testing.T) {
	b := NewBoard()
	fmt.Println(b.String(), b.Boxes)
}
func TestEdgeToXYZ(t *testing.T) {
	x, y, z := EdgeToXYZ(&Edge{1, 0})
	fmt.Println(x, y, z)
	if x != 1 || y != 0 || z != 0 {
		t.Fatal("错误的转换")
	}
	x, y, z = EdgeToXYZ(&Edge{0, 1})
	if x != 0 || y != 0 || z != 0 {
		t.Fatal("错误的转换")
	}
	fmt.Println(x, y, z)
}
func TestBoxToXY(t *testing.T) {
	x, y, _ := BoxToXY(3, 5)
	fmt.Println(x, y)
	if x != 1 || y != 2 {
		t.Fatal("错误的转换")
	}
}
func BenchmarkBoard_BitMove1(b *testing.B) {
	bb := NewBoard()
	for i := 0; i < b.N; i++ {
		bb.State[1][1] = 1

	}
}

func TestBoard_Status(t *testing.T) {
	b := NewBoard()
	fmt.Println(b.Status())
}
func TestBoard_GetEdgeByBI(t *testing.T) {
	b := NewBoard()
	err := b.Move(&Edge{1, 0})
	if err != nil {
		t.Fatal(err)
		return
	}
	edges, _ := b.GetEdgeByBI(1, 1)
	fmt.Println(edges)
}
func TestEdge_String(t *testing.T) {
	fmt.Println(&Edge{1, 2})
}
func TestXYZToEdge(t *testing.T) {
	edge := XYZToEdge(1, 2, 1)
	fmt.Println(edge)
}
func TestBoard_GetChain(t *testing.T) {
	b := NewBoard()
	err := b.Move(&Edge{0, 1}, &Edge{1, 0}, &Edge{0, 3}, &Edge{1, 4}, &Edge{3, 0}, &Edge{4, 1})
	if err != nil {
		t.Fatal(err)
	}
	err = b.CheckoutEdge(&Edge{0, 1}, &Edge{1, 0}, &Edge{0, 3}, &Edge{1, 4}, &Edge{3, 0}, &Edge{4, 1})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(b.String())
	mmap := map[int]bool{}
	c := NewChain()
	err = b.GetChain(1, 1, mmap, c, true)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(c.String(), b.String())

	b1 := NewBoard()
	err = b1.Move(&Edge{0, 1}, &Edge{1, 0}, &Edge{0, 3}, &Edge{1, 4}, &Edge{3, 0}, &Edge{4, 1}, &Edge{4, 3}, &Edge{3, 4})
	if err != nil {
		t.Fatal(err)
	}
	err = b1.CheckoutEdge(&Edge{0, 1}, &Edge{1, 0}, &Edge{0, 3}, &Edge{1, 4}, &Edge{3, 0}, &Edge{4, 1}, &Edge{4, 3}, &Edge{3, 4})
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Println(b1.String())
	mmap1 := map[int]bool{}
	c1 := NewChain()
	err = b1.GetChain(1, 1, mmap1, c1, true)
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Println(c1.String(), b1.String())

	b2 := NewBoard()
	err = b2.Move(&Edge{1, 0}, &Edge{1, 2}, &Edge{3, 0}, &Edge{3, 2})
	if err != nil {
		t.Fatal(err)
	}
	err = b2.CheckoutEdge(&Edge{1, 0}, &Edge{1, 2}, &Edge{3, 0}, &Edge{3, 2})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(b2.String())
	mmap2 := map[int]bool{}
	c2 := NewChain()
	err = b2.GetChain(1, 1, mmap2, c2, true)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(c2.String(), b2.String())
}
func TestBoard_GetChain1(t *testing.T) {
	b := NewBoard()
	err := b.MoveAndCheckout(&Edge{2, 1}, &Edge{1, 2})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(b.String())

	mmap := map[int]bool{}
	c := NewChain()
	err = b.GetChain(1, 1, mmap, c, true)
	if err != nil {
		t.Fatal(err)
	}
	c.CheckChainType()
	fmt.Println(c.String(), b.String())

}
func TestBoard_GetChain2(t *testing.T) {
	b := NewBoard()
	err := b.MoveAndCheckout(&Edge{6, 1}, &Edge{7, 2}, &Edge{9, 2}, &Edge{10, 1})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(b.String())

	mmap := map[int]bool{}
	c := NewChain()
	err = b.GetChain(7, 1, mmap, c, true)
	if err != nil {
		t.Fatal(err)
	}
	c.CheckChainType()
	fmt.Println(c.String(), b.String())

}

func TestBoard_GetDChainEdges(t *testing.T) {
	//长链 双交
	b := NewBoard()
	err := b.Move(&Edge{0, 1}, &Edge{1, 0}, &Edge{0, 3}, &Edge{1, 4}, &Edge{3, 0}, &Edge{4, 1}, &Edge{3, 2})
	if err != nil {
		t.Fatal(err)
	}
	err = b.CheckoutEdge(&Edge{0, 1}, &Edge{1, 0}, &Edge{0, 3}, &Edge{1, 4}, &Edge{3, 0}, &Edge{4, 1}, &Edge{3, 2})
	if err != nil {
		t.Fatal(err)
	}
	c := NewChain()
	boxMark := map[int]bool{}
	err = b.GetChain(1, 1, boxMark, c, true)
	if err != nil {
		return
	}
	fmt.Println(c)
	es, _ := b.GetDChainEdges(3, 1, c, c.Length-2, true)
	err = b.MoveAndCheckout(es...)
	if err != nil {
		return
	}
	fmt.Println(b)
	fmt.Println("---------------------")
}
func TestBoard_GetDChainEdges2(t *testing.T) {
	//长链 全捕获
	b := NewBoard()
	err := b.Move(&Edge{0, 1}, &Edge{1, 0}, &Edge{0, 3}, &Edge{1, 4}, &Edge{3, 0}, &Edge{4, 1}, &Edge{3, 2})
	if err != nil {
		t.Fatal(err)
	}
	err = b.CheckoutEdge(&Edge{0, 1}, &Edge{1, 0}, &Edge{0, 3}, &Edge{1, 4}, &Edge{3, 0}, &Edge{4, 1}, &Edge{3, 2})
	if err != nil {
		t.Fatal(err)
	}
	c := NewChain()
	boxMark := map[int]bool{}
	err = b.GetChain(1, 1, boxMark, c, true)
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Println(c)
	es, _ := b.GetDChainEdges(3, 1, c, c.Length-1, true)
	err = b.MoveAndCheckout(es...)
	if err != nil {
		return
	}
	fmt.Println(b)
	fmt.Println("---------------------")
}
func TestBoard_GetDChainEdges3(t *testing.T) {
	//死短链 双交
	b := NewBoard()
	err := b.MoveAndCheckout(&Edge{0, 1}, &Edge{1, 0}, &Edge{2, 1}, &Edge{0, 3}, &Edge{2, 3})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(b)
	c := NewChain()
	boxMark := map[int]bool{}
	err = b.GetChain(1, 3, boxMark, c, true)
	if err != nil {
		return
	}
	fmt.Println(c)
	es, _ := b.GetDChainEdges(1, 1, c, c.Length-2, true)
	err = b.Move(es...)
	if err != nil {
		return
	}
	fmt.Println(b)
	fmt.Println("---------------------")
}
func TestBoard_GetDChainEdges4(t *testing.T) {
	//死短链 全吃
	b := NewBoard()
	err := b.MoveAndCheckout(&Edge{0, 1}, &Edge{1, 0}, &Edge{2, 1}, &Edge{0, 3}, &Edge{2, 3})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(b)
	c := NewChain()
	boxMark := map[int]bool{}
	err = b.GetChain(1, 3, boxMark, c, true)
	if err != nil {
		return
	}
	fmt.Println(c)
	es, _ := b.GetDChainEdges(1, 1, c, c.Length-1, true)
	err = b.MoveAndCheckout(es...)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(b)
	fmt.Println("---------------------")
}
func TestBoard_GetDChainEdges5(t *testing.T) {
	//长链 全捕获
	b := NewBoard()
	err := b.MoveAndCheckout(&Edge{4, 1}, &Edge{0, 1}, &Edge{1, 0}, &Edge{0, 3}, &Edge{1, 4}, &Edge{3, 0}, &Edge{3, 4}, &Edge{3, 2})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(b)
	if err != nil {
		t.Fatal(err)
		return
	}
	d, a, _ := b.GetDTreeEdges()
	fmt.Println(d, a, b)
	err = b.MoveAndCheckout(d...)
	fmt.Println(b, "---------------------")
}
func TestBoard_GetDChainEdges6(t *testing.T) {
	//死环 全捕获
	b := NewBoard()
	err := b.MoveAndCheckout(&Edge{4, 3}, &Edge{4, 1}, &Edge{0, 1}, &Edge{1, 0}, &Edge{0, 3}, &Edge{1, 2}, &Edge{1, 4}, &Edge{3, 0}, &Edge{3, 4})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(b)
	d, es, err := b.GetDTreeEdges()
	fmt.Println(b, d, es)
	if err != nil {
		t.Fatal(err)
		return
	}
	err = b.MoveAndCheckout(d...)
	if err != nil {
		t.Fatal(err)
		return
	}
	if err != nil {
		return
	}
	fmt.Println(b, d, es)
	fmt.Println("---------------------")
}
func TestBoard_GetDChainEdges7(t *testing.T) {
	//长链 全捕获
	b := NewBoard()
	err := b.MoveAndCheckout(&Edge{0, 1}, &Edge{1, 0}, &Edge{0, 3}, &Edge{1, 2}, &Edge{1, 4}, &Edge{3, 0}, &Edge{3, 4}, &Edge{3, 2})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(b)
	c := NewChain()
	boxMark := map[int]bool{}
	err = b.GetChain(3, 1, boxMark, c, true)
	if err != nil {
		t.Fatal(err)
		return
	}
	d, es, err := b.GetDTreeEdges()
	fmt.Println(b, d, es)
	if err != nil {
		t.Fatal(err)
		return
	}
	err = b.MoveAndCheckout(d...)
	if err != nil {
		t.Fatal(err)
		return
	}
	if err != nil {
		return
	}
	fmt.Println(b, d, es)
	fmt.Println("---------------------")
}
func TestBoard_RandomMoveByCheck(t *testing.T) {
	b := NewBoard()
	_, _ = b.RandomMoveByCheck()
	fmt.Println(b)
}
func BenchmarkBoard_GetEdgesByIdentifyingChains(b *testing.B) {
	bb := NewBoard()
	for i := 0; i < b.N; i++ {
		_, err := bb.RandomMoveByCheck()
		if err != nil {
			return
		}
		if bb.Status() != 0 {
			//fmt.Println(bb)
			bb = NewBoard()

		}
	}
}
func TestCopyBoard(t *testing.T) {
	b := NewBoard()
	b.MoveAndCheckout(&Edge{5, 6}, &Edge{4, 7}, &Edge{5, 8}, &Edge{6, 7})
	fmt.Println(b, b.Boxes)
	b1 := CopyBoard(b)
	fmt.Println(b1, b1.Boxes, b, b.Boxes)
}
func TestBoard_GetMove(t *testing.T) {
	for i := 0; i <= 100000; i++ {
		b := NewBoard()
		for b.Status() == 0 {
			//fmt.Println("front:", b)
			es, err := b.RandomMoveByCheck()
			fmt.Println(b)
			//_, err := b.RandomMoveByCheck()

			if err != nil {
				t.Fatal(err)
				return
			}
			//fmt.Printf("%b%b", b.M[0], b.M[1])
			fmt.Println(es, "\n----------------------")
		}
		//fmt.Println(b.String(), b.Boxes)
	}

}
func TestEdgesToHV(t *testing.T) {
	h, v := EdgesToHV([]*Edge{&Edge{1, 2}, &Edge{6, 9}}...)
	fmt.Printf("%b %b\n", h, v)
}
func TestEdgesToM(t *testing.T) {
	//M, _ := EdgesToM([]*Edge{&Edge{0, 1}, &Edge{0, 3}, &Edge{1, 2}, &Edge{6, 9}, &Edge{10, 9}, &Edge{9, 10}}...)
	//fmt.Printf("%b\n", M)
	es := []*Edge{}
	for i := 0; i < 11; i++ {
		for j := 0; j < 11; j++ {
			if (i+j)&1 == 1 {
				es = append(es, &Edge{i, j})
				//M, _ := EdgesToM(&Edge{i, j})
				//fmt.Printf("%b\n", M)
			}
		}
	}
	M := EdgesToM(es...)
	fmt.Printf("%b\n", M)
	es = MtoEdges(M)
	fmt.Println(es)
}
func TestBoard_GetSafeNo4Edge(t *testing.T) {
	b := NewBoard()
	b.MoveAndCheckout(&Edge{1, 2}, &Edge{3, 2}, &Edge{1, 0}, &Edge{0, 7})
	es, _ := b.GetSafeNo4Edge()
	fmt.Println(es, b)
}
func TestBoard_GetEdgeBy12LChain(t *testing.T) {
	//单格
	b := NewBoard()
	b.MoveAndCheckout(&Edge{1, 2}, &Edge{3, 2}, &Edge{1, 0}, &Edge{0, 7})
	es, _ := b.GetEdgeBy12LChain()
	fmt.Println(es, b)
	//短链
	b1 := NewBoard()
	b1.MoveAndCheckout(&Edge{1, 2}, &Edge{3, 2}, &Edge{1, 0}, &Edge{0, 7}, &Edge{3, 0})
	es1, _ := b1.GetEdgeBy12LChain()
	fmt.Println(es1, b1)
}
func TestBoard_GetSafeAndChain12Edge(t *testing.T) {

	//单格
	b := NewBoard()
	b.MoveAndCheckout(&Edge{1, 2}, &Edge{3, 2}, &Edge{1, 0}, &Edge{0, 7}, &Edge{3, 0})
	es1, _ := b.GetSafeAndChain12Edge()
	fmt.Println(es1, b)

}
