package aliyun

import (
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/kmcqqq/pkg/config"
	"time"
)

type AliCloudClient struct {
	SlsClient sls.ClientInterface
}

const (
	slsEndPoint = "ap-southeast-5-intranet.log.aliyuncs.com"
	ossEndPoint = "oss-ap-southeast-5-internal.aliyuncs.com"
)

func NewAliCloudClient(cfg *config.AliCloudConfig) (*AliCloudClient, error) {
	provider := sls.NewStaticCredentialsProvider(cfg.AccessId, cfg.AccessSecret, "")
	slsClient := sls.CreateNormalInterfaceV2(slsEndPoint, provider)
	return &AliCloudClient{
		SlsClient: slsClient,
	}, nil
}

// QuerySlsLog 查询 sls 日志
func (a *AliCloudClient) QuerySlsLog(project, logStore string, from, to time.Time, query string) ([]map[string]string, error) {
	req := sls.GetLogRequest{
		From:  from.Unix(),
		To:    to.Unix(),
		Query: query,
	}
	resp, err := a.SlsClient.GetLogsV2(project, logStore, &req)
	return resp.Logs, err
}
