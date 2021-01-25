package tools

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type ResponseUtil struct {
	ctx     *gin.Context
	Code    int         `json:"code"`
	Message interface{} `json:"msg"`
	Data    interface{} `json:"data"`
}

//返回json信息
func NewResponse(ctx *gin.Context) *ResponseUtil {
	//设置响应的http数据类型
	ctx.Header("Content-Type", "application/json; charset=utf-8")

	return &ResponseUtil{ctx: ctx}
}

func (r ResponseUtil) JsonSuccess(data interface{}) {
	//拼装返回model,并转换为json
	h := ResponseUtil{
		Code:    0,
		Message: "successful",
		Data:    data,
	}
	jsonStr, err := json.Marshal(h)
	if err != nil {
		log.Println(err)
	}
	//将响应数据写入到返回中进行返回
	r.ctx.String(http.StatusOK, string(jsonStr))
}

func (r ResponseUtil) JsonFailed(message interface{}) {
	//拼装返回model,并转换为json
	h := ResponseUtil{
		Code:    1,
		Message: message,
		Data:    "",
	}
	jsonStr, err := json.Marshal(h)
	if err != nil {
		log.Println(err)
	}
	//将响应数据写入到返回中进行返回
	r.ctx.String(http.StatusOK, string(jsonStr))
}
