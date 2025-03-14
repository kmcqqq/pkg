package utils

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

func StringToInt(s string) int {
	n, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0
	}
	return int(n)
}

func StringToInt64(s string) int64 {
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return int64(n)
}

func IntToString(n int) string {
	return strconv.Itoa(n)
}

func Int64ToString(n int64) string {
	return strconv.FormatInt(n, 10)
}

func StringToFloat64(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}

	return f
}

func Float64ToString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func Float64ToStringByMoney(f float64) string {
	intPart := math.Floor(f)
	if intPart == f {
		return Int64ToString(int64(intPart))
	} else {
		return Float64ToString(f)
	}
}

// 时间戳转日期
func FormatTimestampToDateTime(timestamp int64) time.Time {
	var date time.Time

	if timestamp > 1e12 {
		date = time.Unix(0, timestamp*int64(time.Millisecond))
	} else {
		date = time.Unix(timestamp, 0)
	}
	return date
}

func BankCardFormat(cardId string) string {
	cardId = strings.Replace(cardId, " ", "", -1)
	cardId = strings.Replace(cardId, "_", "", -1)
	cardId = strings.Replace(cardId, "-", "", -1)
	cardId = strings.Replace(cardId, "/", "", -1)
	cardId = strings.Replace(cardId, "+", "", -1)
	return cardId
}

func FormatStringToDateTime(layout, value string) (time.Time, error) {
	parsedTime, err := time.Parse(layout, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse date time: %w", err)
	}
	return parsedTime, nil
}
