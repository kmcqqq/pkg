package httpHelper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.bobbylive.cn/kongmengcheng/pkg/logger"
	"io"
	"net/http"
	"time"
)

// Post 请求
func Post(url string, data interface{}, header map[string]string) (string, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	client := &http.Client{}

	request, err := http.NewRequest("POST", url, bytes.NewReader(dataBytes))

	if err != nil {
		return "", err
	}

	request.Header.Add("Content-Type", "application/json;charset=utf-8")
	for key, val := range header {
		request.Header.Add(key, val)
	}
	client.Timeout = 1 * time.Minute

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)

	if err != nil {
		return "", err
	}

	logger.Info("http", logger.String("url", url), logger.String("method", "POST"), logger.Int("StatusCode", response.StatusCode), logger.String("req", string(dataBytes)), logger.String("header", fmt.Sprintf("%+v", header)), logger.String("resp", string(bodyBytes)))

	if response.StatusCode != 200 {
		return "", errors.New(response.Status)
	}

	return string(bodyBytes), nil
}

// GetHeader Get带 header 请求
func GetHeader(url string, header map[string]string) (string, error) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return "", err
	}

	for key, val := range header {
		request.Header.Add(key, val)
	}
	client.Timeout = 1 * time.Minute

	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	//zap.S().Infow("httpHelper", "url", url, "method", "GET", "header", "StatusCode", resp.StatusCode, fmt.Sprintf("%+v", header), "resp", string(body))
	logger.Info("http", logger.String("url", url), logger.String("method", "GET"), logger.Int("StatusCode", resp.StatusCode), logger.String("header", fmt.Sprintf("%+v", header)), logger.String("resp", string(body)))

	return string(body), nil
}

// Get 请求
func Get(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	//zap.S().Infow("httpHelper", "url", url, "method", "GET", "StatusCode", resp.StatusCode, "resp", string(body))
	logger.Info("http", logger.String("url", url), logger.String("method", "GET"), logger.Int("StatusCode", resp.StatusCode), logger.String("resp", string(body)))

	return string(body), nil
}

// http转发
func HttpTransform(url, method string, body io.Reader, header http.Header) (string, error) {
	client := &http.Client{}
	request, err := http.NewRequest(method, url, body)

	if err != nil {
		return "", err
	}

	request.Header = header

	client.Timeout = 1 * time.Minute

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)

	if err != nil {
		return "", err
	}

	data, _ := io.ReadAll(body)
	//zap.S().Infow("httpHelper", "url", url, "method", method, "StatusCode", response.StatusCode, "req", string(data), "header", fmt.Sprintf("%+v", header), "resp", string(bodyBytes))
	logger.Info("http", logger.String("url", url), logger.String("method", method), logger.Int("StatusCode", response.StatusCode), logger.String("req", string(data)), logger.String("header", fmt.Sprintf("%+v", header)), logger.String("resp", string(bodyBytes)))

	if response.StatusCode != 200 {
		return "", errors.New(response.Status)
	}

	return string(bodyBytes), nil
}

func PostForm(url string, data string) (string, error) {
	//fmt.Printf("request url %s, data : %+v", url, data)

	dataBytes := []byte(data)

	client := &http.Client{}

	request, err := http.NewRequest("POST", url, bytes.NewReader(dataBytes))

	if err != nil {
		return "", err
	}

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client.Timeout = 1 * time.Minute

	response, err := client.Do(request)

	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)

	if err != nil {
		return "", err
	}

	//fmt.Printf("response data : %+v", string(bodyBytes))
	//zap.S().Infow("httpHelper", "url", url, "method", "POST", "req", string(dataBytes), "resp", string(bodyBytes))

	return string(bodyBytes), nil
}
