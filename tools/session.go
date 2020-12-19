package tools

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
)

type SessionOptions struct {
	sessions.Options
}

type SessionUtil struct {
	handler sessions.Session
}

//var globalSessions *session.Manager
func NewSessionUtil(ctx *gin.Context) *SessionUtil {
	handler := sessions.Default(ctx)
	if handler == nil {
		log.Println("session保存失败")
	}

	return &SessionUtil{

		handler: handler,
	}
}

func (r *SessionUtil) GetOne(key string) interface{} {
	if data := r.handler.Get(key); data != nil {
		return data
	}

	return nil
}

func (r *SessionUtil) SetOne(key string, value interface{}) error {
	r.handler.Set(key, value)
	if err := r.handler.Save(); err != nil {
		return err
	}

	return nil
}

func (r *SessionUtil) Del(key string) error {
	r.handler.Delete(key)
	if err := r.handler.Save(); err != nil {
		return err
	}

	return nil
}

func (r *SessionUtil) FlushAll() error {
	r.handler.Clear()
	if err := r.handler.Save(); err != nil {
		return err
	}

	return nil
}

func (r *SessionUtil) GetAll() error {
	//todo 未实现

	return nil
}
