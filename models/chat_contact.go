package models

type ChatContact struct {
	Id        int64  `xorm:"pk autoincr comment('主键') BIGINT(20)"`
	UserId    int    `xorm:"not null comment('用户ID,记录所属人ID') INT(11)"`
	TargetId  int    `xorm:"not null comment('对方用户ID') INT(11)"`
	Type      int    `xorm:"not null comment('类型(用户对用户,用户对群)') TINYINT(4)"`
	Remark    string `xorm:"not null default '' comment('备注') VARCHAR(255)"`
	Created   int    `xorm:"not null default 0 comment('创建时间') INT(15)"`
	Updated   int    `xorm:"not null default 0 comment('更新时间') INT(15)"`
	CreatedIp string `xorm:"not null default '' comment('创建IP') CHAR(15)"`
	UpdatedIp string `xorm:"not null default '' comment('更新IP') CHAR(15)"`
	Deleted   int    `xorm:"not null default 0 comment('删除时间') INT(15)"`
}

func (m *ChatContact) TableName() string {
	return "chat_contact"
}
