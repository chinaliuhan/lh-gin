package controllers

import (
	"github.com/gin-gonic/gin"
	"lh-gin/constants"
	"lh-gin/models"
	"lh-gin/repositories"
	"lh-gin/requests"
	"lh-gin/services"
	"lh-gin/tools"
	"log"
)

type UserController struct {
	context *gin.Context
}

func NewUserController() *UserController {

	return &UserController{}
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
		tools.NewResponse(ctx).JsonFailed("参数错误")
		return
	}

	tmp.Username = requestRegister.Username
	tmp.Mobile = requestRegister.Mobile
	tmp.Mobile = requestRegister.Mobile
	tmp.Email = requestRegister.Email
	//tmp.WechatKey = requestRegister.WechatKey
	//tmp.AppleKey = requestRegister.AppleKey

	//services
	lastId, err = services.NewUserService().AddNew(tmp)
	if err != nil {
		log.Println(err.Error())
		tools.NewResponse(ctx).JsonFailed("注册失败")
		return
	}
	tools.NewResponse(ctx).JsonSuccess(lastId)
	return
}

/**
登录
使用手动验证的形式
*/
func (r *UserController) Login(ctx *gin.Context) {

	//prepare
	var (
		err          error
		serviceCode  int
		info         models.User
		loginRequest requests.LoginRequest
	)

	//payload on json
	if err = ctx.ShouldBindJSON(&loginRequest); err != nil {
		log.Println(err.Error())
		tools.NewResponse(ctx).JsonFailed("参数错误")
		return
	}

	//services
	info, serviceCode = services.NewUserService().GetLogin(loginRequest)
	if serviceCode == constants.SERVICE_SUCCESS {
		//session
		err = tools.NewSessionUtil(ctx).SetOne("user_id", info.Id)
		if err != nil {
			tools.NewResponse(ctx).JsonFailed("登录状态维持失败")
			log.Println(err.Error())
			return
		}

		tools.NewResponse(ctx).JsonSuccess("")
		return
	}
	switch serviceCode {
	case constants.SERVICE_FAILED:
		log.Println("登录失败:", serviceCode)
		fallthrough
	case constants.SERVICE_NO_EXIST:
		log.Println("登录失败:", serviceCode)
		fallthrough
	case constants.SERVICE_DELETED:
		log.Println("登录失败:", serviceCode)
		fallthrough
	case constants.SERVICE_PASSWORD_ERROR:
		tools.NewResponse(ctx).JsonFailed("账号密码错误")
		return
	default:
		log.Println("登录失败: 未捕捉到service code")
		tools.NewResponse(ctx).JsonFailed("系统繁忙请稍后重试")
		return
	}
}

/**
生成GA
*/
func (r UserController) GaSecret(ctx *gin.Context) {

	userID := tools.NewSessionUtil(ctx).GetOne("user_id")
	//断言, 如果成功,则转换为Int
	uid, ok := userID.(int)
	if !ok {
		tools.NewResponse(ctx).JsonFailed(constants.GetApiMsg(constants.API_CODE_NOT_LOGIN))
		return
	}
	info, err := repositories.NewUserManagerRepository().GetInfoByID(uid)
	if err != nil || info.Id <= 0 {
		tools.NewResponse(ctx).JsonFailed(constants.GetApiMsg(constants.API_CODE_NOT_LOGIN))
		return
	}

	gaSecret := tools.NewGoogleAuth().GetSecret()
	data := gin.H{"ga_secret": gaSecret, "name": info.Username}
	tools.NewResponse(ctx).JsonSuccess(data)
}

/**
获取带用户信息的GA
*/
func (r UserController) GaSecretQrcode(ctx *gin.Context) {
	gaSecret := tools.NewGoogleAuth().GetSecret()
	gaSecretQrcode := tools.NewGoogleAuth().GetQrcode("hahah", gaSecret)
	tools.NewResponse(ctx).JsonSuccess(gaSecretQrcode)
}

/**
绑定GA
*/
func (r UserController) GaBind(ctx *gin.Context) {

	tools.NewResponse(ctx).JsonSuccess("")
}

/**
获取登录后的个人信息
*/
func (r *UserController) Info(ctx *gin.Context) {

	userID := tools.NewSessionUtil(ctx).GetOne("user_id")
	//断言, 如果成功,则转换为Int
	uid, ok := userID.(int)
	if !ok {
		tools.NewResponse(ctx).JsonFailed(constants.GetApiMsg(constants.API_CODE_NOT_LOGIN))
		return
	}
	info, err := repositories.NewUserManagerRepository().GetInfoByID(uid)
	if err != nil {
		tools.NewResponse(ctx).JsonFailed(constants.GetApiMsg(constants.API_CODE_NOT_LOGIN))
		return
	}

	tools.NewResponse(ctx).JsonSuccess(info)
	return
}

/**
登录
使用手动验证的形式, todo 未验证数据
*/
//func (r *UserController) LoginBk(ctx *gin.Context) {
//
//	//prepare
//	var (
//		ok           bool
//		err          error
//		username     string
//		password     string
//		info         models.User
//		loginRequest requests.LoginRequest
//	)
//
//	// 普通方式获取参数, todo 暂时不做验证,有数据就行
//	if username, ok = ctx.GetPostForm("username"); !ok {
//		utils.NewResponse(ctx).JsonFailed("请输入账号")
//	}
//
//	if password, ok = ctx.GetPostForm("password"); !ok {
//		utils.NewResponse(ctx).JsonFailed("请输入密码")
//	}
//
//	//services
//	loginRequest = requests.LoginRequest{Username: username, Password: password}
//	info, err = services.NewUserService().GetLogin(loginRequest)
//	if err != nil {
//		utils.NewResponse(ctx).JsonFailed("登录失败")
//		return
//	}
//
//	//session
//	err = utils.NewSessionUtil(ctx).SetOne("user_id", info.Id)
//	if err != nil {
//		utils.NewResponse(ctx).JsonFailed("登录状态维持失败")
//		log.Println(err.Error())
//		return
//	}
//
//	utils.NewResponse(ctx).JsonSuccess("")
//	return
//}
