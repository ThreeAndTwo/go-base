package rocketmq

import (
	"fmt"
	rocketmq "github.com/apache/rocketmq-client-go/core"
	"github.com/deng00/go-base/mq"
)

const AliyunChannel = "ALIYUN"

// Config rocket client config
type Config struct {
	GroupID   string `mapstructure:"group_id"`
	Addr      string `mapstructure:"addr"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Channel   string `mapstructure:"channel"`
	Topic     string `mapstructure:"topic"`
	Broadcast bool   `mapstructure:"broadcast"`
}

func (c *Config) Check() error {
	if c.Addr == "" {
		return fmt.Errorf("invalid addr")
	}
	if c.AccessKey != "" && c.Channel == "" {
		c.Channel = AliyunChannel
	}
	return nil
}

// rocket mq client
var _ mq.Interface = &Client{}

type Client struct {
	config       *Config
	pushConsumer rocketmq.PushConsumer
	producer     rocketmq.Producer
}

// Listen register listener
func (r *Client) Listen(topic, expression string, listener func(data mq.Message) error) (err error) {
	return r.pushConsumer.Subscribe(topic, expression, func(msg *rocketmq.MessageExt) rocketmq.ConsumeStatus {
		err := listener(mq.Message{
			Topic: msg.Topic,
			Tags:  msg.Tags,
			Keys:  msg.Keys,
			Body:  msg.Body,
		})
		if err != nil {
			return rocketmq.ReConsumeLater
		}
		return rocketmq.ConsumeSuccess
	})
}

// Publish publish message
func (r *Client) Publish(msg mq.Message) (err error) {
	_, err = r.producer.SendMessageSync(&rocketmq.Message{
		Topic: msg.Topic,
		Tags:  msg.Tags,
		Keys:  msg.Keys,
		Body:  msg.Body,
	})
	return
}

// New create new rocket mq client
func New(config *Config) (client *Client, err error) {
	clientConfig := rocketmq.ClientConfig{
		GroupID:    config.GroupID,
		NameServer: config.Addr,
		Credentials: &rocketmq.SessionCredentials{
			AccessKey: config.AccessKey,
			SecretKey: config.SecretKey,
			Channel:   config.Channel,
		},
	}
	producer, err := rocketmq.NewProducer(&rocketmq.ProducerConfig{
		ClientConfig:   clientConfig,
		ProducerModel:  rocketmq.CommonProducer,
		MaxMessageSize: 4 * 1024 * 1024, // max message size 4 MB
	})
	if err != nil {
		return nil, fmt.Errorf("create new producer err:%s", err)
	}
	if err := producer.Start(); err != nil {
		return nil, err
	}
	model := rocketmq.Clustering
	if config.Broadcast {
		model = rocketmq.BroadCasting
	}
	consumer, err := rocketmq.NewPushConsumer(&rocketmq.PushConsumerConfig{
		ClientConfig:  clientConfig,
		Model:         model,
		ConsumerModel: rocketmq.CoCurrently,
	})
	if err != nil {
		return nil, fmt.Errorf("create new rocketmq client err:%s\n", err)
	}
	client = &Client{
		config:       config,
		producer:     producer,
		pushConsumer: consumer,
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
