package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"

	"github.com/code-art/gin-im/util"
)

func init() {
	go udpSendProc()
	go udpReceiveProc()
	fmt.Println("init goroutine ")
}

type Message struct {
	util.Model
	FromId   int64  `json:"fromId,omitempty"`   // 发送者
	TargetId int64  `json:"targetId,omitempty"` // 接收者
	Type     int    `json:"type,omitempty"`     // 发送类型 群聊 私聊 广播
	Media    int    `json:"media,omitempty"`    // 消息类型 文字 图片 音频
	Content  string `json:"content,omitempty"`  // 消息内容
	Pic      string `json:"pic,omitempty"`
	Url      string `json:"url,omitempty"`
	Desc     string `json:"desc,omitempty"`
	Amount   int    `json:"amount,omitempty"` // 其他数字统计
}

func (table *Message) TableName() string {
	return "message"
}

// 最大心跳时间
const HeartbeatMaxTime = time.Minute

type Node struct {
	Conn          *websocket.Conn // 连接
	Addr          string          // 客户端地址
	FirstTime     uint64          // 首次连接时间
	HeartbeatTime uint64          // 心跳时间
	LoginTime     uint64          // 登录时间
	DataQueue     chan []byte     // 消息
	GroupSets     set.Interface   // 群
}

// 映射关系
var clientMap = make(map[int64]*Node, 0)

// 读写锁
var rwLocker sync.RWMutex

func Chat(w http.ResponseWriter, r *http.Request) {
	// 1.获取参数并校验token
	// token := query.Get("token")
	query := r.URL.Query()
	Id := query.Get("userId")
	userId, _ := strconv.ParseInt(Id, 10, 64)
	// msgType := query.Get("type")
	// targetId := query.Get("targetId")
	// context := query.Get("context")

	isValid := true // checkToken()
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(request *http.Request) bool {
			return isValid
		},
	}).Upgrade(w, r, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	// 2.获取conn
	currentTime := uint64(time.Now().Unix())
	node := &Node{
		Conn:          conn,
		Addr:          conn.RemoteAddr().String(),
		HeartbeatTime: currentTime,
		LoginTime:     currentTime,
		DataQueue:     make(chan []byte, 50),
		GroupSets:     set.New(set.ThreadSafe),
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
			fmt.Println("[ws]sendProc >>>> msg:", string(data))
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
		fmt.Println("[ws] receiveProc <<<<<", string(data))
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
		Port: 3000,
	})
	defer con.Close()
	if err != nil {
		fmt.Println(err)
	}

	for {
		select {
		case data := <-udpSendChan:
			fmt.Println("udpSendProc data: ", string(data))
			_, err := con.Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

// 完成udp数据接收协程
func udpReceiveProc() {
	con, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
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
		fmt.Println("udpReceiveProc data: ", string(buf[0:n]))
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
		sendMsg(msg.TargetId, bytes)
	case 2:
		sendGroupMsg(msg.TargetId, bytes)
	case 3:
		sendAllMsg()
	case 4:
	}
}

func sendMsg(userId int64, msg []byte) {
	fmt.Printf("send message >>> userId: %d, msg: %s\n", userId, string(msg))
	rwLocker.RLock()
	node, ok := clientMap[userId]
	rwLocker.RUnlock()
	if ok {
		node.DataQueue <- msg
	}
}

func JoinGroup(userId uint, groupId uint) (int, string) {
	co := Community{}
	err := util.DB.Where("id = ?", groupId).First(&co).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return -1, "此群不存在"
	}

	c := Contact{
		OwnerId:  userId,
		TargetId: groupId,
		Type:     2,
	}
	err = util.DB.Where("owner_id = ? and target_id = ? and type = 2", userId, groupId).First(&c).Error
	if err == nil {
		return -1, "已经加入此群"
	}

	util.DB.Create(&c)
	return 1, "加群成功"
}

func sendGroupMsg(targetId int64, msg []byte) {
	fmt.Println("发送群消息")
	userIds := FindUserByGroupId(uint(targetId))
	for _, id := range userIds {
		sendMsg(int64(id), msg)
	}
}

func sendAllMsg() {

}

// 更新用户心跳
func (n *Node) Heartbeat(currentTime uint64) {
	n.HeartbeatTime = currentTime
	return
}

// 判断用户心跳是否超时
func (n *Node) IsHeartbeatTimeOut(currentTime uint64) (timeout bool) {
	heartbeatMaxTime := viper.GetUint64("task.heartbeatMaxTime")
	if n.HeartbeatTime+heartbeatMaxTime <= currentTime {
		fmt.Println("心跳超时...自动下线")
		timeout = true
	}
	return
}

func ClearConnection(params any) (ans bool) {
	ans = true
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("clean connections err: ", e)
		}
	}()
	fmt.Println("定时任务: 清理超时连接", params)
	currentTime := uint64(time.Now().Unix())
	for k := range clientMap {
		node := clientMap[k]
		if node.IsHeartbeatTimeOut(currentTime) {
			fmt.Println("心跳超时...关闭连接")
			_ = node.Conn.Close()
		}
	}
	return ans
}
