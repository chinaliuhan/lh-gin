package models

type UserInfo struct {
	Id                int    `xorm:"not null pk autoincr comment('主键ID') INT(11)"`
	UserId            string `xorm:"not null default '' comment('UID') unique CHAR(15)"`
	Sex               int    `xorm:"not null default 0 comment('性别') TINYINT(2)"`
	Age               int    `xorm:"not null default 0 comment('年龄') TINYINT(2)"`
	Adder             string `xorm:"not null default '' comment('地址') CHAR(32)"`
	IdentityName      string `xorm:"not null default '' comment('身份证姓名') unique CHAR(12)"`
	IdentityNumber    string `xorm:"not null default '' comment('身份证号') unique CHAR(20)"`
	IdentityCardFront string `xorm:"not null default '' comment('身份证正面') VARCHAR(255)"`
	IdentityCardEnd   string `xorm:"not null default '' comment('身份证背面') VARCHAR(255)"`
	IdentityVideo     string `xorm:"not null default '' comment('认证视频') VARCHAR(255)"`
	Created           int    `xorm:"not null default 0 comment('创建时间') INT(15)"`
	Updated           int    `xorm:"not null default 0 comment('更新时间') INT(15)"`
	Deleted           int    `xorm:"not null default 0 comment('删除时间') INT(15)"`
	CreatedIp         string `xorm:"not null default '' comment('创建IP') CHAR(15)"`
	UpdatedIp         string `xorm:"not null default '' comment('更新IP') CHAR(15)"`
	Version           int    `xorm:"not null default 0 comment('版本') INT(11)"`
	Create            int    `xorm:"not null default 0 comment('创建时间') INT(15)"`
}

func (m *UserInfo) TableName() string {
	return "user_info"
}
