package utils

import (
	"encoding/hex"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

// RandomNumber 随机整数 [min,max)
func RandomNumber(min int, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	num := r.Intn(max-min) + min
	return num
}

// 随机字符串生成
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:length], nil
}

// RandStringRunes 返回随机字符串
func RandStringRunes(n int) string {
	var letterRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)

	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandNumberRunes(n int) string {
	var letterRunes = []rune("1234567890")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)

	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GenerateRequestId() string {
	return uuid.New().String()
}
