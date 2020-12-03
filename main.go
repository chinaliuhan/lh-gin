package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"lh-gin/middlewares"
	"lh-gin/routers"
	"lh-gin/utils"
)

func main() {

	/**
	预先声明必备变量,变量多的时候,方便追踪
	*/
	var (
		err          error
		serverConfig *utils.ServerConfig
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
	serverConfig = utils.NewConfigUtil("app.ini").GetServerConfig("server")
	logrus.Infoln("server config: ", utils.NewJsonUtils().Encode(serverConfig))

	dc := utils.NewConfigUtil("db.ini").GetDbConfig("mysql")
	logrus.Infoln("db.ini config: ", utils.NewJsonUtils().Encode(dc))

	/**
	静态资源
	*/
	staticPath := utils.NewCommon().Pwd() + "/public/static"
	logrus.Infoln(staticPath)
	engine.Static("/static", staticPath)
	assetsPath := utils.NewCommon().Pwd() + "/public/assets"
	engine.Static("/assets", assetsPath)

	/**
	执行GIN
	*/
	err = engine.Run(fmt.Sprintf("%s:%d", serverConfig.Address, serverConfig.Port))
	if err != nil {
		fmt.Println("gin 运行失败:", err.Error())
		panic(err)
	}
}
