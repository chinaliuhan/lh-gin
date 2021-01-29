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
	"time"
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

		str := fmt.Sprintf("%s", message)
		println("[ws]receive<=========================== ", str)
		if str == "ping" {
			_ = node.Conn.WriteMessage(websocket.TextMessage, []byte("pong"))
			return
		}

		//进一步处理接收到的消息,单一服务器情况下可以直接发送
		r.Dispatch(node, message)

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

//连接健康检查
func (r *chatService) CheckConnectHealth() {
	for {
		time.Sleep(time.Duration(10) * time.Second)
		for k, v := range ClientMap {
			err := v.Conn.WriteMessage(websocket.TextMessage, []byte("pong..."))
			if err != nil {
				_ = v.Conn.Close()
				Rwlocker.Lock()
				delete(ClientMap, k)
				Rwlocker.Unlock()
			}
		}
	}
}

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
func (r *chatService) Dispatch(node *constants.NodeConstant, data []byte) {
	//解析data 为message
	message := constants.MessageConstant{}
	err := json.Unmarshal(data, &message)
	if err != nil {
		//解析失败
		log.Println(err.Error())
		return
	}

	//保存聊天记录
	go r.SaveChatRecord(node, message)

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

func (r *chatService) SaveChatRecord(node *constants.NodeConstant, message constants.MessageConstant) {
	messageModel := models.ChatRecord{
		UserId:    message.Userid,
		TargetId:  message.TargetID,
		Category:  message.Cmd,
		Mime:      message.Media,
		Url:       message.Url,
		Remark:    "",
		Content:   message.Content,
		Created:   int(time.Now().Unix()),
		Updated:   0,
		CreatedIp: node.Conn.RemoteAddr().String(),
		UpdatedIp: "",
		Deleted:   0,
	}
	db := tools.NewMysqlInstance()
	_, err := db.InsertOne(messageModel)
	if err != nil {
		tools.NewLogUtil().Error("聊天记录入库失败:")
	}
}

func (r *chatService) GetChatRecordByUidAndTimeAndMyFriend(userID int64, startTime int, friends []string) ([]models.ChatRecord, error) {
	ChatRecordList := make([]models.ChatRecord, 0)
	db := tools.NewMysqlInstance()
	fid := tools.NewCommonUtil().Array2String(friends, ",")
	err := db.Where("user_id=?", userID).
		And("created>?", startTime).
		And("target_id in(?)", fid).
		And("deleted=?", 0).
		Find(&ChatRecordList)
	if err != nil {
		return nil, err
	}
	return ChatRecordList, nil
}

func (r *chatService) GetChatRecordByTargetAndTimeAndMyFriend(targetID int64, startTime int, friends []string) ([]models.ChatRecord, error) {
	ChatRecordList := make([]models.ChatRecord, 0)
	fid := tools.NewCommonUtil().Array2String(friends, ",")
	db := tools.NewMysqlInstance()

	err := db.Where("target_id=?", targetID).
		And("user_id in(?)", fid).
		And("created>?", startTime).
		And("deleted=?", 0).
		Find(&ChatRecordList)
	if err != nil {
		return nil, err
	}
	return ChatRecordList, nil
}

func (r *chatService) GetMyFriends(UserID int64) ([]models.ChatContact, error) {
	db := tools.NewMysqlInstance()
	var chatModel []models.ChatContact
	err := db.Where("user_id=?", UserID).And("type=?", constants.CONNECT_TYPE_USER).And("deleted=?", 0).Find(&chatModel)
	if err != nil {
		return nil, err
	}
	return chatModel, nil
}
