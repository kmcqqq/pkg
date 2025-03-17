package coda

import (
	"errors"
	"fmt"
	"github.com/kmcqqq/pkg/config"
	"github.com/kmcqqq/pkg/httpHelper"
	"github.com/kmcqqq/pkg/utils"
)

type Client struct {
	ApiKey   string
	Country  int
	Currency int
	Url      string
}

func NewClient(cfg *config.PayConfig) (map[string]*Client, error) {
	clients := make(map[string]*Client)
	for key, value := range cfg.Coda {
		clients[key] = &Client{
			ApiKey:   value.ApiKey,
			Country:  value.Country,
			Currency: value.Currency,
		}
	}
	return clients, nil
}

type PayOrder struct {
	OrderId string  `json:"orderId"`
	Coin    int     `json:"coin"`
	Amount  float64 `json:"amount"`
	UserIdx int64   `json:"userIdx"`
	PayType int     `json:"payType"`
}

type PayOrderResponse struct {
	InitResult struct {
		ResultCode int    `json:"resultCode"`
		ResultDesc string `json:"resultDesc"`
		TxnId      int64  `json:"txnId"`
	}
}

type RequestParam struct {
	InitRequest struct {
		Country  int    `json:"country"`
		Currency int    `json:"currency"`
		OrderID  string `json:"orderId"`
		APIKey   string `json:"apiKey"`
		PayType  int    `json:"payType"`
		Items    []struct {
			Code  string  `json:"code"`
			Name  string  `json:"name"`
			Price float64 `json:"price"`
			Type  int     `json:"type"`
		} `json:"items"`
		Profile struct {
			Entry []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			} `json:"entry"`
		} `json:"profile"`
	} `json:"initRequest"`
}

func (c *Client) PayOrder(p PayOrder) (*PayOrderResponse, error) {
	var order = RequestParam{
		InitRequest: struct {
			Country  int    `json:"country"`
			Currency int    `json:"currency"`
			OrderID  string `json:"orderId"`
			APIKey   string `json:"apiKey"`
			PayType  int    `json:"payType"`
			Items    []struct {
				Code  string  `json:"code"`
				Name  string  `json:"name"`
				Price float64 `json:"price"`
				Type  int     `json:"type"`
			} `json:"items"`
			Profile struct {
				Entry []struct {
					Key   string `json:"key"`
					Value string `json:"value"`
				} `json:"entry"`
			} `json:"profile"`
		}{
			Country:  c.Country,
			Currency: c.Currency,
			OrderID:  p.OrderId,
			APIKey:   c.ApiKey,
			PayType:  p.PayType,
			Items: []struct {
				Code  string  `json:"code"`
				Name  string  `json:"name"`
				Price float64 `json:"price"`
				Type  int     `json:"type"`
			}{
				{
					Code:  "Momo Live",
					Name:  fmt.Sprintf("%dCoin", p.Coin),
					Price: p.Amount,
					Type:  1,
				},
			},
			Profile: struct {
				Entry []struct {
					Key   string `json:"key"`
					Value string `json:"value"`
				} `json:"entry"`
			}{
				Entry: []struct {
					Key   string `json:"key"`
					Value string `json:"value"`
				}{
					{
						Key:   "user_id",
						Value: fmt.Sprintf("%d", p.UserIdx),
					},
					{
						Key:   "need_mno_id",
						Value: "yes",
					},
				},
			},
		},
	}

	res, err := httpHelper.Post(c.Url, order, nil)
	if err != nil {
		return nil, err
	}
	var response PayOrderResponse
	err = utils.Json2Struct(res, &response)
	if err != nil {
		return nil, err
	}

	if response.InitResult.ResultCode != 0 {
		return nil, errors.New(response.InitResult.ResultDesc)
	}

	return &response, nil
}

type PayNotify struct {
	OrderId    string `json:"OrderId"`
	TxnId      string `json:"TxnId"`
	ResultCode string `json:"ResultCode"`
	Checksum   string `json:"Checksum"`
	TotalPrice string `json:"TotalPrice"`
}

func (c *Client) ValidSign(txnId string, orderId string, resultCode string, checksum string) bool {
	data := fmt.Sprintf("%s%s%s%s", txnId, c.ApiKey, orderId, resultCode)

	return utils.MD5(data) == checksum
}
