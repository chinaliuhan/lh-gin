package middlewares

import (
	"github.com/gin-gonic/gin"
	"lh-gin/constants"
	"lh-gin/tools"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type RequestMiddleware struct {
}

func NewRequestMiddleware() *RequestMiddleware {

	return &RequestMiddleware{}
}

//检查cookie是否存在,不存在则初始化cookie
func (r *RequestMiddleware) CheckCookieAndInit(ctx *gin.Context) {
	//tools.NewLogUtil().Info(fmt.Sprintf("请求header为: %s", ctx.Request.Header))
	key := "lh-gin"
	value, _ := tools.NewCookieUtil(ctx).Get(key)
	if value != "" {
		return
	}
	markMap := make(map[string]string)
	markMap["timestamp"] = strconv.FormatInt(time.Now().Unix(), 10)
	markMap["ip"] = ctx.ClientIP()

	markString, _ := tools.NewCommonUtil().JsonEncode(markMap)
	markString = tools.NewCommonUtil().Base64Encode(markString)
	tools.NewCookieUtil(ctx).Set(key, markString)
	//tools.NewLogUtil().Info(fmt.Sprintf("为IP: %s  设置cookie: %s", ctx.ClientIP(), markString))
}

/**
通过URI自动加载模板
*/
func (r *RequestMiddleware) AutoExecView(ctx *gin.Context) {
	requestUri := ctx.Request.RequestURI
	if ok := strings.Contains(requestUri, "shtml"); ok {
		index := strings.Index(requestUri, "shtml")
		viewName := requestUri[0 : index+5]
		tools.NewLogUtil().Info("访问View URL Path:" + requestUri + " view name:" + viewName)
		ctx.HTML(http.StatusOK, viewName, nil)
	}
}

/**
JWT token维持
*/
func (r *RequestMiddleware) JWTTokenVerify(ctx *gin.Context) {
	//登录判断
	token, _ := tools.NewCookieUtil(ctx).Get(constants.JWT_TOKEN_KEY)
	if token == "" {
		tools.NewResponse(ctx).JsonFailed(constants.GetApiMsg(constants.API_CODE_NOT_LOGIN))
		ctx.Abort()
		return
	}
	jwtClaims := tools.NewJwtUtil().ParseToken(token)
	if jwtClaims.UserID <= 0 {
		tools.NewResponse(ctx).JsonFailed(constants.GetApiMsg(constants.API_CODE_NOT_LOGIN))
		ctx.Abort()
		return
	}
	//token是否过期
	if err := jwtClaims.Valid(); err != nil {
		tools.NewLogUtil().Info(err.Error())
		tools.NewResponse(ctx).JsonFailed(constants.GetApiMsg(constants.API_CODE_LOGIN_TIME_OUT))
		ctx.Abort()
		return
	}
	//刷新token
	tools.NewJwtUtil().Refresh(token)
}
