package tools

import (
	"github.com/gin-gonic/gin"
)

type CookieUtil struct {
	ctx *gin.Context
}

func NewCookieUtil(ctx *gin.Context) *CookieUtil {
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

	return "", err
}

func (r CookieUtil) Set(key string, value string) {

	r.ctx.SetCookie(key, value, 3600*10, "/", "", false, false)
}
