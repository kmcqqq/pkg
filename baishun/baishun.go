package baishun

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"gitlab.bobbylive.cn/kongmengcheng/pkg/config"
	"gitlab.bobbylive.cn/kongmengcheng/pkg/utils"
)

type BSClient struct {
	appId  string
	appKey string
}

func NewBaiShunClient(cfg *config.BaiShunConfig) (*BSClient, error) {
	return &BSClient{
		appId:  cfg.AppID,
		appKey: cfg.AppKey,
	}, nil
}

func (b *BSClient) CompareSignature(sign string, appId int64, signatureNonce string, timestamp int64) bool {
	data := fmt.Sprintf("%s%s%d", signatureNonce, b.appKey, timestamp)
	h := md5.New()
	h.Write([]byte(data))

	generationSign := hex.EncodeToString(h.Sum(nil))

	if utils.Int64ToString(appId) == b.appId && sign == generationSign {
		return true
	} else {
		return false
	}

}
