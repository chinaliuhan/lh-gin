package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"lh-gin/middlewares"
	"lh-gin/models"
	"lh-gin/routers"
	"lh-gin/tools"
)

func main() {

	/**
	预先声明必备变量,变量多的时候,方便追踪
	*/
	var (
		err          error
		serverConfig *tools.ServerConfig
	)

	/**
	启动Gin
	*/
	engine := gin.Default()

	/**
	检查cookie - 通过中间件
	*/
	engine.Use(middlewares.NewRequestMiddleware().CheckCookieAndInit)

	/**
	加载全部的view **代表目录,*代表文件
	*/
	viewsPath := tools.NewCommonUtil().Pwd() + "/public/views/**/*"
	engine.LoadHTMLGlob(viewsPath)
	//engine.Use(middlewares.NewRequestMiddleware().AutoExecView)

	/**
	自定义路由文件
	*/
	routers.ChatRouters(engine)
	routers.UserRouters(engine)
	routers.ArticleRouters(engine)
	routers.DemoRouters(engine)

	/**
	读取服务器配置文件
	*/
	serverConfig = tools.NewConfigUtil("app.ini").GetServerConfig("server")
	logrus.Infoln("server config: ", tools.NewJsonUtil().Encode(serverConfig))

	dc := tools.NewConfigUtil("db.ini").GetDbConfig("mysql")
	logrus.Infoln("db.ini config: ", tools.NewJsonUtil().Encode(dc))

	/**
	静态资源
	*/
	//staticPath := tools.NewCommonUtil().Pwd() + "/public/static"
	//engine.Static("/static", staticPath)
	assetsPath := tools.NewCommonUtil().Pwd() + "/public/assets"
	engine.Static("/assets", assetsPath)
	engine.Static("/chat/assets/", assetsPath)
	uploadPath := tools.NewCommonUtil().Pwd() + "/public/upload"
	engine.Static("/chat/avatar/upload", uploadPath)

	//为单个静态资源文件，绑定url
	favicon := tools.NewCommonUtil().Pwd() + "/public/assets/images/favicon.ico"
	engine.StaticFile("/favicon.ico", favicon)
	//indexHtml := tools.NewCommonUtil().Pwd() + "/public/static/login.html"
	//engine.StaticFile("/", indexHtml)

	/**
	执行GIN
	*/
	err = engine.Run(fmt.Sprintf("%s:%d", serverConfig.Address, serverConfig.Port))
	if err != nil {
		fmt.Println("gin 运行失败:", err.Error())
		panic(err)
	}
}

func init() {
	//同步表结构到MySQL todo
	_ = tools.NewMysqlInstance().Ping()
	_ = tools.NewMysqlInstance().Sync2(new(models.User))
	_ = tools.NewMysqlInstance().Sync2(new(models.UserInfo))
	_ = tools.NewMysqlInstance().Sync2(new(models.UserLoginRecord))
	_ = tools.NewMysqlInstance().Sync2(new(models.ArticleClassify))
	_ = tools.NewMysqlInstance().Sync2(new(models.ArticleContent))
	_ = tools.NewMysqlInstance().Sync2(new(models.ChatCommunity))
	_ = tools.NewMysqlInstance().Sync2(new(models.ChatContact))
	_ = tools.NewMysqlInstance().Sync2(new(models.ChatRecord))
}
