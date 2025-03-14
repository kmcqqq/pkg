package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	DingDing = iota
	FeiShu
	WeChat
)

// 定义消息结构体
type Message struct {
	Platform uint
	Content  string
}

// 钉钉消息格式
type DingDingMessage struct {
	MsgType string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}

// 飞书消息格式
type FeiShuMessage struct {
	Header struct {
		ContentType string `json:"content-type"`
	} `json:"header"`
	MsgType string `json:"msg_type"`
	Content struct {
		Text string `json:"text"`
	} `json:"content"`
}

// 企业微信消息格式
type WeChatMessage struct {
	MsgType string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}

// 通用发送 Webhook 消息函数
func SendWebhook(message Message, webhookURL string) error {
	var reqBody []byte
	var err error

	// 根据平台构造不同的消息格式
	switch message.Platform {
	case DingDing:
		dingMessage := DingDingMessage{
			MsgType: "text",
		}
		dingMessage.Text.Content = message.Content
		reqBody, err = json.Marshal(dingMessage)
		break
	case FeiShu:
		feiShuMessage := FeiShuMessage{
			MsgType: "text",
		}
		feiShuMessage.Content.Text = message.Content
		reqBody, err = json.Marshal(feiShuMessage)
		break
	case WeChat:
		wechatMessage := WeChatMessage{
			MsgType: "text",
		}
		wechatMessage.Text.Content = message.Content
		reqBody, err = json.Marshal(wechatMessage)
		break
	default:
		return fmt.Errorf("unsupported platform: %d", message.Platform)
	}

	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// 发送 POST 请求
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message, status code: %d", resp.StatusCode)
	}

	//log.Println("Message sent successfully!")
	return nil
}
