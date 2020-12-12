package models

type ArticleContent struct {
	Id         int    `xorm:"not null pk autoincr comment('主键ID') unique INT(11)"`
	UserId     int    `xorm:"not null default 0 comment('UID') INT(11)"`
	ClassifyId int    `xorm:"not null default 0 comment('分类ID') INT(15)"`
	Classify   string `xorm:"not null default '' comment('分类') CHAR(15)"`
	Title      string `xorm:"not null default '' comment('标题') CHAR(15)"`
	Content    string `xorm:"not null comment('内容') TEXT"`
	Status     int    `xorm:"not null default 0 comment('状态') TINYINT(2)"`
	AdminId    int    `xorm:"not null default 0 comment('管理员ID') INT(11)"`
	Created    int    `xorm:"not null default 0 comment('创建时间') INT(15)"`
	Updated    int    `xorm:"not null default 0 comment('更新时间') INT(15)"`
	Deleted    int    `xorm:"not null default 0 comment('删除时间') INT(15)"`
	CreatedIp  string `xorm:"not null default '' comment('创建IP') CHAR(15)"`
	UpdatedIp  string `xorm:"not null default '' comment('更新IP') CHAR(15)"`
	Version    int    `xorm:"not null default 0 comment('版本') INT(11)"`
}

func (m *ArticleContent) TableName() string {
	return "article_content"
}
