package rocketmq

import (
	"context"
	"fmt"
	v2 "github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/deng00/go-base/mq"
)

const AliyunChannel = "ALIYUN"

// Config rocket v2 client config
type Config struct {
	GroupID   string `mapstructure:"group_id"`
	Addr      string `mapstructure:"addr"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Channel   string `mapstructure:"channel"`
	Topic     string `mapstructure:"topic"`
	Broadcast bool   `mapstructure:"broadcast"`
	Namespace string `mapstructure:"namespace"`
}

func (c *Config) Check() error {
	if c.Addr == "" {
		return fmt.Errorf("invalid addr")
	}
	if c.AccessKey != "" && c.Channel == "" {
		c.Channel = AliyunChannel
	}
	if c.AccessKey != "" && c.Namespace == "" {
		return fmt.Errorf("invalid namespace")
	}
	return nil
}

// rocket mq client
var _ mq.Interface = &Client{}

type Client struct {
	config       *Config
	pushConsumer v2.PushConsumer
	producer     v2.Producer
}

// Listen register listener
func (r *Client) Listen(topic, expression string, listener func(data mq.Message) error) (err error) {
	return r.pushConsumer.Subscribe(topic, consumer.MessageSelector{Type: consumer.TAG, Expression: expression},
		func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			for _, msg := range msgs {
				if err := listener(mq.Message{
					Topic: msg.Topic,
					Tags:  msg.GetTags(),
					Keys:  msg.GetKeys(),
					Body:  string(msg.Body),
				}); err != nil {
					return consumer.ConsumeRetryLater, err
				}
			}
			return consumer.ConsumeSuccess, nil
		})
}

// Publish publish message
func (r *Client) Publish(msg mq.Message) (err error) {
	_, err = r.producer.SendSync(context.Background(),
		primitive.NewMessage(msg.Topic, []byte(msg.Body)).WithTag(msg.Tags).WithKeys([]string{msg.Keys}))
	return
}

// New create new rocket mq client
func New(config *Config) (client *Client, err error) {
	pd, err := v2.NewProducer(
		producer.WithGroupName(config.GroupID),
		producer.WithNameServer([]string{config.Addr}),
		producer.WithRetry(3),
		producer.WithCredentials(primitive.Credentials{
			AccessKey: config.AccessKey,
			SecretKey: config.SecretKey,
		}),
		producer.WithNamespace(config.Namespace),
	)
	if err != nil {
		return nil, fmt.Errorf("create new pd err:%s", err)
	}
	if err := pd.Start(); err != nil {
		return nil, err
	}
	model := consumer.Clustering
	if config.Broadcast {
		model = consumer.BroadCasting
	}
	cs, err := v2.NewPushConsumer(
		consumer.WithGroupName(config.GroupID),
		consumer.WithNameServer([]string{config.Addr}),
		consumer.WithCredentials(primitive.Credentials{
			AccessKey: config.AccessKey,
			SecretKey: config.SecretKey,
		}),
		consumer.WithConsumerModel(model),
		consumer.WithNamespace(config.Namespace),
		consumer.WithConsumeMessageBatchMaxSize(1), //
	)
	if err != nil {
		return nil, fmt.Errorf("create new rocketmq client err:%s\n", err)
	}
	client = &Client{
		config:       config,
		producer:     pd,
		pushConsumer: cs,
	}
	return
}

// Close close rocket mq client
func (r *Client) Close() {
	_ = r.pushConsumer.Shutdown()
	_ = r.producer.Shutdown()
}

// StartListen start client
func (r *Client) StartListen() (err error) {
	return r.pushConsumer.Start()
}
