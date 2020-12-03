package utils

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"log"
	"os"
	"runtime"
)

type configUtil struct {
	path    string
	handler *ini.File
	info    interface{}
}

/**
构造函数
*/
func NewConfigUtil(fileName string) *configUtil {
	var (
		pwd      string
		filePath string
		err      error
		handler  *ini.File
		cu       *configUtil
	)
	cu = &configUtil{}

	//初始化配置文件
	pwd, _ = os.Getwd()
	filePath = pwd + "/conf/" + fileName
	handler, err = ini.Load(filePath)

	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		logrus.Infof("配置文件读取失败,文件地址: %s 文件: %s 行号: %d 错误信息: %s", filePath, file, line, err.Error())

		return nil
	}

	cu.handler = handler
	cu.path = filePath

	return cu
}

func (receiver *configUtil) GetConfig2Struct(title string, myStructPoint *interface{}) *interface{} {
	//判断配置是否加载成功
	if receiver.handler == nil {
		logrus.Infoln("handler不存在,配置文件读取失败:", receiver.path)
		return nil
	}
	//将配置文件映射到struct中
	if err := receiver.handler.Section(title).MapTo(myStructPoint); err != nil {
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

func (receiver *configUtil) GetServerConfig(title string) *ServerConfig {
	//判断配置是否加载成功
	if receiver.handler == nil {
		log.Println("handler不存在,配置文件读取失败:", receiver.path)
		return nil
	}
	//将配置文件映射到struct中
	sc := &ServerConfig{}
	if err := receiver.handler.Section(title).MapTo(sc); err != nil {
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
	IsSync    string
	IsShowSql string
}

//dsName = "root:root@(127.0.0.1:3306)/lh-moon?charset=utf8"
func (receiver *configUtil) GetDbConfig(title string) *DbConfig {
	//判断配置是否加载成功
	if receiver.handler == nil {
		log.Println("handler不存在,配置文件读取失败:", receiver.path)
		return nil
	}
	dc := &DbConfig{}
	if err := receiver.handler.Section(title).MapTo(dc); err != nil {
		log.Println("映射配置文件失败:", err.Error())
		return nil
	}

	return dc
}
