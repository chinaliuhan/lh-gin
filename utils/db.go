package utils

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"lh-gin/models"
	"log"
	"strconv"
	"xorm.io/xorm"
)

//
//type DBMysqlUtil struct {
//	Model  struct{}
//	Engine *xorm.Engine
//}

//sync.Once能确保实例化对象Do方法在多线程环境只运行一次,内部通过互斥锁实现,它的内部本质上也是双重检查的方式

func init() {
	NewDBMysql()
}

var xormEngine *xorm.Engine

func NewDBMysql() *xorm.Engine {

	var (
		dsn string
		err error
		//DBMysql    *DBMysqlUtil
	)

	//初始化MySQL数据库
	handler := NewConfigUtil("db.ini")
	if handler == nil {
		return nil
	}
	dbConfig := handler.GetDbConfig("mysql")

	log.Println("db config: ", NewJsonUtils().Encode(dbConfig))

	//dsName = "root:root@(127.0.0.1:3306)/lh-moon?charset=utf8"

	dsn = fmt.Sprintf(
		"%s:%s@(%s:%s)/%s?charset=%s",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database, dbConfig.Charset,
	)
	logrus.Infoln("DB dsn: ", dsn)

	//这里的err比较特殊,最好处理一下 err.Error()的错误信息,防止出现意外
	xormEngine, err = xorm.NewEngine(dbConfig.Db, dsn)
	if err != nil && err.Error() != "" {
		log.Fatal("xorm NewEngine 初始化失败:", err.Error())
	}

	//数据库最大打开的连接数
	maxConn, _ := strconv.ParseInt(dbConfig.MaxConn, 10, 0)
	if maxConn > 0 {
		xormEngine.SetMaxOpenConns(int(maxConn))
		log.Println("设置最大连接数:", maxConn)
	}

	//是否显示SQL语句
	isShowSql, _ := strconv.ParseBool(dbConfig.IsShowSql)
	if isShowSql {
		xormEngine.ShowSQL(true)
		log.Println("开启SQL打印")
	}

	//是否同步表, todo 未实现自动化
	IsSync, _ := strconv.ParseBool(dbConfig.IsSync)
	if IsSync {
		log.Println("开启表结构同步")

		xormEngine.Sync2(new(models.User))
		xormEngine.Sync2(new(models.UserInfo))
		xormEngine.Sync2(new(models.ArticleContent))
	}

	return xormEngine
}

//var Db *xorm.Engine
//
//func init() {
//	var (
//		driverName string
//		dsName     string
//		err        error
//	)
//
//	//初始化MySQL数据库
//	dbConfig := NewConfigUtil("db.ini").GetDbConfig("mysql")
//	log.Println("db config: ", NewJsonUtils().Encode(dbConfig))
//	//dsName = "root:root@(127.0.0.1:3306)/lh-moon?charset=utf8"
//	dsName = fmt.Sprintf(
//		"%s:%s@(%s:%d)/%s?charset=%s",
//		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database, dbConfig.charset,
//	)
//	driverName = dbConfig.Db
//	Db, err = xorm.NewEngine(driverName, dsName)
//	if err != nil && err.Error() != "" {
//		log.Fatal(err.Error())
//	}
//	//数据库最大打开的连接数
//	Db.SetMaxOpenConns(10)
//
//	//是否显示SQL语句
//	Db.ShowSQL(true)
//
//	//自动同步struct中的表结构到DB
//	Db.Sync2(new(models.User), new(models.UserInfo))
//
//	println("Db xorm init success")
//}
