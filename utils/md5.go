package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"
)

func MD5(data string) string {
	h := md5.New()
	h.Write([]byte(data))

	sign := hex.EncodeToString(h.Sum(nil))

	return sign
}

// md5签名
func VerySignWithMd5(data map[string]interface{}, secret string) string {
	//fmt.Println(data, secret)
	newKeys := make([]string, 0)
	for k, v := range data {
		if k == "sign" {
			continue
		} else if v == nil || v == "" {
			continue
		}
		newKeys = append(newKeys, k)
	}
	sort.Strings(newKeys)

	var dataStr string
	for _, v := range newKeys {
		dataStr += fmt.Sprintf("%v=%v&", v, data[v])
	}
	dataStr = dataStr[:len(dataStr)-1]
	dataStr += secret
	//fmt.Println(dataStr)
	sign := MD5(dataStr)

	return sign
}
