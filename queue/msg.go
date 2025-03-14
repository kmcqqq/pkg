package queue

import (
	"gitlab.bobbylive.cn/kongmengcheng/pkg/utils"
	"time"
)

type NotifyMsg struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

type MessageService struct {
	queue Queue // 通过接口注入队列
}

// 创建 MessageService 实例，注入队列实例
func NewMessageService(q Queue) *MessageService {
	return &MessageService{queue: q}
}

// Direct 模式
func (r *MessageService) Direct(routingKey string, body string) error {
	err := r.queue.PublishMessageByExchange("amq.direct", routingKey, body)
	return err
}

// PushDelayMsg 延迟队列推送
func (r *MessageService) PushDelayMsg(routingKey string, delayTime time.Duration, data string) error {
	err := r.queue.PublishDelayMessage(routingKey, data, delayTime)
	return err
}

// RefreshCoin 币更新
func (r *MessageService) RefreshCoin(userIdx int64, cash int64) error {
	message := NotifyMsg{
		Code: 101,
		Data: struct {
			UserIdx int64 `json:"useridx"`
			Cash    int64 `json:"cash"`
		}{
			UserIdx: userIdx,
			Cash:    cash,
		},
	}
	messageStr, err := utils.Struct2Json(message)
	if err != nil {
		return err
	}
	err = r.queue.PublishMessageByExchange("notify", "", messageStr)
	return err
}

// RefreshGralcash 果子更新通知
func (r *MessageService) RefreshGralcash(userIdx int64, gralcash float64) error {
	message := NotifyMsg{
		Code: 112,
		Data: struct {
			UserIdx int64   `json:"useridx"`
			Cash    float64 `json:"cash"`
		}{
			UserIdx: userIdx,
			Cash:    gralcash,
		},
	}
	messageStr, err := utils.Struct2Json(message)
	if err != nil {
		return err
	}
	err = r.queue.PublishMessageByExchange("notify", "", messageStr)
	return err
}

// RechargeSuccess 充值成功
func (r *MessageService) RechargeSuccess(userIdx int64, cash int64, money float64, dtype int, productId int, couponId int, currency string, isFirstPay bool, orderId string) error {
	message := NotifyMsg{
		Code: 101,
		Data: struct {
			UserIdx    int64   `json:"useridx"`
			Cash       int64   `json:"cash"`
			Money      float64 `json:"money"`
			Dtype      int     `json:"dtype"`
			ProductId  int     `json:"productId"`
			CouponId   int     `json:"couponId"`
			Currency   string  `json:"currency"`
			IsFirstPay bool    `json:"isFirstPay"`
			OrderId    string  `json:"orderId"`
		}{
			UserIdx:    userIdx,
			Cash:       cash,
			Money:      money,
			Dtype:      dtype,
			ProductId:  productId,
			CouponId:   couponId,
			Currency:   currency,
			IsFirstPay: isFirstPay,
			OrderId:    orderId,
		},
	}
	messageStr, err := utils.Struct2Json(message)
	if err != nil {
		return err
	}
	err = r.queue.PublishMessageByExchange("notify", "", messageStr)
	return err
}

// GameWinFloating 游戏中奖飘条
func (r *MessageService) GameWinFloating(userIdx int64, cash int64, roomId int64, gameId, ntype int, gameIcon string) error {
	message := NotifyMsg{
		Code: 161,
		Data: struct {
			UserIdx  int64  `json:"useridx"`
			Cash     int64  `json:"cash"`
			RoomId   int64  `json:"roomid"`
			GameId   int    `json:"gameId"`
			Type     int    `json:"type"`
			GameIcon string `json:"gameIcon"`
		}{
			UserIdx:  userIdx,
			Cash:     cash,
			RoomId:   roomId,
			GameId:   gameId,
			Type:     ntype,
			GameIcon: gameIcon,
		},
	}
	messageStr, err := utils.Struct2Json(message)
	if err != nil {
		return err
	}
	err = r.queue.PublishMessageByExchange("notify", "", messageStr)
	return err
}

// UpdateBagInfo 更新背包道具 goodsType 2 坐骑 3vip 4 头像框 18
func (r *MessageService) UpdateBagInfo(userIdx int64, goodsType int, param string) error {
	message := NotifyMsg{
		Code: 121,
		Data: struct {
			UserIdx   int64  `json:"useridx"`
			Goodstype int    `json:"goodstype"`
			Param     string `json:"param"`
		}{
			UserIdx:   userIdx,
			Goodstype: goodsType,
			Param:     param,
		},
	}

	messageStr, err := utils.Struct2Json(message)
	if err != nil {
		return err
	}

	err = r.queue.PublishMessageByExchange("notify", "", messageStr)
	return err
}

// SendVipPiaoTiao 开通 VIP 票条
func (r *MessageService) SendVipPiaoTiao(userIdx int64, level int, content string) error {
	message := NotifyMsg{
		Code: 129,
		Data: struct {
			UserIdx int64  `json:"useridx"`
			Level   int    `json:"level"`
			Content string `json:"content"`
		}{
			UserIdx: userIdx,
			Level:   level,
			Content: content,
		}}

	messageStr, err := utils.Struct2Json(message)
	if err != nil {
		return err
	}

	err = r.queue.PublishMessageByExchange("notify", "", messageStr)
	return err
}
