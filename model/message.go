package model

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
)

func init() {
	go udpSendProc()
	go udpReceiveProc()
}

type Message struct {
	gorm.Model
	FromId   int64  // 发送者
	TargetId int64  // 接收者
	Type     int    // 发送类型 群聊 私聊 广播
	Media    int    // 消息类型 文字 图片 音频
	Content  string // 消息内容
	Pic      string
	Url      string
	Desc     string
	Amount   int // 其他数字统计
}

func (table *Message) TableName() string {
	return "message"
}

type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}

// 映射关系
var clientMap = make(map[int64]*Node, 0)

// 读写锁
var rwLocker sync.RWMutex

func Chat(w http.ResponseWriter, r *http.Request) {
	// 1.获取参数并校验token
	//token := query.Get("token")
	query := r.URL.Query()
	Id := query.Get("userId")
	userId, _ := strconv.ParseInt(Id, 10, 64)
	//msgType := query.Get("type")
	//targetId := query.Get("targetId")
	//context := query.Get("context")

	isValid := true // checkToken()
	conn, err := (&websocket.Upgrader{
		// token校验
		CheckOrigin: func(request *http.Request) bool {
			return isValid
		},
	}).Upgrade(w, r, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	// 2.获取conn
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}

	// 3.用户关系
	// 4.userId跟node绑定起来并加锁
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()

	// 5.完成发送逻辑
	go sendProc(node)

	// 6.完成接收逻辑
	go receiveProc(node)

	sendMsg(userId, []byte("欢迎进入聊天室"))
}

func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func receiveProc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		broadMsg(data)
		fmt.Println("")
	}
}

var udpSendChan = make(chan []byte, 1024)

func broadMsg(data []byte) {
	udpSendChan <- data
}

// 完成udp数据发送协程
func udpSendProc() {
	con, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 6379,
	})
	defer con.Close()
	if err != nil {
		fmt.Println(err)
	}

	for {
		select {
		case data := <-udpSendChan:
			_, err := con.Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

// 完成udp数据发送协程
func udpReceiveProc() {
	con, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 6379,
	})
	if err != nil {
		fmt.Println(err)
	}
	defer con.Close()

	for {
		var buf [512]byte
		n, err := con.Read(buf[0:])
		if err != nil {
			fmt.Println(err)
			return
		}
		dispatch(buf[0:n])
	}
}

// 后端调度逻辑
func dispatch(bytes []byte) {
	msg := Message{}
	err := json.Unmarshal(bytes, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch msg.Type {
	case 1:
		sendMsg(msg.FromId, bytes)
	case 2:
		sendGroupMsg()
	case 3:
		sendAllMsg()
	}
}

func sendAllMsg() {

}

func sendGroupMsg() {

}

func sendMsg(userId int64, msg []byte) {
	rwLocker.RLocker()
	node, ok := clientMap[userId]
	rwLocker.Unlock()
	if ok {
		node.DataQueue <- msg
	}
}
