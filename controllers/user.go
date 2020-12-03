package controllers

import (
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"lh-gin/models"
	"lh-gin/requests"
	"lh-gin/services"
	"lh-gin/utils"
	"log"
)

type UserController struct {
	context *gin.Context
}

func NewUserController() *UserController {

	return &UserController{
	}
}

/**
注册
使用GIN自带验证器的自动验证
*/
func (r *UserController) Register(ctx *gin.Context) {

	//prepare
	var (
		err    error
		lastId int64
	)
	var tmp = models.User{}
	tmp.CreatedIp = ctx.ClientIP()
	tmp.UpdatedIp = ctx.ClientIP()

	//如果输入无效，则会写入400错误并在响应中设置Content-Type标头“ text / plain”。
	//err := ctx.Bind(&requests.RegisterRequest{})
	requestRegister := requests.RegisterRequest{}
	err = ctx.ShouldBind(&requestRegister) ////与c.Bind（）类似，但是此方法未将响应状态代码设置为400，并且如果json无效，则中止。
	if err != nil {
		log.Println(err.Error())
		utils.NewResponse(ctx).JsonFailed("参数错误")
		return
	}

	tmp.Username = requestRegister.Username
	tmp.Mobile = requestRegister.Mobile
	tmp.Mobile = requestRegister.Mobile
	tmp.Email = requestRegister.Email
	tmp.WechatKey = requestRegister.WechatKey
	tmp.AppleKey = requestRegister.AppleKey

	//services
	lastId, err = services.NewUserService().AddNew(tmp)
	if err != nil {
		log.Println(err.Error())
		utils.NewResponse(ctx).JsonFailed("注册失败")
		return
	}
	utils.NewResponse(ctx).JsonSuccess(lastId)
	return
}

/**
登录
使用手动验证的形式, todo 未验证数据
*/
func (r *UserController) Login(ctx *gin.Context) {

	//prepare
	var (
		ok           bool
		err          error
		username     string
		password     string
		info         models.User
		loginRequest requests.LoginRequest
	)

	// 普通方式获取参数, todo 暂时不做验证,有数据就行
	if username, ok = ctx.GetPostForm("username"); !ok {
		utils.NewResponse(ctx).JsonFailed("请输入账号")
	}

	if password, ok = ctx.GetPostForm("password"); !ok {
		utils.NewResponse(ctx).JsonFailed("请输入密码")
	}

	//services
	loginRequest = requests.LoginRequest{Username: username, Password: password}
	info, err = services.NewUserService().GetLogin(loginRequest)
	if err != nil {
		utils.NewResponse(ctx).JsonFailed("登录失败")
		return
	}

	//session
	session := sessions.Default(ctx)
	stringInfo, _ := json.Marshal(info)
	session.Set("loginstatus", stringInfo)
	err = session.Save()
	if err != nil {
		utils.NewResponse(ctx).JsonFailed("登录状态维持失败")
		log.Println(err.Error())
		return
	}

	utils.NewResponse(ctx).JsonSuccess(info)
	return
}

func (r *UserController) Info(ctx *gin.Context) {

}
