package models

type ChatRecord struct {
	Id        int    `xorm:"not null pk autoincr comment('主键') INT(11)"`
	UserId    int    `xorm:"not null comment('用户ID') INT(11)"`
	TargetId  int    `xorm:"not null comment('目标ID') INT(11)"`
	Type      int    `xorm:"not null comment('类型(用户对用户,用户对群)') TINYINT(11)"`
	Remark    string `xorm:"not null default '' comment('备注') VARCHAR(255)"`
	Content   string `xorm:"not null comment('内容') TEXT"`
	Created   int    `xorm:"not null default 0 comment('创建时间') INT(15)"`
	Updated   int    `xorm:"not null default 0 comment('更新时间') INT(15)"`
	CreatedIp string `xorm:"not null default '' comment('创建IP') CHAR(15)"`
	UpdatedIp string `xorm:"not null default '' comment('更新IP') CHAR(15)"`
	Deleted   int    `xorm:"not null default 0 comment('删除时间') INT(15)"`
}

func (m *ChatRecord) TableName() string {
	return "chat_record"
}
