package config

type ThirdPartyConfig struct {
	BaiShun  BaiShunConfig  `mapstructure:"bai-shun"`
	HuanXin  HuanXinConfig  `mapstructure:"huan-xin"`
	AliCloud AliCloudConfig `mapstructure:"ali-cloud"`
}

type BaiShunConfig struct {
	AppID  string `mapstructure:"app-id"`
	AppKey string `mapstructure:"app-key"`
}

type AliCloudConfig struct {
	AccessId     string `mapstructure:"access-id"`
	AccessSecret string `mapstructure:"access-secret"`
}

type HuanXinConfig struct {
	ClientId     string `mapstructure:"client-id"`
	ClientSecret string `mapstructure:"client-secret"`
	AppKey       string `mapstructure:"app-key"`
	Url          string `mapstructure:"url"`
}

type PayConfig struct {
	Xendit   map[string]*XenditConfig `mapstructure:"xendit"`
	PayerMax PayerMaxConfig           `mapstructure:"payermax"`
	Coda     map[string]*CodaConfig   `mapstructure:"coda"`
	Binance  BinanceConfig            `mapstructure:"binance"`
}

type XenditConfig struct {
	SecretKey   string `mapstructure:"secret-key"`
	PublicKey   string `mapstructure:"public-key"`
	VerifyToken string `mapstructure:"verify-token"`
	BusinessId  string `mapstructure:"business-id"`
}

type PayerMaxConfig struct {
	AppID           string `mapstructure:"appid" json:"appID"`
	MerchantNo      string `mapstructure:"merchant-no" json:"merchantNo"`
	Url             string `mapstructure:"url" json:"url"`
	ReturnUrl       string `mapstructure:"return-url" json:"returnUrl"`
	NotifyUrl       string `mapstructure:"notify-url" json:"notifyUrl"`
	PayOutNotifyUrl string `mapstructure:"pay-out-notify-url" json:"payOutNotifyUrl"`
	RSAPublicKey    string `json:"rsaPublicKey"`
	RSAPrivateKey   string `json:"rsaPrivateKey"`
	//RSAPublicBytes  []byte `mapstructure:"-" json:"-"`
	//RSAPrivateBytes []byte `mapstructure:"-" json:"-"`
}

type CodaConfig struct {
	ApiKey   string `mapstructure:"api-key"`
	Country  int    `mapstructure:"country"`
	Currency int    `mapstructure:"currency"`
	Url      string `mapstructure:"url"`
}

type BinanceConfig struct {
	ApiKey    string `mapstructure:"api-key"`
	SecretKey string `mapstructure:"secret-key"`
	Url       string `mapstructure:"url"`
}
