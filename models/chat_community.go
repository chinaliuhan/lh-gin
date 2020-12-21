package models

type ChatCommunity struct {
	Id        int64  `xorm:"pk autoincr comment('主键') BIGINT(20)"`
	UserId    int    `xorm:"not null comment('群主ID') INT(11)"`
	Name      string `xorm:"not null default '' comment('群名称') CHAR(15)"`
	Icon      string `xorm:"not null default '' comment('群logo') VARCHAR(255)"`
	Type      int    `xorm:"not null comment('群类型,这里只有一种 通用群') TINYINT(4)"`
	Remark    string `xorm:"not null default '' comment('备注, 描述') VARCHAR(255)"`
	Created   int    `xorm:"not null default 0 comment('创建时间') INT(15)"`
	Updated   int    `xorm:"not null default 0 comment('更新时间') INT(15)"`
	CreatedIp string `xorm:"not null default '' comment('创建IP') CHAR(15)"`
	UpdatedIp string `xorm:"not null default '' comment('更新IP') CHAR(15)"`
	Deleted   int    `xorm:"not null default 0 comment('删除时间') INT(15)"`
}

func (m *ChatCommunity) TableName() string {
	return "chat_community"
}
