package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/fatih/set.v0"
	"lh-gin/constants"
	"lh-gin/models"
	"lh-gin/requests"
	"lh-gin/services"
	"lh-gin/tools"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
)

type chatController struct {
	context *gin.Context
}

func NewChatController() *chatController {

	return &chatController{}
}

/**
聊天注册
*/
func (r *chatController) Register(ctx *gin.Context) {
	//非POST请求直接返回模板
	if ctx.Request.Method != http.MethodPost {
		ctx.HTML(http.StatusOK, "/chat/register.shtml", nil)
		return
	}

	//prepare
	var (
		err          error
		hashPassword []byte
	)

	var tmp = models.User{}
	tmp.CreatedIp = ctx.ClientIP()
	tmp.UpdatedIp = ctx.ClientIP()
	paramsRequest := requests.ChatRegisterRequest{}
	err = ctx.ShouldBind(&paramsRequest)
	if err != nil {
		tools.NewLogUtil().Error(err.Error())
		tools.NewResponse(ctx).JsonFailed("参数错误")
		return
	}

	tmp.Mobile = paramsRequest.Mobile
	hashPassword, err = bcrypt.GenerateFromPassword([]byte(paramsRequest.Password), bcrypt.DefaultCost)
	tmp.Password = string(hashPassword)
	tmp.Memo = paramsRequest.Memo
	tmp.Avatar = paramsRequest.Avatar
	tmp.Sex = paramsRequest.Sex
	tmp.Nickname = paramsRequest.Nickname

	//services
	_, err = services.NewUserService().AddNew(tmp)
	if err != nil {
		tools.NewResponse(ctx).JsonFailed("注册失败")
		return
	}
	tools.NewResponse(ctx).JsonSuccess(nil)
}

/**
聊天登录
*/
func (r *chatController) Login(ctx *gin.Context) {
	//非POST请求直接返回模板
	if ctx.Request.Method != http.MethodPost {
		ctx.HTML(http.StatusOK, "/chat/login.shtml", nil)
		return
	}

	//prepare
	var (
		err              error
		serviceCode      int
		info             models.User
		chatLoginRequest requests.ChatLoginRequest
	)

	//payload on json
	if err = ctx.ShouldBind(&chatLoginRequest); err != nil {
		tools.NewLogUtil().SugarPrint(err.Error())
		tools.NewResponse(ctx).JsonFailed("参数错误")
		return
	}

	//services
	info, serviceCode = services.NewChatService().Login(&chatLoginRequest)
	token := fmt.Sprintf("%s...%d", tools.NewGenerate().GenerateUUID(), info.Id)
	token = tools.NewGenerate().GenerateMd5(token)
	if serviceCode == constants.SERVICE_SUCCESS {
		//session
		userLoginRecordModel := models.UserLoginRecord{}
		userLoginRecordModel.Token = token
		userLoginRecordModel.UserId = info.Id
		userLoginRecordModel.CreatedIp = ctx.ClientIP()
		userLoginRecordModel.Created = time.Now().Unix()
		_, err = tools.NewMysqlInstance().InsertOne(&userLoginRecordModel)
		if err != nil {
			tools.NewResponse(ctx).JsonFailed("登录状态维持失败")
			tools.NewLogUtil().Error(err.Error())
			return
		}

		tokenMap := make(map[string]interface{}, 1)
		tokenMap["token"] = token
		tokenMap["id"] = info.Id
		tokenMap["avatar"] = info.Avatar
		tokenMap["nickname"] = info.Nickname
		tokenMap["memo"] = info.Memo
		//responseMap := tools.NewJsonUtil().Encode(tokenMap)
		tools.NewResponse(ctx).JsonSuccess(tokenMap)
		return
	}

	tools.NewResponse(ctx).JsonFailed("登录失败")
	return
}

/**
首页
*/
func (r *chatController) Index(ctx *gin.Context) {
	//非POST请求直接返回模板
	if ctx.Request.Method != http.MethodPost {
		ctx.HTML(http.StatusOK, "/chat/index.shtml", nil)
		return
	}

	tools.NewResponse(ctx).JsonSuccess(nil)
	return
}

/**
获取好友列表
*/
func (r *chatController) GetFriendList(ctx *gin.Context) {

	var (
		err   error
		token string
	)

	paramsRequest := requests.ChatConnectSendRequest{}
	if err = ctx.ShouldBind(&paramsRequest); err != nil {
		tools.NewLogUtil().SugarPrint(err.Error())
		tools.NewResponse(ctx).JsonFailed("token参数错误")
		return
	}
	token = paramsRequest.Token
	tools.NewLogUtil().Info(token)
	userLoginRecord := models.UserLoginRecord{}
	db := tools.NewMysqlInstance()
	_, err = db.Where("token=?", token).And("deleted=0").Get(&userLoginRecord)
	if err != nil {
		tools.NewLogUtil().SugarPrint(err.Error())
		tools.NewResponse(ctx).JsonFailed("登录状态获取失败")
		return
	}

	userID := userLoginRecord.UserId

	db = tools.NewMysqlInstance()
	userContact := make([]models.ChatContact, 0)
	err = db.Where("user_id=?", userID).Find(&userContact)
	if err != nil {
		tools.NewLogUtil().SugarPrint(err.Error())
		tools.NewResponse(ctx).JsonFailed("系统繁忙")
		return
	}
	if len(userContact) == 0 {
		tools.NewResponse(ctx).JsonFailed("没有好友")
		return
	}
	targetIDList := make([]int64, 0)
	for _, v := range userContact {
		targetIDList = append(targetIDList, v.TargetId)
	}

	//todo 返回了非必要的数据,尤其是密码和GA
	userList := make([]models.User, 0)
	_ = db.In("id", targetIDList).Find(&userList)

	tools.NewResponse(ctx).JsonSuccess(&userList)
	return
}

/**
添加好友
*/
func (r *chatController) AddFriend(ctx *gin.Context) {

	//prepare
	var (
		err           error
		paramsRequest requests.ChatAddFriendRequest
	)

	chatToken := requests.ChatConnectSendRequest{}
	if err = ctx.ShouldBind(&chatToken); err != nil {
		tools.NewLogUtil().SugarPrint(err.Error())
		tools.NewResponse(ctx).JsonFailed("token参数错误")
		return
	}
	token := chatToken.Token
	tools.NewLogUtil().Info(token)
	userLoginRecord := models.UserLoginRecord{}
	db := tools.NewMysqlInstance()
	_, err = db.Where("token=?", token).And("deleted=0").Get(&userLoginRecord)
	if err != nil {
		tools.NewLogUtil().SugarPrint(err.Error())
		tools.NewResponse(ctx).JsonFailed("登录状态获取失败")
		return
	}

	userID := userLoginRecord.UserId

	//payload on json
	if err = ctx.ShouldBind(&paramsRequest); err != nil {
		tools.NewLogUtil().SugarPrint(err.Error())
		tools.NewResponse(ctx).JsonFailed("参数错误")
		return
	}

	targetID, _ := strconv.ParseInt(paramsRequest.TargetID, 10, 64)
	if userID == targetID {
		tools.NewResponse(ctx).JsonFailed("不能添加自己为好友")
		return
	}

	contactModel := models.ChatContact{}
	userModel := models.User{}
	db = tools.NewMysqlInstance()
	if ok, err := db.Where("id=?", targetID).And("deleted=?", constants.DbConstant{}.NotDeleted).Get(&userModel); !ok {
		tools.NewLogUtil().SugarPrint("查询对方失败:", err)
		tools.NewResponse(ctx).JsonFailed("查无此人")
		return
	}
	if _, err := db.Where("user_id=?", userID).And("target_id=?", targetID).Get(&contactModel); err != nil {
		tools.NewLogUtil().SugarPrint("查询好友失败:", err)
		tools.NewResponse(ctx).JsonFailed("系统繁忙")
		return
	}
	if contactModel.Id > 0 {
		tools.NewResponse(ctx).JsonFailed("你们已经是好友啦")
		return
	}

	transaction := db.NewSession()
	err = transaction.Begin()
	if err != nil {
		tools.NewLogUtil().SugarPrint(err)
		tools.NewResponse(ctx).JsonFailed("系统繁忙")
		return
	}

	_, err = transaction.Insert(models.ChatContact{
		UserId:    userID,
		TargetId:  targetID,
		Created:   0,
		Updated:   0,
		CreatedIp: ctx.ClientIP(),
		Deleted:   0,
	})
	if err != nil {
		_ = transaction.Rollback()
		tools.NewLogUtil().SugarPrint(err)
		tools.NewResponse(ctx).JsonFailed("系统繁忙")
		return
	}
	_, err = transaction.Insert(
		models.ChatContact{
			UserId:    targetID,
			TargetId:  userID,
			Created:   0,
			Updated:   0,
			CreatedIp: ctx.ClientIP(),
			Deleted:   0,
		})
	if err != nil {
		_ = transaction.Rollback()
		tools.NewLogUtil().SugarPrint(err)
		tools.NewResponse(ctx).JsonFailed("系统繁忙")
		return
	}
	err = transaction.Commit()
	if err != nil {
		_ = transaction.Rollback()
		tools.NewLogUtil().SugarPrint(err)
		tools.NewResponse(ctx).JsonFailed("系统繁忙")
		return
	}

	tools.NewResponse(ctx).JsonSuccess("")
	return
}

/**
获取群
*/
func (r *chatController) GetCommunityList(ctx *gin.Context) {

	tools.NewResponse(ctx).JsonSuccess(nil)
	return
}

/**
添加群
*/
func (r *chatController) AddCommunityList(ctx *gin.Context) {

	tools.NewResponse(ctx).JsonSuccess(nil)
	return
}

/**
上传头像
*/
func (r *chatController) Upload(ctx *gin.Context) {

	//prepare
	var (
		file     *multipart.FileHeader
		err      error
		fileName string
	)
	//take file
	file, err = ctx.FormFile("file")
	if err != nil {
		tools.NewLogUtil().Error(err.Error())
		tools.NewResponse(ctx).JsonFailed("请选择文件")
		return
	}

	//fileName = file.Filename + tools.NewGenerate().GenerateUUID()
	fileName = file.Filename
	//todo 暂时存储,后面会用云存储,可通过配置使用本地存储还是云存储
	//save file
	uploadPath := tools.NewCommon().Pwd() + "/public/upload/" + fileName
	err = ctx.SaveUploadedFile(file, uploadPath)
	if err != nil {
		tools.NewLogUtil().Error(err.Error())
		tools.NewResponse(ctx).JsonFailed("上传失败")
		return
	}

	tools.NewResponse(ctx).JsonSuccess("avatar/upload/" + fileName)
	return
}

/**
聊天
*/
func (r *chatController) ConnectSend(ctx *gin.Context) {

	//get token
	var (
		token  string
		err    error
		userID int64
	)
	paramsRequest := requests.ChatConnectSendRequest{}
	if err = ctx.ShouldBind(&paramsRequest); err != nil {
		tools.NewLogUtil().SugarPrint(err.Error())
		tools.NewResponse(ctx).JsonFailed("token参数错误")
		return
	}
	//validate token
	token = paramsRequest.Token
	tools.NewLogUtil().Info(token)
	userLoginRecord := models.UserLoginRecord{}
	db := tools.NewMysqlInstance()
	_, err = db.Where("token=?", token).And("deleted=0").Get(&userLoginRecord)
	if err != nil {
		tools.NewLogUtil().SugarPrint(err.Error())
		tools.NewResponse(ctx).JsonFailed("登录状态获取失败")
		return
	}
	userID = userLoginRecord.UserId

	isValid := true
	//第三方包	激活conn 同时判断
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isValid
		},
	}).Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		tools.NewLogUtil().Error("websocket 连接失败: ", err.Error())
		tools.NewResponse(ctx).JsonFailed("websocket 连接失败")
	}
	//only accept websocket
	if ok := websocket.IsWebSocketUpgrade(ctx.Request); !ok {
		tools.NewResponse(ctx).JsonFailed("仅限websocket")
	}
	//获取conn句柄
	node := &constants.NodeConstant{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe), //开启一个线程安全的set
	}
	//获取用户名下全部的群ID
	conconts := make([]models.ChatContact, 0)
	comIds := make([]int64, 0)

	_ = tools.NewMysqlInstance().Where("user_id = ? and type = ?", userID, 1).Find(&conconts)
	for _, v := range conconts {
		comIds = append(comIds, v.TargetId)
	}
	for _, v := range comIds {
		//将获取到的信息缓冲到set中
		node.GroupSets.Add(v)
	}
	//将userId和node做对应关系
	services.Rwlocker.Lock()
	services.ClientMap[userID] = node
	services.Rwlocker.Unlock()

	//发送逻辑
	go services.NewChatService().SendProc(node)

	//接收逻辑
	go services.NewChatService().RecvProc(node)

	//发送消息
	services.NewChatService().SendMsg(userID, []byte("哈哈哈哈哈"))
}
