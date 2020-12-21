package models

type UserChat struct {
	Id            int    `xorm:"not null pk autoincr comment('主键') INT(11)"`
	UserId        int    `xorm:"not null comment('用户ID') INT(11)"`
	Avatar        string `xorm:"not null default '' comment('头像') VARCHAR(255)"`
	Nickname      string `xorm:"not null default '' comment('昵称') VARCHAR(255)"`
	IsOnline      int    `xorm:"not null default 0 comment('是否在线') TINYINT(4)"`
	ChatStatus    int    `xorm:"not null default 0 comment('聊天状态') TINYINT(4)"`
	SendStatus    int    `xorm:"not null default 0 comment('发送状态') TINYINT(11)"`
	ReceiveStatus int    `xorm:"not null default 0 comment('接收状态') TINYINT(4)"`
	Created       int    `xorm:"not null default 0 comment('创建时间') INT(15)"`
	Updated       int    `xorm:"not null default 0 comment('更新时间') INT(15)"`
	CreatedIp     string `xorm:"not null default '' comment('创建IP') CHAR(15)"`
	UpdatedIp     string `xorm:"not null default '' comment('更新IP') CHAR(15)"`
	Deleted       int    `xorm:"not null default 0 comment('删除时间') INT(15)"`
}

func (m *UserChat) TableName() string {
	return "user_chat"
}
