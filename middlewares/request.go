package middlewares

import (
	"github.com/gin-gonic/gin"
	"lh-gin/utils"
	"log"
)

type RequestMiddleware struct {
}

func NewRequestMiddleware() *RequestMiddleware {

	return &RequestMiddleware{}
}

//检查cookie是否存在,不存在则初始化cookie
func (receiver RequestMiddleware) CheckCookieAndInit(ctx *gin.Context) {
	key := "lh-gin"
	if value, err := utils.NewCookie(ctx).Get(key); err == nil {
		log.Printf("IP: %s  的cookie为: %s", ctx.ClientIP(), value)
		return
	}
	uuid := utils.NewGenerate().GenerateUUID()
	log.Printf("为IP: %s  设置cookie: %s", ctx.ClientIP(), uuid)

	utils.NewCookie(ctx).Set(key, uuid)
}
