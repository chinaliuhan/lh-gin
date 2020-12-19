package tools

import (
	"errors"
	"github.com/gin-gonic/gin"
)

type CookieUtil struct {
	ctx *gin.Context
}

func NewCookie(ctx *gin.Context) *CookieUtil {
	return &CookieUtil{
		ctx: ctx,
	}
}

func (receiver CookieUtil) Get(key string) (string, error) {

	if cookie, err := receiver.ctx.Request.Cookie(key); err != nil {
		if cookie != nil {
			return cookie.Value, nil
		}
	}

	return "", errors.New("未获取到任何cookie")
}

func (receiver CookieUtil) Set(key string, value string) {

	receiver.ctx.SetCookie(key, value, 3600*1, "/", "", true, true)
}
