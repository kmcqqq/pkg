package payermax

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"gitlab.bobbylive.cn/kongmengcheng/pkg/config"
	"gitlab.bobbylive.cn/kongmengcheng/pkg/httpHelper"
	"gitlab.bobbylive.cn/kongmengcheng/pkg/utils"
	"time"
)

type Client struct {
	AppId           string
	MerchantNo      string
	PrivateKey      *rsa.PrivateKey
	PublicKey       *rsa.PublicKey
	ReturnUrl       string
	NotifyUrl       string
	PayOutNotifyUrl string
	Url             string
}

func NewClient(cfg *config.PayConfig) (*Client, error) {
	var client = Client{
		AppId:           cfg.PayerMax.AppID,
		MerchantNo:      cfg.PayerMax.MerchantNo,
		ReturnUrl:       cfg.PayerMax.ReturnUrl,
		NotifyUrl:       cfg.PayerMax.NotifyUrl,
		PayOutNotifyUrl: cfg.PayerMax.PayOutNotifyUrl,
		Url:             cfg.PayerMax.Url,
	}

	privateKey, err := utils.LoadPrivateKey(cfg.PayerMax.RSAPrivateKey, false)
	if err != nil {
		return nil, err
	}
	client.PrivateKey = privateKey

	publicKey, err := utils.LoadPublicKey(cfg.PayerMax.RSAPublicKey, false)
	if err != nil {
		return nil, err
	}
	client.PublicKey = publicKey

	return &client, nil
}

const (
	Order               = "orderAndPay"
	OrderQuery          = "orderQuery"
	Refund              = "refund"
	RefundQuery         = "refundQuery"
	OutPayOrder         = "paymentOrderPay"
	OutPayQuery         = "paymentOrderQry"
	CurrentBalanceQuery = "currentBalanceQuery"
)

type RequestData struct {
	Version     string      `json:"version"`
	KeyVersion  string      `json:"keyVersion"`
	RequestTime string      `json:"requestTime"`
	AppID       string      `json:"appId"`
	MerchantNo  string      `json:"merchantNo"`
	Data        interface{} `json:"data"`
}
type Response struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (c *Client) Do(orderInfo interface{}, payType string) (*Response, error) {
	var requestData = RequestData{
		Version:     "1.3",
		KeyVersion:  "1",
		RequestTime: getTimeStr(), //time.Now().Format("2006-01-02T15:04:05.999Z07:00"),
		AppID:       c.AppId,
		MerchantNo:  c.MerchantNo,
		Data:        orderInfo,
	}
	requestStr, err := utils.Struct2Json(requestData)
	if err != nil {
		return nil, err
	}
	sign, err := utils.SignRsa(requestStr, c.PrivateKey)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/%s", c.Url, payType)
	header := map[string]string{
		"sign": sign,
	}

	res, err := httpHelper.Post(url, requestData, header)

	var response Response
	err = utils.Json2Struct(res, &response)
	return &response, err
}

type PayOrder struct {
	OutTradeNo    string  `json:"outTradeNo"`
	Subject       string  `json:"subject"`
	TotalAmount   float64 `json:"totalAmount"`
	Currency      string  `json:"currency"`
	Country       string  `json:"country"`
	UserID        string  `json:"userId"`
	PaymentDetail struct {
		PaymentMethodType string `json:"paymentMethodType"`
		TargetOrg         string `json:"targetOrg"`
	} `json:"paymentDetail"`
	Language         string `json:"language"`
	Reference        string `json:"reference"`
	FrontCallbackURL string `json:"frontCallbackURL"`
	NotifyURL        string `json:"notifyUrl"`
	Integrate        string `json:"integrate"`
}

type PayOrderResponseData struct {
	OutTradeNo  string `mapstructure:"outTradeNo" json:"outTradeNo"`
	TradeToken  string `mapstructure:"tradeToken" json:"tradeToken"`
	Status      string `mapstructure:"status" json:"status"`
	RedirectUrl string `mapstructure:"redirectUrl" json:"redirectUrl"`
}

func (c *Client) PayOrder(pay *PayOrder) (*PayOrderResponseData, error) {
	pay.Language = "en"
	pay.FrontCallbackURL = c.ReturnUrl
	pay.NotifyURL = c.NotifyUrl
	pay.Integrate = "Hosted_Checkout"

	paymentType := pay.PaymentDetail.PaymentMethodType
	if paymentType == "CARD" || paymentType == "APPLEPAY" || paymentType == "GOOGLEPAY" {
		pay.PaymentDetail.TargetOrg = ""
	}

	res, err := c.Do(pay, Order)
	if err != nil {
		return nil, err
	}

	if res.Code == "APPLY_SUCCESS" {
		var payOrder PayOrderResponseData
		if err := mapstructure.Decode(res.Data, &payOrder); err != nil {
			return nil, errors.New("res.Date format error")
		}

		return &payOrder, nil
	}

	return nil, errors.New(res.Msg)
}

type PayQuery struct {
	OutTradeNo string `json:"outTradeNo"`
}

type PayQueryResponseData struct {
	OutTradeNo string `mapstructure:"outTradeNo" json:"outTradeNo"`
	TradeNo    string `mapstructure:"tradeNo" json:"tradeNo"`
	Status     string `mapstructure:"status" json:"status"`
	Trade      struct {
		Amount   string `mapstructure:"amount" json:"amount"`
		Currency string `mapstructure:"currency" json:"currency"`
	} `json:"trade"`
	Reference    string `mapstructure:"reference" json:"reference"`
	NotifyEmail  string `mapstructure:"notifyEmail" json:"notifyEmail"`
	NotifyPhone  string `mapstructure:"notifyPhone" json:"notifyPhone"`
	ResponseCode string `mapstructure:"responseCode" json:"responseCode"`
	ResponseMsg  string `mapstructure:"responseMsg" json:"responseMsg"`
}

func (c *Client) PayQuery(orderId string) (*PayQueryResponseData, error) {
	req := PayQuery{
		OutTradeNo: orderId,
	}
	res, err := c.Do(req, OrderQuery)
	if err != nil {
		return nil, err
	}

	if res.Code == "APPLY_SUCCESS" {
		var pay PayQueryResponseData
		if err := mapstructure.Decode(res.Data, &pay); err != nil {
			return nil, errors.New("res.Date format error")
		}

		return &pay, nil
	}

	return nil, errors.New(res.Msg)
}

type PayNotify struct {
	Code       string    `json:"code"`
	Msg        string    `json:"msg"`
	KeyVersion string    `json:"keyVersion"`
	AppID      string    `json:"appId"`
	MerchantNo string    `json:"merchantNo"`
	NotifyTime time.Time `json:"notifyTime"`
	NotifyType string    `json:"notifyType"`
	Data       struct {
		OutTradeNo     string  `json:"outTradeNo"`
		TradeToken     string  `json:"tradeToken"`
		TotalAmount    float64 `json:"totalAmount"`
		Currency       string  `json:"currency"`
		ChannelNo      string  `json:"channelNo"`
		ThirdChannelNo string  `json:"thirdChannelNo"`
		PaymentCode    string  `json:"paymentCode"`
		Country        string  `json:"country"`
		Status         string  `json:"status"`
		PaymentDetails []struct {
			PaymentMethodType string `json:"paymentMethodType"`
			TargetOrg         string `json:"targetOrg"`
		} `json:"paymentDetails"`
		Reference string `json:"reference"`
	} `json:"data"`
}

type RemitRequest struct {
	OutTradeNo string `json:"outTradeNo"`
	Country    string `json:"country"`
	Trade      struct {
		Amount   string `json:"amount"`
		Currency string `json:"currency"`
	} `json:"trade"`
	PayeeInfo struct {
		PaymentMethodType string `json:"paymentMethodType"`
		TargetOrg         string `json:"targetOrg"`
		AccountInfo       struct {
			AccountNo string `json:"accountNo"`
		} `json:"accountInfo"`
		Name struct {
			FullName string `json:"fullName"`
		} `json:"name"`
		PayeePhone string `json:"payeePhone"`
		BankInfo   struct {
			BankCode string `json:"bankCode"`
		}
	} `json:"payeeInfo"`
	Remark      string `json:"remark"`
	Reference   string `json:"reference"`
	NotifyURL   string `json:"notifyUrl"`
	NotifyEmail string `json:"notifyEmail"`
}

type RemitResponseData struct {
	OutTradeNo string `mapstructure:"outTradeNo" json:"outTradeNo"`
	TradeNo    string `mapstructure:"tradeNo" json:"tradeNo"`
	Status     string `mapstructure:"status" json:"status"`
}

func (c *Client) PayOut(p *RemitRequest) (*RemitResponseData, error) {
	p.NotifyURL = c.PayOutNotifyUrl
	res, err := c.Do(p, OutPayOrder)
	if err != nil {
		return nil, err
	}

	if res.Code == "APPLY_SUCCESS" {
		var pay RemitResponseData
		if err := mapstructure.Decode(res.Data, &pay); err != nil {
			return nil, errors.New("res.Date format error")
		}

		return &pay, nil
	}

	return nil, errors.New(res.Msg)
}

type PayOutNotify struct {
	Code       string `json:"code"`
	Msg        string `json:"msg"`
	KeyVersion string `json:"keyVersion"`
	AppID      string `json:"appId"`
	MerchantNo string `json:"merchantNo"`
	NotifyTime string `json:"notifyTime"`
	NotifyType string `json:"notifyType"`
	Data       struct {
		OutTradeNo string `json:"outTradeNo"`
		TradeNo    string `json:"tradeNo"`
		Status     string `json:"status"`
		Trade      struct {
			Amount   string `json:"amount"`
			Currency string `json:"currency"`
		} `json:"trade"`
		TransactionUtcTime string `json:"transactionUtcTime"`
		PayFinishTime      string `json:"payFinishTime"`
		BounceBackTime     string `json:"bounceBackTime"`
		ExpiryTime         string `json:"expiryTime"`
		RedeemCode         string `json:"redeemCode"`
		Source             struct {
			Amount       string `json:"amount"`
			Currency     string `json:"currency"`
			Fee          string `json:"fee"`
			FeeCurrency  string `json:"feeCurrency"`
			Tax          string `json:"tax"`
			TaxCurrency  string `json:"taxCurrency"`
			ExchangeRate string `json:"exchangeRate"`
		} `json:"source"`
		Destination struct {
			Amount       string `json:"amount"`
			Currency     string `json:"currency"`
			Fee          string `json:"fee"`
			FeeCurrency  string `json:"feeCurrency"`
			Tax          string `json:"tax"`
			TaxCurrency  string `json:"taxCurrency"`
			ExchangeRate string `json:"exchangeRate"`
		} `json:"destination"`
		Reference    string `json:"reference"`
		NotifyEmail  string `json:"notifyEmail"`
		NotifyPhone  string `json:"notifyPhone"`
		ResponseCode string `json:"responseCode"`
		ResponseMsg  string `json:"responseMsg"`
	} `json:"data"`
}

func getTimeStr() string {
	unixMilli := time.Now().UnixMilli()
	if unixMilli%10 == 0 {
		unixMilli++
	}
	t := time.UnixMilli(unixMilli)
	return t.Format("2006-01-02T15:04:05.999Z07:00")
}
