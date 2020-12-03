package models

type User struct {
	Id          int    `xorm:"not null pk autoincr comment('主键ID') unique INT(11)"`
	Username    string `xorm:"not null default '' comment('自定义账户') CHAR(32)"`
	Mobile      string `xorm:"not null default '' comment('手机号') CHAR(15)"`
	Email       string `xorm:"not null default '' comment('邮箱') CHAR(32)"`
	WechatKey   string `xorm:"not null default '' comment('微信ID') VARCHAR(255)"`
	AppleKey    string `xorm:"not null default '' comment('苹果ID') VARCHAR(255)"`
	Nickname    string `xorm:"not null default '' comment('昵称') CHAR(12)"`
	Password    string `xorm:"not null default '' comment('登录密码') CHAR(32)"`
	LoginStatus int    `xorm:"not null default 0 comment('登录状态') TINYINT(2)"`
	Create      int    `xorm:"not null default 0 comment('创建时间') INT(15)"`
	Updated     int    `xorm:"not null default 0 comment('更新时间') INT(15)"`
	Deleted     int    `xorm:"not null default 0 comment('删除时间') INT(15)"`
	CreatedIp   string `xorm:"not null default '' comment('创建IP') CHAR(15)"`
	UpdatedIp   string `xorm:"not null default '' comment('更新IP') CHAR(15)"`
	Version     int    `xorm:"not null default 0 comment('版本') INT(11)"`
	Created     int    `xorm:"not null default 0 comment('创建时间') INT(15)"`
}
