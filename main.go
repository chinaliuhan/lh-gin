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
	engine := gin.New()

	/**
	检查cookie - 通过中间件
	*/
	engine.Use(middlewares.NewRequestMiddleware().CheckCookieAndInit)

	/**
	自定义路由文件
	*/
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
	staticPath := tools.NewCommon().Pwd() + "/public/static"
	engine.Static("/static", staticPath)
	assetsPath := tools.NewCommon().Pwd() + "/public/assets"
	engine.Static("/assets", assetsPath)

	//为单个静态资源文件，绑定url
	favicon := tools.NewCommon().Pwd() + "/public/assets/images/favicon.ico"
	engine.StaticFile("/favicon.ico", favicon)

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
	_, _ = tools.NewMysqlInstance().QueryString("select * from user limit 100")
	_ = tools.NewMysqlInstance().Sync2(new(models.User))
	_ = tools.NewMysqlInstance().Sync2(new(models.UserInfo))
	_ = tools.NewMysqlInstance().Sync2(new(models.ArticleClassify))
	_ = tools.NewMysqlInstance().Sync2(new(models.ArticleContent))
}
