package rocketmq

import (
	"fmt"
	"github.com/deng00/go-base/mq"
	"testing"
	"time"
)

var client *Client

func init() {
	var err error
	client, err = New(&Config{
		GroupID: "test_group",
		Addr:    "127.0.0.1:9876",
	})
	if err != nil {
		panic(err)
	}
}
func TestMain(m *testing.M) {
	m.Run()
	client.Close()
}

func TestRocketMQ_ListenAndPublish(t *testing.T) {
	type args struct {
		topic      string
		expression string
		listener   func(data mq.Message) error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test_listenandpublish",
			args: args{
				topic:      "demo",
				expression: "*",
				listener: func(data mq.Message) error {
					t.Log(data)
					return nil
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := client.Listen(tt.args.topic, tt.args.expression, tt.args.listener); err != nil {
				t.Fatalf("listen failed:%s", err)
			}
			if err := client.StartListen(); err != nil {
				panic(fmt.Sprintf("mq start failed:%s", err))
				return
			}
			if err := client.Publish(mq.Message{
				Topic: tt.args.topic,
				Tags:  "test_topic",
				Keys:  "test_key",
				Body:  "hello world",
			}); err != nil {
				t.Fatalf("publish msg failed:%s", err)
			}
			time.Sleep(1 * time.Second)
		})
	}
}
