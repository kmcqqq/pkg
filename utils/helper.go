package utils

import "encoding/json"

// Contains 函数使用泛型来接受任意类型的切片和要查找的元素
func Contains[T comparable](arr []T, element T) bool {
	for _, v := range arr {
		if v == element {
			return true
		}
	}
	return false
}

func GetJsonValue(jsonData string, key string) (interface{}, error) {
	var result interface{}
	err := json.Unmarshal([]byte(jsonData), &result)
	if err != nil {
		return nil, err
	}
	return searchKey(result, key), nil
}

func searchKey(data interface{}, key string) interface{} {
	switch value := data.(type) {
	case map[string]interface{}:
		for k, v := range value {
			if k == key {
				return v
			}
			if result := searchKey(v, key); result != nil {
				return result
			}
		}
	case []interface{}:
		for _, item := range value {
			if result := searchKey(item, key); result != nil {
				return result
			}
		}
	}
	return nil
}
