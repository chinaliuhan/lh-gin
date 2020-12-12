package requests

/**
登录参数绑定与验证
*/
type LoginRequest struct {
	Username string `json:"username" form:"username"  binding:"required,alphanum"` //required 为必须
	Password string `json:"password" form:"password"  binding:"required,min=6,max=32"`
	GaCode   string `json:"ga_code" form:"ga_code"  binding:"numeric"`
	//Mobile   string `json:"mobile" form:"mobile"  binding:"numeric"`
	//Email    string `json:"email" form:"email"  binding:"email"`
}

/**
注册绑定
*/
type RegisterRequest struct {
	Username string `json:"username" form:"username"  binding:"required,alphanum"` //required 为必须
	Password string `json:"password" form:"password"  binding:"required,min=6,max=32"`
	Mobile   string `json:"mobile" form:"mobile"  binding:"numeric"`
	Email    string `json:"email" form:"email"  binding:"email"`
	//WechatKey string `json:"wechat_key" form:"wechat_key"  binding:"alphanum"`
	//AppleKey  string `json:"apple_key" form:"apple_key"  binding:"alphanum"`
}
