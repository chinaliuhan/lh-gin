package services

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"lh-gin/constants"
	"lh-gin/models"
	"lh-gin/repositories"
	"lh-gin/requests"
	"lh-gin/tools"
	"log"
	"sync"
)

type chatService struct {
}

func NewChatService() *chatService {
	return &chatService{}
}

func (r *chatService) Login(params *requests.ChatLoginRequest) (models.User, int) {

	//get db
	info, err := repositories.NewUserManagerRepository().GetInfoByMobile(params.Mobile)
	if err != nil {
		return models.User{}, constants.SERVICE_FAILED
	}
	if info.Id <= 0 {
		return models.User{}, constants.SERVICE_NO_EXIST
	}

	//validate password
	if tools.NewGenerate().GenerateMd5(params.Password) == info.Password {
		return models.User{}, constants.SERVICE_PASSWORD_ERROR
	}

	return info, constants.SERVICE_SUCCESS
}

//接收协程
func (r *chatService) ReceiveProcess(node *constants.NodeConstant) {
	for {
		//读取数据包
		_, message, err := node.Conn.ReadMessage()
		if err != nil {
			log.Println(err.Error())
			return
		}

		fmt.Printf("[ws]recv<=%s\n", message)

		//进一步处理接收到的消息,单一服务器情况下可以直接发送
		r.Dispatch(message)

		//分布式情况下, 需要将消息广播到局域网,或者使用消息队列,或者使用NGINX中间件做分发
		//boardMsg(data)

	}
}

//发送协程
func (r *chatService) SendProcess(node *constants.NodeConstant) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Println(err.Error())
				return
			}
		}
	}
}

//映射关系表
var ClientMap map[int64]*constants.NodeConstant = make(map[int64]*constants.NodeConstant, 0)

//读写锁
var Rwlocker sync.RWMutex

//发送消息
func (r *chatService) SendMsg(userId int64, msg []byte) {
	//获取到信息就发送
	Rwlocker.Lock()
	node, ok := ClientMap[userId]
	Rwlocker.Unlock()
	if ok {
		node.DataQueue <- msg

	}
}

const (
	CMD_HEART      = 0 //心跳
	CMD_SINGLE_MSG = 1 //单聊
	CMD_ROOM_MSG   = 2 //群聊
)

//接收消息后的调度逻辑处理
func (r *chatService) Dispatch(data []byte) {
	//解析data 为message
	message := constants.MessageConstant{}
	err := json.Unmarshal(data, &message)
	if err != nil {
		//解析失败
		log.Println(err.Error())
		return
	}
	//根据cmd字段对逻辑进行处理
	switch message.Cmd {
	case CMD_SINGLE_MSG: //单聊
		//发送消息
		r.SendMsg(message.TargetID, data)
	case CMD_ROOM_MSG: //群聊
		//群聊转发逻辑
		for _, v := range ClientMap {
			if v.GroupSets.Has(message.TargetID) {
				v.DataQueue <- data
			}
		}
	case CMD_HEART: //心跳
		//一般不用管
		tools.NewLogUtil().Info("websocket 心跳: ", message.Id)
	}
}
