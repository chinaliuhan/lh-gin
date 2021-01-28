package routers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"lh-gin/controllers"
	"lh-gin/middlewares"
	"lh-gin/tools"
)

/**
路由
*/

func UserRouters(engine *gin.Engine) *gin.RouterGroup {
	//注册session中间件
	store := cookie.NewStore([]byte("secret"))
	SessionConfig := tools.NewConfigUtil("app.ini").GetCookieConfig("session")
	engine.Use(sessions.Sessions(SessionConfig.SessionKey, store))

	//绑定路由
	engineHandler := engine.Group("/user/")
	{
		controllerHandler := controllers.NewUserController()
		//获取
		engineHandler.GET("info", controllerHandler.Info)
		//登录
		engineHandler.POST("login", controllerHandler.Login)
		//注册
		engineHandler.POST("register", controllerHandler.Register)

		//GA
		engineHandler.GET("generateGa", controllerHandler.GaSecret)
		engineHandler.GET("generateGaQrcode", controllerHandler.GaSecretQrcode)
		engineHandler.POST("gaBind", controllerHandler.GaBind)
		engineHandler.POST("gaUnbind", controllerHandler.GaUnbind)
	}

	return engineHandler
}

/**
文章路由
*/
func ArticleRouters(engine *gin.Engine) *gin.RouterGroup {

	//注册session中间件
	store := cookie.NewStore([]byte("secret"))
	SessionConfig := tools.NewConfigUtil("app.ini").GetCookieConfig("session")
	engine.Use(sessions.Sessions(SessionConfig.SessionKey, store))

	//绑定路由
	engineHandler := engine.Group("/article/")
	{
		controllerHandler := controllers.NewArticleController()
		//获取
		engineHandler.GET("info", controllerHandler.Info)
		//添加
		engineHandler.POST("add", controllerHandler.Add)
		//删除
		engineHandler.POST("del", controllerHandler.Del)
		//修改不
		engineHandler.POST("modify", controllerHandler.Modify)
	}

	return engineHandler
}

/**
聊天
*/
func ChatRouters(engine *gin.Engine) *gin.RouterGroup {
	controllerHandler := controllers.NewChatController()
	//注册 登录
	engine.GET("/", controllerHandler.LoginAction)
	engine.Any("/chat/register", controllerHandler.RegisterAction)
	engine.Any("/chat/login", controllerHandler.LoginAction)
	engine.POST("/chat/upload", controllerHandler.UploadAction)
	engine.GET("/chat/index", controllerHandler.IndexAction)

	//绑定路由
	engineHandler := engine.Group("/chat/", middlewares.NewRequestMiddleware().JWTTokenVerify)
	{
		controllerHandler := controllers.NewChatController()
		//注册
		//engineHandler.Any("register", controllerHandler.RegisterAction)
		//登录
		//engineHandler.Any("login", controllerHandler.LoginAction)
		//上传
		//engineHandler.POST("upload", controllerHandler.UploadAction)
		//首页
		//engineHandler.GET("index", controllerHandler.IndexAction)
		//获取好友列表
		engineHandler.POST("getMyFriendList", controllerHandler.GetMyFriendListAction)
		//添加好友
		engineHandler.POST("addFriend", controllerHandler.AddFriendAction)
		//获取群列表
		engineHandler.POST("getCommunityList", controllerHandler.GetCommunityListAction)
		//添加群
		engineHandler.POST("createCommunity", controllerHandler.CreateCommunityAction)
		//加入群
		engineHandler.POST("joinCommunity", controllerHandler.JoinCommunityAction)
		//聊天
		engineHandler.Any("connectSend", controllerHandler.ConnectSendAction)
		//聊天记录
		engineHandler.Any("getMyChatRecord", controllerHandler.GetMyChatRecord)
	}

	return engineHandler
}

/**
案例
*/
func DemoRouters(engine *gin.Engine) *gin.RouterGroup {

	//注册session中间件
	store := cookie.NewStore([]byte("secret"))
	SessionConfig := tools.NewConfigUtil("app.ini").GetCookieConfig("session")
	engine.Use(sessions.Sessions(SessionConfig.SessionKey, store))

	//路由组
	engineHandler := engine.Group("/demo/")
	{
		//添加
		engineHandler.GET("set", controllers.SetSession)
		//获取
		engineHandler.GET("get", controllers.GetSession)

	}

	return engineHandler
}
