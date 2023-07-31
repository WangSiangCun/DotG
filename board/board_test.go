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
func TestEdge_String(t *testing.T) {
	fmt.Println(&Edge{1, 2})
}
func TestXYZToEdge(t *testing.T) {
	edge := XYZToEdge(1, 2, 1)
	fmt.Println(edge)
}
func TestCopyBoard(t *testing.T) {
	b := NewBoard()
	b.MoveAndCheckout(&Edge{5, 6}, &Edge{4, 7}, &Edge{5, 8}, &Edge{6, 7})
	fmt.Println(b, b.Boxes)
	b1 := CopyBoard(b)
	fmt.Println(b1, b1.Boxes, b, b.Boxes)
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
func TestBoard_GetDGridEdges(t *testing.T) {
	b := NewBoard()
	b.MoveAndCheckout(&Edge{2, 7}, &Edge{3, 6}, &Edge{3, 8}, &Edge{3, 10}, &Edge{4, 9})
	fmt.Println(b)
	fmt.Println(b.GetDGridEdges())
	_, f1 := b.Get2FEdgeAndMessage()
	fmt.Println(b.GetDGridEdgesByMessage(f1))
}
