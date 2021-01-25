package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"lh-gin/tools"
)

type RequestMiddleware struct {
}

func NewRequestMiddleware() *RequestMiddleware {

	return &RequestMiddleware{}
}

//检查cookie是否存在,不存在则初始化cookie
func (receiver RequestMiddleware) CheckCookieAndInit(ctx *gin.Context) {
	tools.NewLogUtil().Info(fmt.Sprintf("请求header为: %s", ctx.Request.Header))
	key := "lh-gin"
	value, _ := tools.NewCookie(ctx).Get(key)
	if value != "" {
		return
	}
	uuid := tools.NewGenerate().GenerateUUID()
	tools.NewCookie(ctx).Set(key, uuid)
	tools.NewLogUtil().Info(fmt.Sprintf("为IP: %s  设置cookie: %s", ctx.ClientIP(), uuid))

}
