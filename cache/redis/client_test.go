package redis

import (
	"fmt"
	"testing"
	"time"
)

var redisClient *Redis

func init() {
	var err error
	redisClient, err = New(&Config{
		Addr:          "127.0.0.1:6379",
		Pass:          "",
		TlsSkipVerify: false,
	})
	if err != nil {
		panic(err)
	}
}
func TestRedis_PubSub(t *testing.T) {
	type args struct {
		channel    string
		subscriber func(event Event)
		event      Event
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test_subcribe",
			args: args{
				channel: "testchan",
				event: Event{
					Channel: "testchan",
					Payload: "hello world",
				},
				subscriber: func(event Event) {
					fmt.Println(event)
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redisClient.Subscribe(tt.args.subscriber, tt.args.channel)
			_, err := redisClient.Publish(tt.args.event)
			if err != nil {
				t.Error(err)
				return
			}
			time.Sleep(2 * time.Second)
		})
	}
}

func TestRedis_SetAndGet(t *testing.T) {
	type args struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test_redis_setandget01",
			args: args{
				key:   "name",
				value: "xiaoming",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := redisClient.Set(tt.args.key, tt.args.value, time.Minute); err != nil {
				t.Fatalf("set failed:%s", err)
			}
			got, err := redisClient.Get(tt.args.key)
			if err != nil {
				t.Fatalf("get failed:%s", err)
			}
			if got != tt.args.value {
				t.Fatalf("got value(%s) not equal set value(%s)", got, tt.args.value)
			}
		})
	}
}
