package tools

import (
	"encoding/json"
)

type JsonUtil struct {
}

func NewJsonUtil() *JsonUtil {
	return &JsonUtil{}
}

func (receiver *JsonUtil) Encode(data interface{}) string {

	bytes, err := json.Marshal(data)
	if err == nil {
		return string(bytes)
	}
	return ""
}

func (receiver *JsonUtil) Decode(jsonStr string, data interface{}) interface{} {

	if err := json.Unmarshal([]byte(jsonStr), data); err == nil {
		return data
	}
	return nil
}
