package main

import (
	"dotg/algorithm/uct"
	"dotg/board"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	Address string = ":8222"
	TimeOut int    = 30
)

// 使用 Gorilla WebSocket 库
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var (
	b      *board.Board
	AITurn int
)

func sendEdges(conn *websocket.Conn, es []*board.Edge) {
	builder := strings.Builder{}
	builder.WriteString("Moves-")
	for i := 0; i < len(es); i++ {
		if i < len(es)-1 {
			builder.WriteString(strconv.Itoa(es[i].X) + "," + strconv.Itoa(es[i].Y) + "-")
		} else {
			builder.WriteString(strconv.Itoa(es[i].X) + "," + strconv.Itoa(es[i].Y))
		}
	}
	// 向客户端发送消息
	if err := conn.WriteMessage(websocket.TextMessage, []byte(builder.String())); err != nil {
		log.Println("发送消息错误:", err)
	}
}
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("无法升级 WebSocket 连接:", err)
		return
	}
	defer conn.Close()

	for {
		// 读取客户端发送的消息
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("读取消息错误:", err)
			break
		}
		msg := string(message)
		log.Printf("接收到消息: %s\n", msg)
		if strings.HasPrefix(msg, "NewGame") {
			b = board.NewBoard()
			str := strings.Split(msg, "NewGame,")
			AITurn, err = strconv.Atoi(str[1])
			if err != nil {
				log.Println(err)
			}
			if AITurn == 1 {
				es := uct.Move(b, TimeOut, 1000000, AITurn, true, true)
				fmt.Println(es)
				sendEdges(conn, es)
			}
		} else if strings.HasPrefix(msg, "Moves") {
			strs := strings.Split(msg, "Moves-")
			str := strs[1]
			moves := strings.Split(str, "-")
			es := []*board.Edge{}
			for _, m := range moves {
				ij := strings.Split(m, ",")
				i, _ := strconv.Atoi(ij[0])
				j, _ := strconv.Atoi(ij[1])
				es = append(es, &board.Edge{i, j})
			}
			fmt.Println(es)
			b.MoveAndCheckout(es...)
			es = uct.Move(b, TimeOut, 1000000, AITurn, true, true)
			sendEdges(conn, es)

		}
	}
}
func main() {
	http.HandleFunc("/ws", handleWebSocket)
	log.Fatal(http.ListenAndServe(Address, nil))
}
