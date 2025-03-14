package queue

import (
	"context"
	"github.com/streadway/amqp"
	"gitlab.bobbylive.cn/kongmengcheng/pkg/config"
	"time"
)

type Queue interface {
	Init(cfg *config.ServerInfo) error                                             // 初始化队列连接和通道
	PublishMessageByExchange(exchangeName, routingKey, message string) error       // 发送消息
	PublishDelayMessage(routingKey, message string, delayTime time.Duration) error // 发送消息
	ConsumeMessages(queueName string) (<-chan amqp.Delivery, error)                // 消费消息
	Close()                                                                        // 关闭连接和通道
	GetConsumer(queue string, handler Handler) *RabbitMQConsumer
}

type Message struct {
	Topic string
	Key   string
	Data  string
}

type Handler func(context.Context, *Message) error
