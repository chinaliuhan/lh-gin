package tools

import (
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

func (r CookieUtil) Get(key string) (string, error) {
	cookie, err := r.ctx.Cookie(key)
	if err != nil {
		return "", err
	}

	if cookie != "" {
		return cookie, nil
	}

	return "", nil
}

func (r CookieUtil) Set(key string, value string) {

	r.ctx.SetCookie(key, value, 3600*1, "/", "", true, true)
}
