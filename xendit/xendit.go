package xendit

import (
	"errors"
	"fmt"
	"github.com/kmcqqq/pkg/config"
	"github.com/kmcqqq/pkg/httpHelper"
	"github.com/kmcqqq/pkg/utils"
	"time"
)

type Client struct {
	SecretKey   string
	PublicKey   string
	VerifyToken string
	BusinessId  string
}

func NewClient(cfg *config.PayConfig) (map[string]*Client, error) {
	clients := make(map[string]*Client)
	for key, value := range cfg.Xendit {
		clients[key] = &Client{
			SecretKey:   value.SecretKey,
			PublicKey:   value.PublicKey,
			VerifyToken: value.VerifyToken,
			BusinessId:  value.BusinessId,
		}
	}
	return clients, nil
}

type InvoiceRequest struct {
	ExternalId     string   `json:"external_id"`
	Amount         float64  `json:"amount"`
	Currency       string   `json:"currency"`
	Description    string   `json:"description"`
	PaymentMethods []string `json:"payment_methods"`
	//Fees            []Fees   `json:"fees"`
	InvoiceDuration uint `json:"invoice_duration"`
}

type InvoiceResponse struct {
	ID           string    `json:"id"`
	UserId       string    `json:"user_id"`
	ExternalID   string    `json:"external_id"`
	Status       string    `json:"status"`
	Amount       float64   `json:"amount"`
	ExpiryDate   time.Time `json:"expiry_date"`
	InvoiceURL   string    `json:"invoice_url"`
	ErrorCode    string    `json:"error_code"`
	ErrorMessage string    `json:"message"`
}

func (x *Client) CreateInvoice(invoice InvoiceRequest) (*InvoiceResponse, error) {
	url := "https://api.xendit.co/v2/invoices"
	header := map[string]string{
		"Authorization": fmt.Sprintf("Basic %s", utils.EncodeStr2Base64(x.SecretKey)),
	}
	result, err := httpHelper.Post(url, invoice, header)
	if err != nil {
		return nil, err
	}

	var response *InvoiceResponse
	err = utils.Json2Struct(result, &response)

	if response.ErrorCode != "" {
		return nil, errors.New(response.ErrorMessage)
	}

	if response.Status == "PENDING" && response.ID != "" {
		return response, nil
	}

	return nil, errors.New(response.Status)
}

func (x *Client) GetInvoice(orderId string) (*InvoiceResponse, error) {
	url := "https://api.xendit.co/v2/invoices"
	header := map[string]string{
		"Authorization": fmt.Sprintf("Basic %s", utils.EncodeStr2Base64(x.SecretKey)),
	}
	result, err := httpHelper.GetHeader(fmt.Sprintf("%s/?external_id=%s", url, orderId), header)
	if err != nil {
		return nil, err
	}

	var res InvoiceResponse
	if utils.IsJSONArray(result) {
		var arr []InvoiceResponse
		err = utils.Json2Struct(result, &arr)
		if len(arr) > 0 {
			res = arr[0]
		}
	} else {
		err = utils.Json2Struct(result, &res)
	}
	return &res, err
}

type InvoiceNotify struct {
	ID                     string  `json:"id"`
	ExternalID             string  `json:"external_id"`
	UserID                 string  `json:"user_id"`
	Status                 string  `json:"status"`
	MerchantName           string  `json:"merchant_name"`
	Amount                 float64 `json:"amount"`
	PayerEmail             string  `json:"payer_email"`
	Description            string  `json:"description"`
	FeesPaidAmount         float64 `json:"fees_paid_amount"`
	AdjustedReceivedAmount float64 `json:"adjusted_received_amount"`
	BankCode               string  `json:"bank_code"`
	RetailOutletName       string  `json:"retail_outlet_name"`
	EwalletType            string  `json:"ewallet_type"`
	OnDemandLink           string  `json:"on_demand_link"`
	RecurringPaymentID     string  `json:"recurring_payment_id"`
	PaidAmount             float64 `json:"paid_amount"`
	Updated                string  `json:"updated"`
	Created                string  `json:"created"`
	MidLabel               string  `json:"mid_label"`
	Currency               string  `json:"currency"`
	SuccessRedirectURL     string  `json:"success_redirect_url"`
	FailureRedirectURL     string  `json:"failure_redirect_url"`
	PaidAt                 string  `json:"paid_at"`
	CreditCardChargeID     string  `json:"credit_card_charge_id"`
	PaymentMethod          string  `json:"payment_method"`
	PaymentChannel         string  `json:"payment_channel"`
	PaymentDestination     string  `json:"payment_destination"`
}

func (x *Client) InvoiceNotify(notify InvoiceNotify, token string) (bool, error) {
	if x.VerifyToken != token {
		return false, nil
	}

	return true, nil
}

type PayoutRequest struct {
	ReferenceID         string              `json:"reference_id"`
	ChannelCode         string              `json:"channel_code"`
	ChannelProperties   ChannelProperties   `json:"channel_properties"`
	Amount              float64             `json:"amount"`
	Description         string              `json:"description"`
	Currency            string              `json:"currency"`
	ReceiptNotification ReceiptNotification `json:"receipt_notification"`
}
type ReceiptNotification struct {
	EmailTo []string `json:"email_to"`
	EmailCc []string `json:"email_cc"`
}

type ChannelProperties struct {
	AccountHolderName string `json:"account_holder_name"`
	AccountNumber     string `json:"account_number"`
}

type PayoutResponse struct {
	ID                   string              `json:"id"`
	Amount               float64             `json:"amount"`
	ChannelCode          string              `json:"channel_code"`
	Currency             string              `json:"currency"`
	Description          string              `json:"description"`
	ReferenceID          string              `json:"reference_id"`
	Status               string              `json:"status"`
	Created              time.Time           `json:"created"`
	Updated              time.Time           `json:"updated"`
	EstimatedArrivalTime time.Time           `json:"estimated_arrival_time"`
	BusinessID           string              `json:"business_id"`
	ChannelProperties    ChannelProperties   `json:"channel_properties"`
	ReceiptNotification  ReceiptNotification `json:"receipt_notification"`
	ErrorCode            string              `json:"error_code"`
	Message              string              `json:"message"`
}

func (x *Client) PayOut(payout PayoutRequest) (*PayoutResponse, error) {
	url := "https://api.xendit.co/v2/payouts"
	requestId := utils.GenerateRequestId()

	header := map[string]string{
		"Authorization":   fmt.Sprintf("Basic %s", utils.EncodeStr2Base64(x.SecretKey)),
		"Idempotency-key": requestId,
	}
	result, err := httpHelper.Post(url, payout, header)
	if err != nil {
		return nil, err
	}

	var response PayoutResponse
	err = utils.Json2Struct(result, &response)

	return &response, err
}

type BalanceResponse struct {
	Balance float64 `json:"balance"`
}

func (x *Client) GetBalance() (*BalanceResponse, error) {
	url := "https://api.xendit.co/balance"
	header := map[string]string{
		"Authorization": fmt.Sprintf("Basic %s", utils.EncodeStr2Base64(x.SecretKey)),
	}
	result, err := httpHelper.GetHeader(url, header)
	if err != nil {
		return nil, err
	}

	var res BalanceResponse
	err = utils.Json2Struct(result, &res)

	return &res, err
}

type InvoiceCallback struct {
	ID                     string  `json:"id"`
	ExternalID             string  `json:"external_id"`
	UserID                 string  `json:"user_id"`
	Status                 string  `json:"status"`
	MerchantName           string  `json:"merchant_name"`
	Amount                 float64 `json:"amount"`
	PayerEmail             string  `json:"payer_email"`
	Description            string  `json:"description"`
	FeesPaidAmount         float64 `json:"fees_paid_amount"`
	AdjustedReceivedAmount float64 `json:"adjusted_received_amount"`
	BankCode               string  `json:"bank_code"`
	RetailOutletName       string  `json:"retail_outlet_name"`
	EwalletType            string  `json:"ewallet_type"`
	OnDemandLink           string  `json:"on_demand_link"`
	RecurringPaymentID     string  `json:"recurring_payment_id"`
	PaidAmount             float64 `json:"paid_amount"`
	Updated                string  `json:"updated"`
	Created                string  `json:"created"`
	MidLabel               string  `json:"mid_label"`
	Currency               string  `json:"currency"`
	SuccessRedirectURL     string  `json:"success_redirect_url"`
	FailureRedirectURL     string  `json:"failure_redirect_url"`
	PaidAt                 string  `json:"paid_at"`
	CreditCardChargeID     string  `json:"credit_card_charge_id"`
	PaymentMethod          string  `json:"payment_method"`
	PaymentChannel         string  `json:"payment_channel"`
	PaymentDestination     string  `json:"payment_destination"`
}

type PayOutNotify struct {
	Event      string         `json:"event"`
	BusinessId string         `json:"business_id"`
	Created    string         `json:"created"`
	Data       PayoutResponse `json:"data"`
}
