package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct {
	context gin.Context
}

func NewUserController() UserController {
	return UserController{}
}

func (this *UserController) Register(username string, password string, ga uint) (uint, error) {

	return 1, nil
}

func Demo(ctx *gin.Context) {

	ctx.String(http.StatusOK, "MMP")
}
