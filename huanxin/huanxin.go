package huanxin

import (
	"errors"
	"fmt"
	"gitlab.bobbylive.cn/kongmengcheng/pkg/config"
	"gitlab.bobbylive.cn/kongmengcheng/pkg/httpHelper"
	"gitlab.bobbylive.cn/kongmengcheng/pkg/utils"
	"strings"
)

type HxClient struct {
	clientId     string
	clientSecret string
	appKey       string
	baseUrl      string
}

func NewHuanXinClient(cfg *config.HuanXinConfig) (*HxClient, error) {
	return &HxClient{
		clientId:     cfg.ClientId,
		clientSecret: cfg.ClientSecret,
		appKey:       cfg.AppKey,
		baseUrl:      cfg.Url,
	}, nil
}

func (h *HxClient) GetOrgName() string {
	return strings.Split(h.appKey, "#")[0]
}

func (h *HxClient) GetAppName() string {
	return strings.Split(h.appKey, "#")[1]
}

func (h *HxClient) GetAccessToken() (string, error) {
	url := fmt.Sprintf("%s/%s/%s/token", h.baseUrl, h.GetOrgName(), h.GetAppName())
	data := map[string]string{
		"grant_type":    "client_credentials",
		"client_id":     h.clientId,
		"client_secret": h.clientSecret,
	}
	res, err := httpHelper.Post(url, data, nil)
	if err != nil {
		return "", err
	}

	if strings.Contains(res, "error") {
		return "", errors.New(res)
	}

	accessToken, err := utils.GetJsonValue(res, "access_token")

	if err != nil {
		return "", err
	}
	return accessToken.(string), nil
}

func (h *HxClient) RegisterUser(userName string, password string, accessToken string) (string, error) {
	url := fmt.Sprintf("%s/%s/%s/users", h.baseUrl, h.GetOrgName(), h.GetAppName())
	data := map[string]string{
		"username": userName,
		"password": password,
		"nickname": userName,
	}

	headers := map[string]string{
		"Authorization": "Bearer " + accessToken,
	}
	res, err := httpHelper.Post(url, data, headers)
	if err != nil {
		return "", err
	}

	if strings.Contains(res, "error") {
		return "", errors.New(res)
	}

	uuId, err := utils.GetJsonValue(res, "uuid")
	return uuId.(string), err
}
