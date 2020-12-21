package routers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"lh-gin/controllers"
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
