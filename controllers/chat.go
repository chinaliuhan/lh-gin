package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/fatih/set.v0"
	"lh-gin/constants"
	"lh-gin/models"
	"lh-gin/repositories"
	"lh-gin/requests"
	"lh-gin/services"
	"lh-gin/tools"
	"mime/multipart"
	"net/http"
	"sort"
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
func (r *chatController) RegisterAction(ctx *gin.Context) {
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

	paramsRequest := requests.ChatRegisterRequest{}
	err = ctx.ShouldBind(&paramsRequest)
	if err != nil {
		tools.NewLogUtil().Error(err.Error())
		tools.NewResponse(ctx).JsonFailed("参数错误")
		return
	}
	tmp := models.User{}
	_, _ = tools.NewMysqlInstance().Where("mobile=? and deleted=0", paramsRequest.Mobile).Get(&tmp)
	if tmp.Id > 0 {
		tools.NewResponse(ctx).JsonFailed("手机号已被占用")
		return
	}

	tmp.CreatedIp = ctx.ClientIP()
	tmp.UpdatedIp = ctx.ClientIP()

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
func (r *chatController) LoginAction(ctx *gin.Context) {
	//非POST请求直接返回模板
	if ctx.Request.Method != http.MethodPost {
		ctx.HTML(http.StatusOK, "/chat/login.shtml", nil)
		return
	}

	//prepare
	var (
		err              error
		serviceCode      int
		userInfo         models.User
		chatLoginRequest requests.ChatLoginRequest
	)

	//payload on json
	if err = ctx.ShouldBind(&chatLoginRequest); err != nil {
		tools.NewLogUtil().SugarPrint(err.Error())
		tools.NewResponse(ctx).JsonFailed("参数错误")
		return
	}

	//services
	userInfo, serviceCode = services.NewChatService().Login(&chatLoginRequest)
	token := tools.NewJwtUtil().GenerateToken(&tools.JWTClaims{
		UserID:   userInfo.Id,
		Nickname: userInfo.Nickname,
		Mobile:   userInfo.Mobile,
	})

	if serviceCode == constants.SERVICE_SUCCESS {
		//login history
		_, err = repositories.NewUserLoginRecordRepository().AddNew(userInfo.Id, token, ctx.ClientIP(), int(time.Now().Unix()))
		if err != nil {
			tools.NewLogUtil().Error(err)
			tools.NewResponse(ctx).JsonFailed("登录失败")
			return
		}

		//jwt token 2 cookie
		tools.NewCookieUtil(ctx).Set("token", token)

		//response json
		tokenMap := make(map[string]interface{}, 1)
		tokenMap["token"] = token //jwt token
		tokenMap["id"] = userInfo.Id
		tokenMap["avatar"] = userInfo.Avatar
		tokenMap["nickname"] = userInfo.Nickname
		tokenMap["memo"] = userInfo.Memo

		tools.NewResponse(ctx).JsonSuccess(tokenMap)
		return
	}
	switch serviceCode {
	case constants.SERVICE_FAILED:
		tools.NewLogUtil().Info("登录失败:", serviceCode)
		fallthrough
	case constants.SERVICE_NO_EXIST:
		tools.NewLogUtil().Info("登录失败:", serviceCode)
		fallthrough
	case constants.SERVICE_DELETED:
		tools.NewLogUtil().Info("登录失败:", serviceCode)
		fallthrough
	case constants.SERVICE_PASSWORD_ERROR:
		tools.NewResponse(ctx).JsonFailed("账号密码错误")
		return
	default:
		tools.NewLogUtil().Info("登录失败: 未捕捉到service code")
		tools.NewResponse(ctx).JsonFailed("系统繁忙请稍后重试")
		return
	}

}

/**
首页
*/
func (r *chatController) IndexAction(ctx *gin.Context) {
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
func (r *chatController) GetMyFriendListAction(ctx *gin.Context) {

	var (
		err error
	)

	//登录判断
	token, _ := tools.NewCookieUtil(ctx).Get(constants.JWT_TOKEN_KEY)
	if token == "" {
		tools.NewResponse(ctx).JsonFailed(constants.GetApiMsg(constants.API_CODE_NOT_LOGIN))

		return
	}
	jwtClaims := tools.NewJwtUtil().ParseToken(token)
	userID := jwtClaims.UserID
	if token == "" || userID <= 0 {
		tools.NewResponse(ctx).JsonFailed(constants.GetApiMsg(constants.API_CODE_NOT_LOGIN))
		return
	}

	db := tools.NewMysqlInstance()
	userContact := make([]models.ChatContact, 0)
	err = db.Where("user_id=? and type=?", userID, constants.CONNECT_TYPE_USER).Find(&userContact)
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
func (r *chatController) AddFriendAction(ctx *gin.Context) {

	//prepare
	var (
		err           error
		paramsRequest requests.ChatAddFriendRequest
	)

	//登录判断
	token, _ := tools.NewCookieUtil(ctx).Get(constants.JWT_TOKEN_KEY)
	if token == "" {
		tools.NewResponse(ctx).JsonFailed(constants.GetApiMsg(constants.API_CODE_NOT_LOGIN))

		return
	}
	jwtClaims := tools.NewJwtUtil().ParseToken(token)
	userID := jwtClaims.UserID
	if token == "" || userID <= 0 {
		tools.NewResponse(ctx).JsonFailed(constants.GetApiMsg(constants.API_CODE_NOT_LOGIN))
		return
	}

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
	db := tools.NewMysqlInstance()
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
		Type:      constants.CONNECT_TYPE_USER,
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
			Type:      constants.CONNECT_TYPE_USER,
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

func (r *chatController) GetMyChatRecord(ctx *gin.Context) {
	//登录判断
	token, _ := tools.NewCookieUtil(ctx).Get(constants.JWT_TOKEN_KEY)
	if token == "" {
		tools.NewResponse(ctx).JsonFailed(constants.GetApiMsg(constants.API_CODE_NOT_LOGIN))
		return
	}
	jwtClaims := tools.NewJwtUtil().ParseToken(token)
	userID := jwtClaims.UserID
	if token == "" || userID <= 0 {
		tools.NewResponse(ctx).JsonFailed(constants.GetApiMsg(constants.API_CODE_NOT_LOGIN))
		return
	}
	myFriends, err := services.NewChatService().GetMyFriends(userID)
	if err != nil {
		tools.NewResponse(ctx).JsonFailed("获取聊天记录失败")
		return
	}
	var friendsID []string
	for _, v := range myFriends {
		friendsID = append(friendsID, strconv.FormatInt(v.TargetId, 10))
	}

	myRecordList, err := services.NewChatService().GetChatRecordByUidAndTimeAndMyFriend(userID, 0, friendsID)
	if err != nil {
		tools.NewResponse(ctx).JsonFailed("获取聊天记录失败")
		return
	}
	targetRecordList, err := services.NewChatService().GetChatRecordByTargetAndTimeAndMyFriend(userID, 0, friendsID)
	if err != nil {
		tools.NewResponse(ctx).JsonFailed("获取聊天记录失败")
		return
	}
	recordList := make([]models.ChatRecord, 0)

	for _, v := range myRecordList {
		recordList = append(recordList, v)
	}
	for _, v := range targetRecordList {
		recordList = append(recordList, v)
	}

	recordIDList := []int{}
	for _, v := range recordList {
		recordIDList = append(recordIDList, v.Id)
	}

	//相当于PHP的多维数组排序
	sort.SliceStable(recordList, func(i, j int) bool { return recordList[i].Id < recordList[j].Id })

	//todo 数据未过滤
	tools.NewResponse(ctx).JsonSuccess(recordList)
	return
}

/**
获取群
*/
func (r *chatController) GetCommunityListAction(ctx *gin.Context) {
	//登录判断
	token, _ := tools.NewCookieUtil(ctx).Get(constants.JWT_TOKEN_KEY)
	if token == "" {
		tools.NewResponse(ctx).JsonFailed(constants.GetApiMsg(constants.API_CODE_NOT_LOGIN))
		return
	}
	jwtClaims := tools.NewJwtUtil().ParseToken(token)
	userID := jwtClaims.UserID
	if token == "" || userID <= 0 {
		tools.NewResponse(ctx).JsonFailed(constants.GetApiMsg(constants.API_CODE_NOT_LOGIN))
		return
	}

	contacts := make([]models.ChatContact, 0)
	comIds := make([]int64, 0)

	_ = tools.NewMysqlInstance().Where("user_id = ? and type = ? and deleted=0", userID, 2).Find(&contacts)
	for _, v := range contacts {
		comIds = append(comIds, v.TargetId)
	}
	coms := make([]models.ChatCommunity, 0)
	if len(comIds) == 0 {
		tools.NewResponse(ctx).JsonFailed("你没有加入群聊")
		return
	}
	_ = tools.NewMysqlInstance().In("id", comIds).Find(&coms)

	//todo 数据未过滤
	tools.NewResponse(ctx).JsonSuccess(coms)
	return
}

/**
创建群
*/
func (r *chatController) CreateCommunityAction(ctx *gin.Context) {

	var (
		err error
	)

	//prepare
	requestPrams := requests.ChatCreateCommunity{}
	if err := ctx.ShouldBind(&requestPrams); err != nil {
		tools.NewLogUtil().Error(err.Error())
		tools.NewResponse(ctx).JsonFailed("参数错误")
		return
	}

	//登录判断
	token, _ := tools.NewCookieUtil(ctx).Get(constants.JWT_TOKEN_KEY)
	if token == "" {
		tools.NewResponse(ctx).JsonFailed(constants.GetApiMsg(constants.API_CODE_NOT_LOGIN))

		return
	}
	jwtClaims := tools.NewJwtUtil().ParseToken(token)
	userID := jwtClaims.UserID
	if token == "" || userID <= 0 {
		tools.NewResponse(ctx).JsonFailed(constants.GetApiMsg(constants.API_CODE_NOT_LOGIN))
		return
	}
	db := tools.NewMysqlInstance()

	//create group
	if len(requestPrams.Name) == 0 {
		tools.NewResponse(ctx).JsonFailed("请输入群名称")
	}
	groupModel := models.ChatCommunity{}
	groupModel.UserId = userID

	countGroup, err := db.Count(&groupModel)
	if countGroup > 10 {
		tools.NewResponse(ctx).JsonFailed("一个用户最多最多创建10个群")
		return
	} else {
		dbSession := db.NewSession()
		_ = dbSession.Begin()
		groupModel.Created = int(time.Now().Unix())
		groupModel.CreatedIp = ctx.ClientIP()
		groupModel.Type = constants.CONNECT_TYPE_USER
		groupModel.Name = requestPrams.Name
		groupModel.Icon = requestPrams.Icon
		groupModel.Memo = requestPrams.Memo

		_, err = dbSession.InsertOne(&groupModel)
		if err != nil {
			_ = dbSession.Rollback()
			tools.NewLogUtil().Error(err)
			tools.NewResponse(ctx).JsonFailed("系统繁忙")
		}
		_, err := dbSession.InsertOne(
			models.ChatContact{
				UserId:   groupModel.UserId,
				TargetId: groupModel.Id,
				Type:     2,
				Created:  int(time.Now().Unix()),
			})
		if err != nil {
			_ = dbSession.Rollback()
			tools.NewLogUtil().Error(err)
			tools.NewResponse(ctx).JsonFailed("系统繁忙")
		}

		_ = dbSession.Commit()
	}

	tools.NewResponse(ctx).JsonSuccess(nil)
	return
}

/**
加入群
*/
func (r *chatController) JoinCommunityAction(ctx *gin.Context) {

	var (
		err error
	)

	//prepare
	requestPrams := requests.ChatJoinCommunity{}
	if err := ctx.ShouldBind(&requestPrams); err != nil {
		tools.NewLogUtil().Error(err.Error())
		tools.NewResponse(ctx).JsonFailed("参数错误")
		return
	}
	//登录判断
	token, _ := tools.NewCookieUtil(ctx).Get(constants.JWT_TOKEN_KEY)
	if token == "" {
		tools.NewResponse(ctx).JsonFailed(constants.GetApiMsg(constants.API_CODE_NOT_LOGIN))
		return
	}
	jwtClaims := tools.NewJwtUtil().ParseToken(token)
	userID := jwtClaims.UserID
	if token == "" || userID <= 0 {
		tools.NewResponse(ctx).JsonFailed(constants.GetApiMsg(constants.API_CODE_NOT_LOGIN))
		return
	}
	db := tools.NewMysqlInstance()

	//判断群是否存在
	chatCommunityModel := &models.ChatCommunity{}
	_, err = db.Where("id=? and deleted=0", requestPrams.TargetID).Get(chatCommunityModel)
	if err != nil || chatCommunityModel.Id <= 0 {
		tools.NewResponse(ctx).JsonFailed("不存在该群")
		return
	}
	//判断是否已入群
	chatContactModel := models.ChatContact{
		UserId:   userID,
		TargetId: requestPrams.TargetID,
		Type:     constants.CONNECT_TYPE_USER,
		Deleted:  0,
	}
	if _, err = db.Get(&chatContactModel); err == nil && chatContactModel.Id > 0 {
		tools.NewResponse(ctx).JsonFailed("已入群")
		return
	}

	//添加聊天关系
	chatContactModel.Created = int(time.Now().Unix())
	_, err = db.InsertOne(chatContactModel)
	if err != nil {
		tools.NewLogUtil().SugarPrint(err)
		tools.NewResponse(ctx).JsonFailed("入群失败")
		return
	}

	//取得node,添加gid到set
	services.Rwlocker.Lock()
	node, ok := services.ClientMap[userID]
	if ok {
		node.GroupSets.Add(requestPrams.TargetID)
	}
	//clientMap[userId] = node
	services.Rwlocker.Unlock()

	tools.NewResponse(ctx).JsonSuccess(nil)
	return
}

/**
上传头像
*/
func (r *chatController) UploadAction(ctx *gin.Context) {

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
	fileName = tools.NewGenerate().GenerateMd5(file.Filename)
	//todo 暂时存储,后面会用云存储,可通过配置使用本地存储还是云存储
	//save file
	uploadPath := tools.NewCommonUtil().Pwd() + "/public/upload/" + fileName
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
func (r *chatController) ConnectSendAction(ctx *gin.Context) {

	var (
		err error
	)
	//only accept websocket
	if ok := websocket.IsWebSocketUpgrade(ctx.Request); !ok {
		tools.NewResponse(ctx).JsonFailed("仅限websocket")
		return
	}

	//validate token
	token, _ := tools.NewCookieUtil(ctx).Get(constants.JWT_TOKEN_KEY)
	if token == "" {
		tools.NewResponse(ctx).JsonFailed(constants.GetApiMsg(constants.API_CODE_NOT_LOGIN))
		return
	}
	jwtClaims := tools.NewJwtUtil().ParseToken(token)
	userID := jwtClaims.UserID
	if token == "" || userID <= 0 {
		tools.NewResponse(ctx).JsonFailed(constants.GetApiMsg(constants.API_CODE_NOT_LOGIN))
		return
	}

	isValid := true
	//第三方websocket包,激活conn 同时判断
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isValid
		},
	}).Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		tools.NewLogUtil().Error("websocket 连接失败: ", err.Error())
		tools.NewResponse(ctx).JsonFailed("websocket 连接失败")
	}

	//获取conn句柄
	node := &constants.NodeConstant{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),   //并行转串行的队列,Conn 是一个IO型的资源 存在竞争关系
		GroupSets: set.New(set.ThreadSafe), //第三方包,可以快速获取 并集 交集 差集等,开启一个线程安全的set, 用以存储群组信息
	}
	//获取用户名下全部的群ID,只处理群即可,c2c的不用管,c2c存在一个Map中只发即可不需要群发
	chatContactModels := make([]models.ChatContact, 0)
	groupIDs := make([]int64, 0)

	_ = tools.NewMysqlInstance().Where("user_id = ? and type = ? ", userID, constants.CONNECT_TYPE_GROUP).Find(&chatContactModels)
	for _, v := range chatContactModels {
		groupIDs = append(groupIDs, v.TargetId)
	}
	for _, v := range groupIDs {
		//将获取到的信息缓冲到set中
		node.GroupSets.Add(v)
	}
	//将userId即连接者和node做对应关系,c2c存在此MAP中,直接发即可
	services.Rwlocker.Lock()
	services.ClientMap[userID] = node
	services.Rwlocker.Unlock()

	//发送逻辑
	go services.NewChatService().SendProcess(node)

	//接收逻辑
	go services.NewChatService().ReceiveProcess(node)

	//发送消息
	services.NewChatService().SendMsg(userID, []byte("welcome connect"))
}
