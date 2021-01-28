package models

type ChatRecord struct {
	Id        int    `xorm:"not null pk autoincr comment('主键') INT(11)"`
	UserId    int64  `xorm:"not null comment('用户ID') INT(11)"`
	TargetId  int64  `xorm:"not null comment('目标ID') INT(11)"`
	Category  int    `xorm:"not null default 0 comment('消息分类(用户对用户,用户对群)') TINYINT(2)"`
	Mime      int    `xorm:"not null default 0 comment('消息类型(文字,图片,音频等)') TINYINT(2)"`
	Url       string `xorm:"not null default '' comment('表情或图片的URL地址') VARCHAR(255)"`
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
