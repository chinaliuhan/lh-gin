package constants

const (
	API_CODE_SUCCESS   = 0
	API_CODE_FAILED    = 1
	API_CODE_NOT_LOGIN = 2
	API_CODE_LOCK      = 3
	API_CODE_EXIST     = 1000
	API_CODE_NO_EXIST  = 1001
)

func GetApiMsg(code int) string {

	codeList := map[int]string{
		API_CODE_SUCCESS:   "成功",
		API_CODE_FAILED:    "失败",
		API_CODE_NOT_LOGIN: "未登录",
		API_CODE_LOCK:      "已锁定",
		API_CODE_EXIST:     "已存在",
		API_CODE_NO_EXIST:  "不已存在",
	}

	return codeList[code]
}
