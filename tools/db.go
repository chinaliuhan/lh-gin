package tools

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"strconv"
	"sync"
	"xorm.io/xorm"
)

//继承xorm
type mysqlInstance struct {
	*xorm.Engine
}

//连接缓存
//var instance map[string]*mysqlInstance
//var instance *mysqlInstance
//var ism sync.Map
var instanceMap map[string]*mysqlInstance

//safe map
//type instanceSafeMap struct {
//	instanceMap map[string]*mysqlInstance
//	sync.RWMutex
//}

//sync.Once能确保实例化对象Do方法在多线程环境只运行一次,内部通过互斥锁实现,它的内部本质上也是双重检查的方式
var lockOnce sync.Once

/**
通过可扩展单例获取MySQL连接
第一个参数为: Section name 必填
第二个参数为: file name  app.ini 可选
第三个参数为: folder name /home/liuhao/go/src/lh-gin/conf 可选
*/
func NewMysqlInstance(params ...string) *mysqlInstance {
	NewLogUtil().Info("获取可变参数的长度", len(params))

	//参数多的时候最好在顶层初始化
	var (
		section string
	)

	//同时指定了section, filename folder
	if len(params) > 2 {
		section = params[0]
	} else if len(params) > 1 {
		//仅指定了section filename
		section = params[0]

	} else {
		if len(params) > 0 {
			//仅仅指定了section
			section = params[0]

		} else {
			//未指定任何参数,采用默认值
			section = "mysql"
		}
	}
	NewLogUtil().Info("当前 section", section)

	lockOnce.Do(func() {
		NewLogUtil().Info("调用一次单例")
		if instanceMap == nil {
			instanceMap = make(map[string]*mysqlInstance)
		}
		//params... 会将参数打散,逐个传入
		instanceMap[section] = NewMysqlUtil(params...).GetConnect()

	})

	return instanceMap[section]
}

type MysqlUtil struct {
	sectionName string
	fileName    string
	folderName  string
	dbConfig    *DbConfig
}

/**
直接获取MySQL连接
第一个参数为: Section name 必填
第二个参数为: file name  app.ini 可选
第三个参数为: folder name /home/liuhao/go/src/lh-gin/conf 可选
*/
func NewMysqlUtil(params ...string) *MysqlUtil {

	//参数多的时候最好在顶层初始化
	var (
		err      error
		filename string
		folder   string
		section  string
	)

	//同时指定了section, filename folder
	if len(params) > 2 {
		section = params[0]
		filename = params[1]
		folder = params[2]

	} else if len(params) > 1 {
		//仅指定了section filename
		section = params[0]
		filename = params[1]

	} else {
		if len(params) > 0 {
			//仅仅指定了section
			section = params[0]

		} else {
			//未指定任何参数,采用默认值
			section = "mysql"
		}

		//自动识别配置文件
		folder = NewCommon().Pwd() + "/conf/"
		_, err = os.Stat(folder + "db.ini")
		if os.IsNotExist(err) {
			filename = "db.ini"
		} else {
			filename = "app.ini"
		}
	}
	NewLogUtil().Info(filename, folder)

	//初始化MySQL数据库
	configHandler := NewConfigUtil(filename, folder)
	if configHandler == nil {
		NewLogUtil().Warning("读取MySQL配置文件失败")
		return nil
	}
	dbConfig := configHandler.GetDbConfig(section)
	NewLogUtil().Info("db config: ", NewJsonUtil().Encode(dbConfig))

	return &MysqlUtil{
		sectionName: section,
		fileName:    filename,
		folderName:  folder,
		dbConfig:    dbConfig,
	}

}

func (r *MysqlUtil) GetConnect() *mysqlInstance {
	//参数多的时候最好在顶层初始化
	var (
		dsn        string
		err        error
		ormHandler *xorm.Engine
	)
	//dsName = "root:root@(127.0.0.1:3306)/lh-moon?charset=utf8"
	dsn = fmt.Sprintf(
		"%s:%s@(%s:%s)/%s?charset=%s",
		r.dbConfig.User, r.dbConfig.Password, r.dbConfig.Host, r.dbConfig.Port, r.dbConfig.Database, r.dbConfig.Charset,
	)
	NewLogUtil().Info("DB dsn: ", dsn)

	//这里的err比较特殊,最好处理一下 err.Error()的错误信息,防止出现意外
	ormHandler, err = xorm.NewEngine(r.dbConfig.Db, dsn)
	if err != nil && err.Error() != "" {
		NewLogUtil().Warning("xorm NewEngine 初始化失败:", err.Error())
		return nil
	}
	if ormHandler == nil {
		NewLogUtil().Warning("xorm NewEngine engine 初始化失败: engine 为 nil")
		return nil
	}

	//数据库最大打开的连接数
	maxConn, _ := strconv.ParseInt(r.dbConfig.MaxConn, 10, 0)
	if maxConn > 0 {
		maxConn, _ := strconv.ParseInt(r.dbConfig.MaxConn, 10, 0)
		ormHandler.SetMaxOpenConns(int(maxConn))
		log.Println("设置最大连接数:", maxConn)
		NewLogUtil().Info("设置最大连接数")
	}

	//是否显示SQL语句
	isShowSql, _ := strconv.ParseBool(r.dbConfig.IsShowSql)
	if isShowSql {
		ormHandler.ShowSQL(true)
		NewLogUtil().Info("开启SQL打印")
	}

	return &mysqlInstance{ormHandler}
}

//var xormEngine *xorm.Engine

/**
第一个参数为: Section name
第二个参数为: file name  app.ini
第三个参数为: folder name /home/liuhao/go/src/lh-gin/conf
*/
//func NewDBMysql(params ...string) *xorm.Engine {
//	log.Println("获取可变参数的长度", len(params))
//
//	//参数多的时候最好在顶层初始化
//	var (
//		dsn      string
//		err      error
//		filename string
//		folder   string
//		section  string
//		//DBMysql  *MysqlUtil
//	)
//
//	//同时指定了section, filename folder
//	if len(params) > 2 {
//		section = params[0]
//		filename = params[1]
//		folder = params[2]
//
//	} else if len(params) > 1 {
//		//仅指定了section filename
//		section = params[0]
//		filename = params[1]
//
//	} else {
//		if len(params) > 0 {
//			//仅仅指定了section
//			section = params[0]
//
//		} else {
//			//未指定任何参数,采用默认值
//			section = "mysql"
//		}
//
//		//自动识别配置文件
//		folder = NewCommon().Pwd() + "/conf/"
//		_, err = os.Stat(folder + "db.ini")
//		if os.IsNotExist(err) {
//			filename = "db.ini"
//		} else {
//			filename = "app.ini"
//		}
//	}
//	logrus.Println(filename, folder)
//
//	//初始化MySQL数据库
//	handler := NewConfigUtil(filename, folder)
//	if handler == nil {
//		return nil
//	}
//	dbConfig := handler.GetDbConfig(section)
//
//	log.Println("db config: ", NewJsonUtil().Encode(dbConfig))
//
//	//dsName = "root:root@(127.0.0.1:3306)/lh-moon?charset=utf8"
//	dsn = fmt.Sprintf(
//		"%s:%s@(%s:%s)/%s?charset=%s",
//		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database, dbConfig.Charset,
//	)
//	logrus.Infoln("DB dsn: ", dsn)
//
//	//这里的err比较特殊,最好处理一下 err.Error()的错误信息,防止出现意外
//	xormEngine, err = xorm.NewEngine(dbConfig.Db, dsn)
//	if err != nil && err.Error() != "" {
//		log.Fatal("xorm NewEngine 初始化失败:", err.Error())
//	}
//
//	//数据库最大打开的连接数
//	maxConn, _ := strconv.ParseInt(dbConfig.MaxConn, 10, 0)
//	if maxConn > 0 {
//		xormEngine.SetMaxOpenConns(int(maxConn))
//		log.Println("设置最大连接数:", maxConn)
//	}
//
//	//是否显示SQL语句
//	isShowSql, _ := strconv.ParseBool(dbConfig.IsShowSql)
//	if isShowSql {
//		xormEngine.ShowSQL(true)
//		log.Println("开启SQL打印")
//	}
//
//	return xormEngine
//}

/**
* 自动检测和创建表，这个检测是根据表的名字
* 自动检测和新增表中的字段，这个检测是根据字段名，同时对表中多余的字段给出警告信息
* 自动检测，创建和删除索引和唯一索引，这个检测是根据索引的一个或多个字段名，而不根据索引名称。因此这里需要注意，如果在一个有大量数据的表中引入新的索引，数据库可能需要一定的时间来建立索引。
* 自动转换varchar字段类型到text字段类型，自动警告其它字段类型在模型和数据库之间不一致的情况。
* 自动警告字段的默认值，是否为空信息在模型和数据库之间不匹配的情况
 */
//func SyncTable(tableStruct interface{}) error {
//	if err := xormEngine.Sync2(tableStruct); err != nil {
//
//		log.Println("表结构同步失败")
//		return err
//	}
//
//	return nil
//}

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
