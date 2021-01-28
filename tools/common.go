package tools

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"os"
	"reflect"
	"strings"
)

type CommonUtil struct {
}

func NewCommonUtil() *CommonUtil {
	return &CommonUtil{}
}

func (r *CommonUtil) Pwd() string {
	pwd, err := os.Getwd()
	if err != nil {
		log.Println("获取当前路径失败: ", err)
		return ""
	}

	return pwd
}

func (r *CommonUtil) Struct2Map(myStruct interface{}) map[string]interface{} {
	t := reflect.TypeOf(myStruct)
	v := reflect.ValueOf(myStruct)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

func (r *CommonUtil) StructPointer2Map2(myStruct *interface{}) map[string]interface{} {
	v := reflect.ValueOf(myStruct).Elem()
	typeOfType := v.Type()
	var data = make(map[string]interface{})
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		data[typeOfType.Field(i).Name] = field.Interface()
	}
	return data
}

func (r *CommonUtil) JsonEncode(data interface{}) (string, error) {

	bytes, err := json.Marshal(data)
	if err == nil {
		return string(bytes), nil
	}
	return "", err
}

func (r *CommonUtil) JsonDecode(jsonStr string, data interface{}) (interface{}, error) {

	err := json.Unmarshal([]byte(jsonStr), data)
	if err == nil {
		return data, nil
	}
	return nil, err
}

func (r *CommonUtil) Base64Encode(str string) string {
	strBytes := []byte(str)
	return base64.StdEncoding.EncodeToString(strBytes)
}

func (r *CommonUtil) Base64Decode(str string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", nil
	}
	return string(decoded), nil
}

func (r *CommonUtil) Array2String(arr []string, limit string) string {
	var str string
	for _, v := range arr { //遍历数组中所有元素追加成string
		str = str + limit + v
	}
	return strings.Trim(str, limit)
}
