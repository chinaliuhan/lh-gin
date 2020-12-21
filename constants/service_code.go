package constants

const (
	SERVICE_SUCCESS        = 0
	SERVICE_FAILED         = 1
	SERVICE_NO_EXIST       = 10000
	SERVICE_DELETED        = 10001
	SERVICE_PASSWORD_ERROR = 10002
	SERVICE_GA_WRONG       = 10003
)

func GetServiceMsg(code int) string {

	codeList := map[int]string{
		SERVICE_SUCCESS:        "成功",
		SERVICE_FAILED:         "失败",
		SERVICE_NO_EXIST:       "不存在",
		SERVICE_DELETED:        "已被删除",
		SERVICE_PASSWORD_ERROR: "密码错误",
		SERVICE_GA_WRONG:       "GA输入错误",
	}

	return codeList[code]
}
