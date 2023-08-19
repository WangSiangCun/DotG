package record

import (
	"dotg/board"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	r       = ""
	b       = ""
	content = ""
)

func ClearContent() {
	r = ""
	b = ""
	content = ""
}
func SetR(rr string) {
	r = rr
}
func SetB(bb string) {
	b = bb
}
func PrintContentStart(rScore, bScore, date string) {
	t := "{"
	t += "\"R\": " + "\"" + r + "\","
	t += "\"B\": " + "\"" + b + "\","
	if rScore > bScore {
		t += "\"winner\": " + "\"" + "R" + "\","
	} else {
		t += "\"winner\": " + "\"" + "B" + "\","
	}
	t += "\"RScore\": " + rScore + ","
	t += "\"BScore\": " + bScore + ","
	t += "\"Date\": " + "\"" + date + "\","
	t += "\"Event\": " + "\"" + "" + "\","
	t += "\"game\": ["
	content = t + content
}

func PrintContentMiddle(b *board.Board, moves []*board.Edge) {
	remain := len(moves)
	res := []*board.Edge{}
	i := 0
	tb := board.CopyBoard(b)
	nB := board.CopyBoard(tb)
	for remain != 1 {

		if nB.MoveAndCheckoutForPrint(moves[i]) {
			tb.MoveAndCheckout(moves[i])
			res = append(res, moves[i])

			tM := moves[len(moves)-1]
			moves[len(moves)-1] = moves[i]
			moves[i] = tM
			moves = moves[:len(moves)-1]

			i = 0
			remain--

		} else {
			nB = board.CopyBoard(tb)
			i++
		}

	}
	res = append(res, moves[0])

	c := ""
	if b.Now == 1 {
		c = "b"
	} else {
		c = "r"
	}
	for _, m := range res {
		x, y, z := board.EdgeToXYZ(m)
		aStr, bStr := "", ""
		if x == 0 {
			bStr = strconv.Itoa(6 - y)
			aStr = string('a' + z)
			content += "{\"piece\": \"" + c + "(" + aStr + bStr + ",h)\"},"

		} else {
			bStr = strconv.Itoa(5 - y)
			aStr = string('a' + z)
			content += "{\"piece\": \"" + c + "(" + aStr + bStr + ",v)\"},"
		}

	}

}
func PrintContentBack() {
	content = content[0 : len(content)-1]
	content += "]}"
}
func WriteToFile(bb *board.Board) {

	desiredFileName := ""
	if bb.Status() == 1 {

		desiredFileName = "DB：" + r + " vs " + b + "：先手胜.txt"
	} else {
		desiredFileName = "DB：" + r + " vs " + b + "：后手胜.txt"
	}
	newFileName := desiredFileName

	for {
		_, err := os.Stat(newFileName)
		if os.IsNotExist(err) {
			break
		}

		// 文件存在，添加后缀
		newFileName = addSuffix(newFileName)
	}

	// 创建文件并写入内容
	file, err := os.Create(newFileName)
	if err != nil {
		fmt.Println("无法创建文件:", err)
		return
	}
	defer file.Close()
	_, err = file.WriteString(content)
	if err != nil {
		fmt.Println("无法写入文件内容:", err)
		return
	}

	fmt.Printf("文件 %s 创建成功并已写入内容。\n", newFileName)

}
func addSuffix(fileName string) string {
	strs := strings.Split(fileName, ".")
	s := strs[0] + "+." + strs[1]
	return s
}
