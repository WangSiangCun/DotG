package board

import (
	"fmt"
	"testing"
	"time"
)

func TestNewBoard(t *testing.T) {
	b := NewBoard()
	fmt.Println(b.String(), b.Boxes)
}
func TestEdgeToXYZ(t *testing.T) {
	x, y, z, _ := EdgeToXYZ(&Edge{1, 0})
	fmt.Println(x, y, z)
	if x != 1 || y != 0 || z != 0 {
		t.Fatal("错误的转换")
	}
	x, y, z, _ = EdgeToXYZ(&Edge{0, 1})
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
func TestBoard_Move(t *testing.T) {
	b := NewBoard()
	edge, _ := XYZToEdge(0, 0, 0)
	err := b.Move(edge)
	if err != nil {
		t.Fatal(err)

	}
	fmt.Println(b.String())
	edge, _ = XYZToEdge(1, 0, 0)
	err = b.Move(edge)
	if err != nil {
		t.Fatal(err)

	}
	fmt.Println(b.String())
	edge, _ = XYZToEdge(1, 1, 0)
	err = b.Move(edge)
	if err != nil {
		return
	}
	fmt.Println(b.String())
	edge, _ = XYZToEdge(0, 1, 0)
	err = b.Move(edge)
	if err != nil {
		t.Fatal(err)

	}
	fmt.Println(b.String())
}
func TestBoard_String(t *testing.T) {
	b := NewBoard()
	ms, _ := b.GetAllMoves()
	fmt.Println(len(ms))
	fmt.Println(b.String())
	for _, m := range ms {
		err := b.Move(m)
		if err != nil {
			t.Fatal(err)

		}
	}
	b.State[1][1] = 1
	fmt.Println(b.String())
}
func TestBoard_GetAllMoves(t *testing.T) {
	b := NewBoard()
	ms, _ := b.GetAllMoves()
	fmt.Println(len(ms))
	fmt.Println(b.String())
	for _, m := range ms {
		err := b.Move(m)
		if err != nil {
			t.Fatal(err)

		}
	}
	fmt.Println(b.String())
}
func TestBoard_EatAllCBox(t *testing.T) {
	now := time.Now()
	for i := 0; i < 100; i++ {
		b := NewBoard()
		for b.Status() == 0 {
			_, err := b.RandomMove()
			if err != nil {
				t.Fatal(err)
			}
			err = b.EatAllCBox()
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println(b)
		}
	}
	fmt.Println(time.Since(now))
}
func BenchmarkBoard_EatAllCBox(b *testing.B) {
	now := time.Now()
	for i := 0; i < b.N; i++ {
		b := NewBoard()
		for b.Status() == 0 {
			_, _ = b.RandomMove()
			err := b.EatAllCBox()
			if err != nil {
				return
			}
			fmt.Println(b)
			//fmt.Println(b)
		}
	}
	fmt.Println(time.Since(now))
}
func TestBoard_GetFByBI(t *testing.T) {
	b := NewBoard()
	edge, _ := XYZToEdge(0, 0, 0)
	err := b.Move(edge)
	if err != nil {
		t.Fatal(err)
	}
	edge, _ = XYZToEdge(1, 0, 0)
	err = b.Move(edge)
	if err != nil {
		t.Fatal(err)
	}
	edge, _ = XYZToEdge(1, 0, 1)
	err = b.Move(edge)
	if err != nil {
		t.Fatal(err)
	}
	f, err := b.GetFByBI(1, 1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(b.String(), f)
	if f != 1 {
		t.Fatal("自由度1监测错误")
	}

	b1 := NewBoard()
	edge, _ = XYZToEdge(0, 0, 0)
	err = b1.Move(edge)
	if err != nil {
		t.Fatal(err)
	}
	edge, _ = XYZToEdge(1, 0, 1)
	err = b1.Move(edge)
	if err != nil {
		t.Fatal(err)
	}
	f, _ = b1.GetFByBI(1, 1)
	fmt.Println(b1.String(), f)
	if f != 2 {
		t.Fatal("自由度2监测错误")
	}

	b2 := NewBoard()
	edge, _ = XYZToEdge(0, 0, 0)
	err = b2.Move(edge)
	if err != nil {
		t.Fatal(err)
	}
	f, _ = b2.GetFByBI(1, 1)
	fmt.Println(b2.String(), f)
	if f != 3 {
		t.Fatal("自由度3监测错误")
	}

	b3 := NewBoard()
	f, _ = b3.GetFByBI(1, 1)
	fmt.Println(b3.String(), f)
	if f != 4 {
		t.Fatal("自由度4监测错误")
	}

	fmt.Println(b3.String(), f)
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
func TestBoard_RandomMove(t *testing.T) {
	now := time.Now()
	for i := 0; i < 100; i++ {
		b := NewBoard()
		for b.Status() == 0 {
			_, err := b.RandomMove()
			if err != nil {
				t.Fatal(err)
			}
			err = b.EatAllCBox()
			if err != nil {
				t.Fatal(err)
				return
			}
			//fmt.Println(b)
		}
	}
	fmt.Println(time.Since(now))
}
func TestEdge_String(t *testing.T) {
	fmt.Println(&Edge{1, 2})
}
func TestXYZToEdge(t *testing.T) {
	edge, _ := XYZToEdge(1, 2, 1)
	fmt.Println(edge)
}
func TestBoard_GetBoxType(t *testing.T) {

	b := NewBoard()
	edge, _ := XYZToEdge(1, 0, 0)
	err := b.Move(edge)
	if err != nil {
		t.Fatal(err)
	}
	edge, _ = XYZToEdge(1, 0, 1)
	err = b.Move(edge)
	if err != nil {
		t.Fatal(err)
	}
	num, _ := b.GetBoxType(1, 1)
	b.Boxes[0].Type = num
	fmt.Println(num)
	fmt.Println(b.String())
	fmt.Println(b.Boxes)
	if num < 1 && num > 6 {
		t.Fatal("错误的类型")
	}

	b1 := NewBoard()
	edge, _ = XYZToEdge(0, 0, 0)
	err = b1.Move(edge)
	if err != nil {
		t.Fatal(err)
	}
	edge, _ = XYZToEdge(0, 1, 0)
	err = b1.Move(edge)
	if err != nil {
		t.Fatal(err)
	}
	num, _ = b1.GetBoxType(1, 1)
	b1.Boxes[0].Type = num
	fmt.Println(num)
	fmt.Println(b1.String())
	fmt.Println(b1.Boxes)
	if num < 1 && num > 6 {
		t.Fatal("错误的类型")
	}

	b2 := NewBoard()
	edge, _ = XYZToEdge(1, 0, 0)
	err = b2.Move(edge)
	if err != nil {
		t.Fatal(err)
	}
	edge, _ = XYZToEdge(0, 1, 0)
	err = b2.Move(edge)
	if err != nil {
		t.Fatal(err)
	}
	num, _ = b2.GetBoxType(1, 1)
	b2.Boxes[0].Type = num
	fmt.Println(num)
	fmt.Println(b2.String())
	fmt.Println(b2.Boxes)
	if num < 2 && num > 7 {
		t.Fatal("错误的类型")
	}

	b3 := NewBoard()
	edge, _ = XYZToEdge(0, 0, 0)
	err = b3.Move(edge)
	if err != nil {
		t.Fatal(err)
	}
	edge, _ = XYZToEdge(1, 0, 0)
	err = b3.Move(edge)
	if err != nil {
		t.Fatal(err)
	}
	num, _ = b3.GetBoxType(1, 1)
	b3.Boxes[0].Type = num
	fmt.Println(num)
	fmt.Println(b3.String())
	fmt.Println(b3.Boxes)
	if num < 2 && num > 7 {
		t.Fatal("错误的类型")
	}

	b4 := NewBoard()
	edge, _ = XYZToEdge(0, 0, 0)
	err = b4.Move(edge)
	if err != nil {
		t.Fatal(err)
	}
	edge, _ = XYZToEdge(1, 0, 1)
	err = b4.Move(edge)
	if err != nil {
		t.Fatal(err)
	}
	num, _ = b4.GetBoxType(1, 1)
	b4.Boxes[0].Type = num
	fmt.Println(num)
	fmt.Println(b4.String())
	fmt.Println(b4.Boxes)
	if num < 2 && num > 7 {
		t.Fatal("错误的类型")
	}

	b5 := NewBoard()
	edge, _ = XYZToEdge(0, 1, 0)
	err = b5.Move(edge)
	if err != nil {
		t.Fatal(err)
	}
	edge, _ = XYZToEdge(1, 0, 1)
	err = b5.Move(edge)
	if err != nil {
		t.Fatal(err)
	}
	num, _ = b5.GetBoxType(1, 1)
	b5.Boxes[0].Type = num
	fmt.Println(num)
	fmt.Println(b5.String())
	fmt.Println(b5.Boxes)
	if num < 2 && num > 7 {
		t.Fatal("错误的类型")
	}

	b6 := NewBoard()
	edge, _ = XYZToEdge(0, 0, 0)
	err = b6.Move(edge)
	if err != nil {
		t.Fatal(err)
	}
	num, _ = b6.GetBoxType(1, 1)
	b6.Boxes[0].Type = num
	fmt.Println(num)
	fmt.Println(b6.String())
	fmt.Println(b6.Boxes)
	if num != 1 {
		t.Fatal("错误的类型")
	}

	b7 := NewBoard()
	edge, _ = XYZToEdge(0, 1, 0)
	err = b7.Move(edge)
	if err != nil {
		t.Fatal(err)
	}
	num, _ = b7.GetBoxType(1, 1)
	b7.Boxes[0].Type = num
	fmt.Println(num)
	fmt.Println(b7.String())
	fmt.Println(b7.Boxes)
	if num != 1 {
		t.Fatal("错误的类型")
	}

	b8 := NewBoard()
	edge, _ = XYZToEdge(1, 0, 0)
	err = b8.Move(edge)
	if err != nil {
		t.Fatal(err)
	}
	num, _ = b8.GetBoxType(1, 1)
	b8.Boxes[0].Type = num
	fmt.Println(num)
	fmt.Println(b8.String())
	fmt.Println(b8.Boxes)
	if num != 1 {
		t.Fatal("错误的类型")
	}
	b9 := NewBoard()
	edge, _ = XYZToEdge(1, 0, 1)
	err = b9.Move(edge)
	if err != nil {
		t.Fatal(err)
	}
	num, _ = b9.GetBoxType(1, 1)
	b9.Boxes[0].Type = num
	fmt.Println(num)
	fmt.Println(b9.String())
	fmt.Println(b9.Boxes)
	if num != 1 {
		t.Fatal("错误的类型")
	}

	b10 := NewBoard()
	edge, _ = XYZToEdge(0, 0, 0)
	err = b10.Move(edge)
	if err != nil {
		t.Fatal(err)
	}
	edge, _ = XYZToEdge(1, 0, 1)
	err = b10.Move(edge)
	if err != nil {
		t.Fatal(err)
	}
	edge, _ = XYZToEdge(1, 0, 0)
	err = b10.Move(edge)
	if err != nil {
		t.Fatal(err)
	}
	num, _ = b10.GetBoxType(1, 1)
	b10.Boxes[0].Type = num
	fmt.Println(num)
	fmt.Println(b10.String())
	fmt.Println(b10.Boxes)
	if num != 8 {
		t.Fatal("错误的类型")
	}

}
func TestBoard_CheckoutEdge(t *testing.T) {

	b := NewBoard()
	for b.Status() == 0 {

		edge, err := b.RandomMove()
		if err != nil {
			t.Fatal(err)
		}
		err = b.CheckoutEdge(edge)
		if err != nil {
			t.Fatal(err)
			return
		}
		//fmt.Println(b.String(), b.Boxes)

	}
	fmt.Println(b.String(), b.Boxes)
	b1 := NewBoard()
	edges := []*Edge{}
	for i := 0; i < 20; i++ {
		edge, _ := b1.RandomMove()
		fmt.Println(edge)
		edges = append(edges, edge)
	}
	fmt.Println(edges)
	err := b1.CheckoutEdge(edges...)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(b1, b1.Boxes)
}
func BenchmarkBoard_CheckoutEdge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bb := NewBoard()
		for bb.Status() == 0 {
			edge, _ := bb.RandomMove()
			err := bb.CheckoutEdge(edge)
			if err != nil {
				return
			}
		}

	}
}
func TestBoard_GetFByE(t *testing.T) {

	b := NewBoard()
	for b.Status() == 0 {

		edge, err := b.RandomMove()
		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(b.GetFByE(edge))
		err = b.CheckoutEdge(edge)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(b)
		//fmt.Println(b)
	}

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

func TestBoard_GetChains(t *testing.T) {
	b := NewBoard()
	for b.Status() == 0 {
		m, err := b.RandomMove()
		if err != nil {
			t.Fatal(err)
		}
		err = b.CheckoutEdge(m)
		if err != nil {
			t.Fatal(err)
		}
		chains, _ := b.GetChains()
		for _, chain := range chains {
			fmt.Println(chain)
		}
		fmt.Println(b.String())
		fmt.Println("-------------------------------")
	}
}
func TestChain_CheckChainType(t *testing.T) {
	//长链
	b := NewBoard()
	err := b.Move(&Edge{0, 1}, &Edge{1, 0}, &Edge{0, 3}, &Edge{1, 4}, &Edge{3, 0}, &Edge{4, 1})
	if err != nil {
		t.Fatal(err)
	}
	err = b.CheckoutEdge(&Edge{0, 1}, &Edge{1, 0}, &Edge{0, 3}, &Edge{1, 4}, &Edge{3, 0}, &Edge{4, 1})
	if err != nil {
		t.Fatal(err)
	}
	mmap := map[int]bool{}
	c := NewChain()
	err = b.GetChain(1, 1, mmap, c, true)
	if err != nil {
		t.Fatal(err)
	}
	err = c.CheckChainType()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(c.String(), b.String())
	fmt.Println("---------------------")
	//环
	b1 := NewBoard()
	err = b1.Move(&Edge{0, 1}, &Edge{1, 0}, &Edge{0, 3}, &Edge{1, 4}, &Edge{3, 0}, &Edge{4, 1}, &Edge{4, 3}, &Edge{3, 4})
	if err != nil {
		t.Fatal(err)
	}
	err = b1.CheckoutEdge(&Edge{0, 1}, &Edge{1, 0}, &Edge{0, 3}, &Edge{1, 4}, &Edge{3, 0}, &Edge{4, 1}, &Edge{4, 3}, &Edge{3, 4})
	if err != nil {
		t.Fatal(err)
	}
	mmap1 := map[int]bool{}
	c1 := NewChain()
	err = b1.GetChain(1, 1, mmap1, c1, true)
	if err != nil {
		t.Fatal(err)
	}
	err = c1.CheckChainType()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(c1.String(), b1.String())
	fmt.Println("---------------------")

	b11 := NewBoard()
	err = b11.Move(&Edge{0, 1}, &Edge{1, 0}, &Edge{0, 3}, &Edge{0, 5}, &Edge{1, 6}, &Edge{3, 6}, &Edge{3, 0}, &Edge{2, 3}, &Edge{4, 1}, &Edge{4, 3}, &Edge{4, 5})
	if err != nil {
		t.Fatal(err)
	}
	err = b11.CheckoutEdge(&Edge{0, 1}, &Edge{1, 0}, &Edge{0, 3}, &Edge{0, 5}, &Edge{1, 6}, &Edge{3, 6}, &Edge{3, 0}, &Edge{2, 3}, &Edge{4, 1}, &Edge{4, 3}, &Edge{4, 5})
	if err != nil {
		t.Fatal(err)
	}
	mmap11 := map[int]bool{}
	c11 := NewChain()
	err = b11.GetChain(1, 1, mmap11, c11, true)
	if err != nil {
		t.Fatal(err)
	}
	err = c11.CheckChainType()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(c11.String(), b11.String())
	fmt.Println("---------------------")

	//二格短链
	b2 := NewBoard()
	err = b2.Move(&Edge{1, 0}, &Edge{1, 2}, &Edge{3, 0}, &Edge{3, 2})
	if err != nil {
		t.Fatal(err)
	}
	err = b2.CheckoutEdge(&Edge{1, 0}, &Edge{1, 2}, &Edge{3, 0}, &Edge{3, 2})
	if err != nil {
		t.Fatal(err)
	}
	mmap2 := map[int]bool{}
	c2 := NewChain()
	err = b2.GetChain(1, 1, mmap2, c2, true)
	if err != nil {
		t.Fatal(err)
	}
	err = c2.CheckChainType()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(c2.String(), b2.String())
	fmt.Println("---------------------")

	//模拟
	b3 := NewBoard()
	for b3.Status() == 0 {
		m, err := b3.RandomMove()
		if err != nil {
			t.Fatal(err)
		}
		err = b3.CheckoutEdge(m)
		if err != nil {
			t.Fatal(err)
		}
		chains, _ := b3.GetChains()

		for _, chain := range chains {

			err = chain.CheckChainType()
			if err != nil {
				t.Fatal(err)
			}
			if chain.Type == 4 {
				fmt.Println(1)
			}
			fmt.Println(chain)
		}
		fmt.Println(b3.String())
		fmt.Println("-------------------------------")
	}

}
func TestBoard_Get2FEdge(t *testing.T) {
	//模拟
	b3 := NewBoard()
	for b3.Status() == 0 {
		edges, err := b3.Get2FEdge()

		m, err := b3.RandomMove()
		if err != nil {
			t.Fatal(err)
		}
		err = b3.CheckoutEdge(m)
		if err != nil {
			t.Fatal(err)
		}
		if len(edges) == 0 {
			fmt.Println(b3.String())
			fmt.Println("-------------------------------")
			break
		}
	}
}
func TestBoard_GetDGridEdges(t *testing.T) {
	//模拟
	b3 := NewBoard()
	for b3.Status() == 0 {
		m, err := b3.RandomMove()
		if err != nil {
			t.Fatal(err)
		}
		if err := b3.CheckoutEdge(m); err != nil {
			t.Fatal(err)
		}
		edges, err := b3.GetDGridEdges()
		if err != nil {
			t.Fatal(err)
		}
		err = b3.Move(edges...)
		if err != nil {
			t.Fatal(err)
		}
		err = b3.CheckoutEdge(edges...)
		if err != nil {
			t.Fatal(err)
		}
		cs, err := b3.GetChains()
		if err != nil {
			t.Fatal(err)
		}
		for _, c := range cs {
			err = c.CheckChainType()
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println(c)
		}
		fmt.Println(b3, "--------------------------------")

	}

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
	err = b.Move(es...)
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
	fmt.Println(c)
	es, _ := b.GetDChainEdges(1, 1, c, c.Length-1, true)
	err = b.MoveAndCheckout(es...)
	if err != nil {
		return
	}
	fmt.Println(b)
	fmt.Println("---------------------")
}
func TestBoard_IsDCircle(t *testing.T) {
	//环
	b1 := NewBoard()
	err := b1.Move(&Edge{1, 2}, &Edge{0, 1}, &Edge{1, 0}, &Edge{0, 3}, &Edge{1, 4}, &Edge{3, 0}, &Edge{4, 1}, &Edge{4, 3}, &Edge{3, 4})
	if err != nil {
		t.Fatal(err)
	}
	err = b1.CheckoutEdge(&Edge{1, 2}, &Edge{0, 1}, &Edge{1, 0}, &Edge{0, 3}, &Edge{1, 4}, &Edge{3, 0}, &Edge{4, 1}, &Edge{4, 3}, &Edge{3, 4})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("---------------------")
	l, _ := b1.IsDCircle(1, 1)
	fmt.Println(b1, l)

}
func TestBoard_IsDCircle2(t *testing.T) {
	//环
	b1 := NewBoard()
	err := b1.MoveAndCheckout(&Edge{3, 0}, &Edge{2, 1}, &Edge{2, 3}, &Edge{3, 4}, &Edge{4, 1}, &Edge{5, 0}, &Edge{7, 0}, &Edge{8, 1}, &Edge{8, 3}, &Edge{5, 4}, &Edge{7, 4})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(b1)
	fmt.Println("---------------------")
	l, _ := b1.IsDCircle(3, 1)
	fmt.Println(b1, l)

}
func TestBoard_IsDCircle3(t *testing.T) {
	//环
	b1 := NewBoard()
	err := b1.MoveAndCheckout(&Edge{5, 2}, &Edge{3, 0}, &Edge{2, 1}, &Edge{2, 3}, &Edge{3, 4}, &Edge{4, 1}, &Edge{5, 0}, &Edge{7, 0}, &Edge{8, 1}, &Edge{8, 3}, &Edge{5, 4}, &Edge{7, 4})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(b1)
	fmt.Println("---------------------")
	l, _ := b1.IsDCircle(3, 1)
	fmt.Println(b1, l)

}
func TestBoard_GetDCircleEdges(t *testing.T) {

	//环 全吃
	b1 := NewBoard()
	err := b1.MoveAndCheckout(&Edge{5, 2}, &Edge{3, 0}, &Edge{2, 1}, &Edge{2, 3}, &Edge{3, 4}, &Edge{4, 1}, &Edge{5, 0}, &Edge{7, 0}, &Edge{8, 1}, &Edge{8, 3}, &Edge{5, 4}, &Edge{7, 4})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(b1)
	fmt.Println("---------------------")
	l, _ := b1.IsDCircle(3, 1)

	es, _ := b1.GetDCircleEdges(3, 1, l-1, false)
	_ = b1.MoveAndCheckout(es...)
	fmt.Println(b1, l)

}
func TestBoard_GetDCircleEdges1(t *testing.T) {

	//环 造双交
	b1 := NewBoard()
	err := b1.MoveAndCheckout(&Edge{5, 2}, &Edge{3, 0}, &Edge{2, 1}, &Edge{2, 3}, &Edge{3, 4}, &Edge{4, 1}, &Edge{5, 0}, &Edge{7, 0}, &Edge{8, 1}, &Edge{8, 3}, &Edge{5, 4}, &Edge{7, 4})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(b1)
	fmt.Println("---------------------")
	l, _ := b1.IsDCircle(3, 1)

	es, _ := b1.GetDCircleEdges(3, 1, l-4, true)
	_ = b1.MoveAndCheckout(es...)
	fmt.Println(b1, l)

}
func TestBoard_GetDCircleEdges2(t *testing.T) {
	//环
	b1 := NewBoard()
	err := b1.Move(&Edge{1, 2}, &Edge{0, 1}, &Edge{1, 0}, &Edge{0, 3}, &Edge{1, 4}, &Edge{3, 0}, &Edge{4, 1}, &Edge{4, 3}, &Edge{3, 4})
	if err != nil {
		t.Fatal(err)
	}
	err = b1.CheckoutEdge(&Edge{1, 2}, &Edge{0, 1}, &Edge{1, 0}, &Edge{0, 3}, &Edge{1, 4}, &Edge{3, 0}, &Edge{4, 1}, &Edge{4, 3}, &Edge{3, 4})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("---------------------")
	l, _ := b1.IsDCircle(1, 1)
	es, _ := b1.GetDCircleEdges(1, 1, l-1, false)
	_ = b1.MoveAndCheckout(es...)
	fmt.Println(b1, l)

}
func TestBoard_GetDCircleEdges3(t *testing.T) {
	//环
	b1 := NewBoard()
	err := b1.Move(&Edge{1, 2}, &Edge{0, 1}, &Edge{1, 0}, &Edge{0, 3}, &Edge{1, 4}, &Edge{3, 0}, &Edge{4, 1}, &Edge{4, 3}, &Edge{3, 4})
	if err != nil {
		t.Fatal(err)
	}
	err = b1.CheckoutEdge(&Edge{1, 2}, &Edge{0, 1}, &Edge{1, 0}, &Edge{0, 3}, &Edge{1, 4}, &Edge{3, 0}, &Edge{4, 1}, &Edge{4, 3}, &Edge{3, 4})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("---------------------")
	l, _ := b1.IsDCircle(1, 1)
	es, _ := b1.GetDCircleEdges(1, 1, l-4, true)
	_ = b1.MoveAndCheckout(es...)
	fmt.Println(b1, l)

}
func TestBoard_GetDTreeEdges(t *testing.T) {
	//环
	b1 := NewBoard()
	err := b1.MoveAndCheckout(&Edge{1, 2}, &Edge{0, 1}, &Edge{1, 0}, &Edge{0, 3}, &Edge{1, 4}, &Edge{3, 0}, &Edge{4, 1}, &Edge{4, 3}, &Edge{3, 4})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(b1)
	ees, err := b1.GetDTreeEdges()
	if err != nil {
		t.Fatal(err)
	}
	for _, es := range ees {
		b1.MoveAndCheckout(es...)
		fmt.Println()
	}
	fmt.Println(b1)
}
func TestBoard_RandomMoveByCheck(t *testing.T) {
	b := NewBoard()
	_, _ = b.RandomMoveByCheck()
	fmt.Println(b)
}
func TestBoard_GetEdgesByIdentifyingChains(t *testing.T) {

	for i := 0; i <= 100000; i++ {
		b := NewBoard()
		for b.Status() == 0 {
			//fmt.Println("front:", b)
			//es, err := b.RandomMoveByCheck()
			//fmt.Println("end:", b)
			_, err := b.RandomMoveByCheck()

			if err != nil {
				t.Fatal(err)
				return
			}
			//fmt.Println(es, "\n----------------------")
		}
		//fmt.Println(b.String(), b.Boxes)
	}

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
