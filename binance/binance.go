package binance

import (
	"crypto/hmac"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/kmcqqq/pkg/config"
	"github.com/kmcqqq/pkg/httpHelper"
	"github.com/kmcqqq/pkg/logger"
	"github.com/kmcqqq/pkg/utils"
	"github.com/mitchellh/mapstructure"
	"reflect"
	"strings"
	"time"
)

type Client struct {
	ApiKey    string
	SecretKey string
	Url       string
	PublicKey *rsa.PublicKey
}

func NewClient(cfg *config.PayConfig) (*Client, error) {
	var client = Client{
		ApiKey:    cfg.Binance.ApiKey,
		SecretKey: cfg.Binance.SecretKey,
		Url:       cfg.Binance.Url,
	}
	publicKey, err := client.getPublicKey()
	if err != nil {
		logger.Error("error", logger.String("text", "binance getPublicKey error"), logger.Err(err))
	}
	client.PublicKey = publicKey

	return &client, nil
}

type PayOrderRequest struct {
	Env struct {
		TerminalType string `json:"terminalType"`
	} `json:"env"`
	MerchantTradeNo string         `json:"merchantTradeNo"`
	OrderAmount     float64        `json:"orderAmount"`
	Currency        string         `json:"currency"`
	Description     string         `json:"description"`
	GoodsDetails    []GoodsDetails `json:"goodsDetails"`
}

type GoodsDetails struct {
	GoodsType        string `json:"goodsType"`
	GoodsCategory    string `json:"goodsCategory"`
	GoodsName        string `json:"goodsName"`
	GoodsDetail      string `json:"goodsDetail"`
	ReferenceGoodsId string `json:"referenceGoodsId"`
}

type Response struct {
	Status       string      `json:"status"`
	Code         string      `json:"code"`
	Data         interface{} `json:"data"`
	ErrorMessage string      `json:"errorMessage"`
}

func (b *Client) Do(url string, req interface{}) (*Response, error) {
	payload, err := utils.Struct2Json(req)
	if err != nil {
		return nil, err
	}
	timestamp := time.Now().UnixMilli()
	nonce := utils.RandStringRunes(32)
	signature := b.signPayload(utils.Int64ToString(timestamp), nonce, payload)

	fullUrl := fmt.Sprintf("%s%s", b.Url, url)
	header := map[string]string{
		"BinancePay-Timestamp":      utils.Int64ToString(timestamp),
		"BinancePay-Nonce":          nonce,
		"BinancePay-Certificate-SN": b.ApiKey,
		"BinancePay-Signature":      signature,
	}
	res, err := httpHelper.Post(fullUrl, req, header)
	if err != nil {
		return nil, err
	}

	var response Response
	err = utils.Json2Struct(res, &response)
	return &response, err
}

type PayOrderResponse struct {
	PrepayId     string `json:"prepayId"`
	TerminalType string `json:"terminalType"`
	ExpireTime   int64  `json:"expireTime"`
	QrcodeLink   string `json:"qrcodeLink"`
	QrContent    string `json:"qrContent"`
	CheckoutUrl  string `json:"checkoutUrl"`
	Deeplink     string `json:"deeplink"`
	UniversalUrl string `json:"universalUrl"`
}

// PayOrder 下单
func (b *Client) PayOrder(req *PayOrderRequest) (*PayOrderResponse, error) {
	url := "/binancepay/openapi/v3/order"
	response, err := b.Do(url, req)
	if err != nil {
		return nil, err
	}

	if response.Status == "SUCCESS" {
		var payOrder PayOrderResponse
		if err := mapstructure.Decode(response.Data, &payOrder); err != nil {
			return nil, errors.New("response.Date format error")
		}

		return &payOrder, nil
	} else {
		return nil, errors.New(response.ErrorMessage)
	}
}

type QueryPayOrderRequest struct {
	PrepayId        string `json:"prepayId"`
	MerchantTradeNo string `json:"merchantTradeNo"`
}

type QueryOrderResult struct {
	MerchantId      int64  `json:"merchantId"`
	PrepayId        string `json:"prepayId"`
	TransactionId   string `json:"transactionId"`
	MerchantTradeNo string `json:"merchantTradeNo"`
	Status          string `json:"status"`
	Currency        string `json:"currency"`
	OrderAmount     string `json:"orderAmount"`
}

// PayOrderQuery 支付订单查询
func (b *Client) PayOrderQuery(req *QueryPayOrderRequest) (*QueryOrderResult, error) {
	url := "/binancepay/openapi/v2/order/query"
	response, err := b.Do(url, req)
	if err != nil {
		return nil, err
	}

	if response.Status == "SUCCESS" {
		var payOrder QueryOrderResult
		if err := mapstructure.Decode(response.Data, &payOrder); err != nil {
			return nil, errors.New("response.Date format error")
		}

		return &payOrder, nil
	} else {
		return nil, errors.New(response.ErrorMessage)
	}
}

type certificates struct {
	CertSerial string `json:"certSerial"`
	CertPublic string `json:"certPublic"`
}

// 获取公钥
func (b *Client) getPublicKey() (*rsa.PublicKey, error) {
	url := "/binancepay/openapi/certificates"
	response, err := b.Do(url, nil)
	if err != nil {
		return nil, err
	}

	if response.Status == "SUCCESS" {
		if rv := reflect.ValueOf(response.Data); rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
			if rv.Len() > 0 {
				var cert certificates
				if err := mapstructure.Decode(rv.Index(0).Interface(), &cert); err != nil {
					return nil, errors.New("response.Date format error")
				}
				publicKey, err := utils.LoadPublicKey(cert.CertPublic, false)
				if err != nil {
					return nil, err
				}
				return publicKey, nil
			}
		}
		return nil, errors.New("response.Date format error")
	} else {
		return nil, errors.New(response.ErrorMessage)
	}
}

type PayNotify struct {
	BizType   string `json:"bizType"`
	BizId     string `json:"bizId"`
	BizIdStr  string `json:"bizIdStr"`
	BizStatus string `json:"bizStatus"`
	Data      struct {
		MerchantTradeNo string  `json:"merchantTradeNo"`
		ProductType     string  `json:"productType"`
		ProductName     string  `json:"productName"`
		TransactTime    int64   `json:"transactTime"`
		TradeType       string  `json:"tradeType"`
		TotalFee        float64 `json:"totalFee"`
		Currency        string  `json:"currency"`
	} `json:"data"`
}

type PayNotifySign struct {
	Payload   string `json:"payload"`
	Timestamp string `json:"timestamp"`
	Nonce     string `json:"nonce"`
}

// VerifyWebhookSignature 验证webhook通知
func (b *Client) VerifyWebhookSignature(payload, timestamp, nonce, signature string) (bool, error) {
	signString := timestamp + "\n" + nonce + "\n" + payload + "\n"
	verify, err := utils.VerySignWithRsa(signString, signature, b.PublicKey)
	return verify, err
}

// 签名
func (b *Client) signPayload(timestamp, nonce, payload string) string {
	signString := timestamp + "\n" + nonce + "\n" + payload + "\n"
	h := hmac.New(sha512.New, []byte(b.SecretKey))
	h.Write([]byte(signString))
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}
