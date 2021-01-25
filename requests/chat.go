package requests

/**
聊天注册
*/
type ChatRegisterRequest struct {
	Mobile     string `json:"mobile" form:"mobile"  binding:"required,alphanum"` //required 为必须
	Password   string `json:"password" form:"password"  binding:"required,min=6,max=32"`
	RePassword string `json:"rePassword" form:"rePassword"  binding:"required,min=6,max=32"`
	Memo       string `json:"memo" form:"memo"  binding:"required,min=6,max=32"`
	Avatar     string `json:"avatar" form:"avatar"  binding:"required,min=1,max=32"`
	Sex        string `json:"sex" form:"sex"  binding:"required,min=1,max=32"`
	Nickname   string `json:"nickname" form:"nickname"  binding:"required,min=2,max=32"`
}

/**
聊天登录
*/
type ChatLoginRequest struct {
	Mobile   string `json:"mobile" form:"mobile"  binding:"required,alphanum"` //required 为必须
	Password string `json:"password" form:"password"  binding:"required,min=6,max=32"`
}

/**
添加好友
*/
type ChatAddFriendRequest struct {
	TargetID string `json:"target_id" form:"target_id"  binding:"required,alphanum"` //required 为必须
}
