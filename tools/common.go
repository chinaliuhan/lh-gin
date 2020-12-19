package tools

import (
	"log"
	"os"
	"reflect"
)

type Common struct {
}

func NewCommon() *Common {
	return &Common{}
}

func (r *Common) Pwd() string {
	pwd, err := os.Getwd()
	if err != nil {
		log.Println("获取当前路径失败: ", err)
		return ""
	}

	return pwd
}

func (r *Common) Struct2Map(myStruct interface{}) map[string]interface{} {
	t := reflect.TypeOf(myStruct)
	v := reflect.ValueOf(myStruct)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

func (r *Common) StructPointer2Map2(myStruct *interface{}) map[string]interface{} {
	v := reflect.ValueOf(myStruct).Elem()
	typeOfType := v.Type()
	var data = make(map[string]interface{})
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		data[typeOfType.Field(i).Name] = field.Interface()
	}
	return data
}
