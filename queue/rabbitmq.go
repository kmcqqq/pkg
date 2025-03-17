package queue

import (
	"context"
	"fmt"
	"github.com/kmcqqq/pkg/config"
	"github.com/kmcqqq/pkg/logger"
	"github.com/streadway/amqp"
	"time"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func (r *RabbitMQ) PublishDelayMessage(routingKey, message string, delayTime time.Duration) error {
	return r.channel.Publish(
		"delay-exchange", //exchangeName
		routingKey,       //routing key
		true,             //mandatory
		false,            //immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(message),
			Headers: amqp.Table{
				"x-delay": int(delayTime / time.Millisecond),
			},
		},
	)
}

var _ Queue = &RabbitMQ{}

func (r *RabbitMQ) Init(cfg *config.ServerInfo) error {
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.User, cfg.Pwd, cfg.Host, cfg.Port)
	var err error
	r.conn, err = amqp.Dial(dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	r.channel, err = r.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %w", err)
	}

	return nil
}

func (r *RabbitMQ) PublishMessageByExchange(exchangeName, routingKey, message string) error {
	return r.channel.Publish(
		exchangeName, //exchangeName
		routingKey,   //routing key
		true,         //mandatory
		false,        //immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}

func (r *RabbitMQ) ConsumeMessages(queueName string) (<-chan amqp.Delivery, error) {
	//	实现 rabbitmq 消费
	return r.channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
}

func (r *RabbitMQ) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}

type RabbitMQConsumer struct {
	client  *RabbitMQ
	queue   string
	handler Handler
}

func (r *RabbitMQ) GetConsumer(queue string, handler Handler) *RabbitMQConsumer {
	return &RabbitMQConsumer{
		client:  r,
		queue:   queue,
		handler: handler,
	}
}

func (c *RabbitMQConsumer) Start(ctx context.Context) error {
	msgs, err := c.client.channel.Consume(
		c.queue,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %w", err)
	}

	go func() {
		for msg := range msgs {
			//logger.Debug("consumer", logger.String("topic", c.queue), logger.String("key", msg.RoutingKey), logger.String("data", string(msg.Body)))

			message := &Message{
				Topic: c.queue,
				Key:   msg.RoutingKey,
				Data:  string(msg.Body),
			}
			if err := c.handler(ctx, message); err != nil {
				logger.Error("error", logger.String("title", "consumer error"), logger.String("topic", c.queue), logger.String("key", msg.RoutingKey), logger.String("data", string(msg.Body)), logger.Err(err))
			}
		}
	}()

	return nil
}
