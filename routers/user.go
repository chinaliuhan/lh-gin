package routers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"lh-gin/controllers"
)

/**
用户路由
*/
var sessionKey = "GinSessionID"

func UserRouters(engine *gin.Engine) *gin.RouterGroup {
	//注册session中间件
	store := cookie.NewStore([]byte("secret"))
	sessionName := sessionKey
	engine.Use(sessions.Sessions(sessionName, store))

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
	}

	return engineHandler
}

/**
文章路由
*/
func ArticleRouters(engine *gin.Engine) *gin.RouterGroup {

	//注册session中间件
	store := cookie.NewStore([]byte("secret"))
	sessionName := sessionKey
	engine.Use(sessions.Sessions(sessionName, store))

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
案例
*/
func DemoRouters(engine *gin.Engine) *gin.RouterGroup {

	//注册session中间件
	store := cookie.NewStore([]byte("secret"))
	sessionName := sessionKey
	engine.Use(sessions.Sessions(sessionName, store))

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
