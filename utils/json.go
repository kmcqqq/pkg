package utils

import "encoding/json"

func Struct2Json(obj interface{}) (string, error) {
	str, err := json.Marshal(obj)
	return string(str), err
}

func Json2Struct(str string, obj interface{}) error {
	// 将json转为结构体
	err := json.Unmarshal([]byte(str), obj)
	return err
}

func IsJSONArray(s string) bool {
	var arr []interface{}
	err := json.Unmarshal([]byte(s), &arr)
	return err == nil
}
