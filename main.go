package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"lh-gin/controllers"
	"net/http"
)

func main() {

	//newRouter := gin.New()
	//newRouter.Handle("GET", "/user/demo", controllers.Demo)
	//
	//_ = newRouter.Run()

	router := gin.Default()
	router.Handle(http.MethodGet, "/user/demo", controllers.Demo)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.GET("/user/register", func(context *gin.Context) {
		username := context.Query("username")
		password := context.Query("password")
		ga := context.Query("ga")
		//context.String(http.StatusOK, "Hello %s")
		context.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "successful",
			"data":    fmt.Sprintf("request params is : %s; %s; %s", username, password, ga),
		})
	})

	_ = router.Run() // listen and serve on 0.0.0.0:8080
}
