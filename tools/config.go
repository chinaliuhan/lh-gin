package tools

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"log"
	"os"
	"runtime"
)

type ConfigUtil struct {
	path    string
	handler *ini.File
	info    interface{}
}

/**
构造函数
*/
//func NewConfigUtil(fileName string) *ConfigUtil {

//两个参数时, 第一个参数为文件名,第二个参数为路径
//一个参数时,为文件名
func NewConfigUtil(path ...string) *ConfigUtil {
	var (
		pwd      string
		filePath string
		err      error
		handler  *ini.File
		cu       *ConfigUtil
	)

	if len(path) > 1 {
		filePath = fmt.Sprintf("%s%s", path[1], path[0])
	} else {
		pwd, _ = os.Getwd()
		filePath = pwd + "/conf/" + path[0]
	}
	//初始化配置文件
	handler, err = ini.Load(filePath)
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		logrus.Infof("配置文件读取失败,文件地址: %s 文件: %s 行号: %d 错误信息: %s", filePath, file, line, err.Error())

		return nil
	}

	cu = &ConfigUtil{}
	cu.handler = handler
	cu.path = filePath

	return cu
}

func (r *ConfigUtil) GetConfig2Struct(title string, myStructPoint *interface{}) *interface{} {
	//判断配置是否加载成功
	if r.handler == nil {
		logrus.Infoln("handler不存在,配置文件读取失败:", r.path)
		return nil
	}
	//将配置文件映射到struct中
	if err := r.handler.Section(title).MapTo(myStructPoint); err != nil {
		log.Println("映射配置文件失败:", err.Error())
		return nil
	}

	return myStructPoint
}

/**
读取服务器配置信息
*/
type ServerConfig struct {
	Address string
	Port    int
}

func (r *ConfigUtil) GetServerConfig(title string) *ServerConfig {
	//判断配置是否加载成功
	if r.handler == nil {
		log.Println("handler不存在,配置文件读取失败:", r.path)
		return nil
	}
	//将配置文件映射到struct中
	sc := &ServerConfig{}
	if err := r.handler.Section(title).MapTo(sc); err != nil {
		log.Println("映射配置文件失败:", err.Error())
		return nil
	}

	return sc
}

/**
读取数据库配置
*/
type DbConfig struct {
	Db        string
	User      string
	Password  string
	Host      string
	Port      string
	Database  string
	Charset   string
	MaxConn   string
	IsShowSql string
}

//dsName = "root:root@(127.0.0.1:3306)/lh-moon?charset=utf8"
func (r *ConfigUtil) GetDbConfig(section string) *DbConfig {
	//判断配置是否加载成功
	if r.handler == nil {
		log.Println("handler不存在,配置文件读取失败:", r.path)
		return nil
	}
	dc := &DbConfig{}
	if err := r.handler.Section(section).MapTo(dc); err != nil {
		log.Println("映射配置文件失败:", err.Error())
		return nil
	}

	return dc
}

type Session struct {
	SaveType   string
	SessionKey string
	Path       string
	Domain     string
	MaxAge     int
	Secure     bool
	HttpOnly   bool
}

func (r *ConfigUtil) GetCookieConfig(section string) *Session {
	//判断配置是否加载成功
	if r.handler == nil {
		log.Println("handler不存在,配置文件读取失败:", r.path)
		return nil
	}
	s := &Session{}
	if err := r.handler.Section(section).MapTo(s); err != nil {
		log.Println("映射配置文件失败:", err.Error())
		return nil
	}

	return s
}

type RabbitMQ struct {
	UserName string
	Password string
	Host     string
	Port     int
	Virtual  string
}

func (r *ConfigUtil) GetRabbitMQConfig(section string) *RabbitMQ {
	//判断配置是否加载成功
	if r.handler == nil {
		log.Println("handler不存在,配置文件读取失败:", r.path)
		return nil
	}
	s := &RabbitMQ{}
	if err := r.handler.Section(section).MapTo(s); err != nil {
		log.Println("映射配置文件失败:", err.Error())
		return nil
	}

	return s
}
