package controllers

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"lh-gin/constants"
	"lh-gin/models"
	"lh-gin/requests"
	"lh-gin/services"
	"lh-gin/tools"
	"mime/multipart"
	"net/http"
	"strconv"
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
	if serviceCode == constants.SERVICE_SUCCESS {
		//session
		err = tools.NewSessionUtil(ctx).SetOne("user_id", info.Id)
		if err != nil {
			tools.NewResponse(ctx).JsonFailed("登录状态维持失败")
			tools.NewLogUtil().Error(err.Error())
			return
		}

		tools.NewResponse(ctx).JsonSuccess("")
		return
	}

	tools.NewResponse(ctx).JsonSuccess(nil)
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
	userID := tools.NewSessionUtil(ctx).GetOne("user_id")
	userID = userID.(int)

	db := tools.NewMysqlInstance()
	userContact := make([]models.ChatContact, 0)
	err := db.Where("user_id=?", userID).Find(&userContact)
	if err != nil {
		tools.NewLogUtil().SugarPrint(err.Error())
		tools.NewResponse(ctx).JsonFailed("系统繁忙")
		return
	}
	if len(userContact) == 0 {
		tools.NewResponse(ctx).JsonFailed("没有好友")
		return
	}
	targetIDList := make([]int, 0)
	for _, v := range userContact {
		targetIDList = append(targetIDList, v.TargetId)
	}

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

	//payload on json
	if err = ctx.ShouldBind(&paramsRequest); err != nil {
		tools.NewLogUtil().SugarPrint(err.Error())
		tools.NewResponse(ctx).JsonFailed("参数错误")
		return
	}

	userID := tools.NewSessionUtil(ctx).GetOne("user_id")
	userID = userID.(int)
	targetID, _ := strconv.Atoi(paramsRequest.TargetID)
	if userID == targetID {
		tools.NewResponse(ctx).JsonFailed("不能添加自己为好友")
		return
	}

	contactModel := models.ChatContact{}
	userModel := models.User{}
	db := tools.NewMysqlInstance()
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
		UserId:    userID.(int),
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
			TargetId:  userID.(int),
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
