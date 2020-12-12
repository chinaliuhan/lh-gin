package utils

import (
	"encoding/json"
)

type JsonUtils struct {
}

func NewJsonUtils() *JsonUtils {
	return &JsonUtils{}
}

func (receiver *JsonUtils) Encode(data interface{}) string {

	bytes, err := json.Marshal(data)
	if err == nil {
		return string(bytes)
	}
	return ""
}

func (receiver *JsonUtils) Decode(jsonStr string, data interface{}) interface{} {

	if err := json.Unmarshal([]byte(jsonStr), data); err == nil {
		return data
	}
	return nil
}
