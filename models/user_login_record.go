package models

type UserLoginRecord struct {
	Id        int    `xorm:"not null pk autoincr comment('主键ID') unique INT(11)"`
	UserId    int64  `xorm:"not null default '' comment('UID') CHAR(15)"`
	Token     string `xorm:"not null default '' comment('token') VARCHAR(255)"`
	Created   int64  `xorm:"not null default 0 comment('创建时间') INT(15)"`
	Updated   int64  `xorm:"not null default 0 comment('更新时间') INT(15)"`
	CreatedIp string `xorm:"not null default '' comment('创建IP') CHAR(15)"`
	UpdatedIp string `xorm:"not null default '' comment('更新IP') CHAR(15)"`
	Deleted   int    `xorm:"not null default 0 comment('删除时间') INT(15)"`
	Version   int    `xorm:"not null default 0 comment('版本') INT(11)"`
}

func (m *UserLoginRecord) TableName() string {
	return "user_login_record"
}
