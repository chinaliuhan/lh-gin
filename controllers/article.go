package controllers

import (
	"github.com/gin-gonic/gin"
	"lh-gin/models"
	"lh-gin/requests"
	"lh-gin/services"
	"lh-gin/utils"
	"log"
	"strconv"
	"time"
)

type ArticleController struct {
	context *gin.Context
}

func NewArticleController() *ArticleController {
	return &ArticleController{
	}
}

func (r *ArticleController) Info(ctx *gin.Context) {
	var (
		uid  int
		err  error
		info models.ArticleContent
	)

	uid = ctx.GetInt("uid")
	uid, err = strconv.Atoi(string(uid))
	if err != nil {
		log.Println("类型转换出错")
		utils.NewResponse(ctx).JsonFailed("参数错误")
	}

	info, _ = services.NewArticleService().GetInfoByUid(int(uid))
	utils.NewResponse(ctx).JsonSuccess(info)
}

func (r *ArticleController) Add(ctx *gin.Context) {

	//prepare
	var (
		err    error
		lastId int64
	)
	var tmp = models.ArticleContent{}
	tmp.CreatedIp = ctx.ClientIP()
	tmp.UpdatedIp = ctx.ClientIP()
	tmp.Updated = int(time.Now().Unix())
	tmp.Created = int(time.Now().Unix())

	//如果输入无效，则会写入400错误并在响应中设置Content-Type标头“ text / plain”。
	//err := ctx.Bind(&requests.RegisterRequest{})
	bindRequest := requests.AddArticleContentRequest{}
	err = ctx.ShouldBind(&bindRequest) ////与c.Bind（）类似，但是此方法未将响应状态代码设置为400，并且如果json无效，则中止。
	if err != nil {
		log.Println(err.Error())
		utils.NewResponse(ctx).JsonFailed("参数错误")
		return
	}

	tmp.UserId = bindRequest.UserId
	tmp.Classify = bindRequest.Classify
	tmp.Title = bindRequest.Title
	tmp.Content = bindRequest.Content

	//services
	lastId, err = services.NewArticleService().AddNew(tmp)
	if err != nil {
		log.Println(err.Error())
		utils.NewResponse(ctx).JsonFailed("添加失败")
		return
	}
	utils.NewResponse(ctx).JsonSuccess(lastId)
	return
}

func (r *ArticleController) Del(ctx *gin.Context) {

}

func (r *ArticleController) Modify(ctx *gin.Context) {

}
