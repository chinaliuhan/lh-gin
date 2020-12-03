package utils

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"log"
)

type SessionUtil struct {
	SessionManager sessions.Session
	cookie         cookie.Store
}

//var globalSessions *session.Manager

func NewSessionUtil(ctx *gin.Context) *SessionUtil {

	sessionHandler := sessions.Default(ctx)
	sessionHandler.Set("name", "liuhao")

	err := sessionHandler.Save()
	if err != nil {
		log.Println("session保存失败")
	}

	//store := cookie.NewStore([]byte(cookieKey))

	return &SessionUtil{

		SessionManager: sessionHandler,
	}
}

func (receiver SessionUtil) GetAll() (string, error) {
	return "", nil
}

func (receiver SessionUtil) GetOne(key string) (string, error) {

	return "", nil
}

func (receiver SessionUtil) SetOne(key string, value string) (string, error) {

	return "", nil
}
