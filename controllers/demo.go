package controllers

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"lh-gin/tools"
	"net/http"
)

func SetSession(ctx *gin.Context) {

	key := ctx.Query("key")
	value := ctx.Query("value")

	sessionHandler := sessions.Default(ctx)
	sessionHandler.Set(key, value)
	err := sessionHandler.Save()
	if err != nil {
		fmt.Println("set failed", err)
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 1, "message": "set"})
}
func GetSession(ctx *gin.Context) {

	tools.NewMysqlInstance()

	key := ctx.Query("key")

	session := sessions.Default(ctx)
	//value := session.Get(key)
	keyvalue := session.Get(key)
	//读出来的结果是base64的
	value := session.Get("sessionID")

	fmt.Println(keyvalue, value)
	session.Clear()
	ctx.JSON(http.StatusOK, gin.H{"code": 1, "message": "get", "data": value, "keyData": keyvalue})
}
