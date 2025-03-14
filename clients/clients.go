package clients

import (
	"momo-server-go/pkg/aliyun"
	"momo-server-go/pkg/baishun"
	"momo-server-go/pkg/config"
	"momo-server-go/pkg/huanxin"
	"sync"
)

var (
	once     sync.Once
	instance *Clients
)

// Clients 用于管理所有第三方客户端
type Clients struct {
	BaiShun  *baishun.BSClient
	HuanXin  *huanxin.HxClient
	AliCloud *aliyun.AliCloudClient
}

// Initialize 初始化所有客户端
func Initialize(cfg *config.Config) error {
	var err error
	once.Do(func() {
		instance = &Clients{}
		// 初始化各个客户端
		if instance.BaiShun, err = baishun.NewBaiShunClient(&cfg.ThirdParty.BaiShun); err != nil {
			return
		}
		if instance.HuanXin, err = huanxin.NewHuanXinClient(&cfg.ThirdParty.HuanXin); err != nil {
			return
		}
		if instance.AliCloud, err = aliyun.NewAliCloudClient(&cfg.ThirdParty.AliCloud); err != nil {
			return
		}
	})
	return err
}

// GetClients 获取客户端实例
func GetClients() *Clients {
	return instance
}
