package mq

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

type Message struct {
	Topic string
	Tags  string
	Keys  string
	Body  string
}

func (m Message) String() string {
	rawStr := `
	topic:%s
	tags:%s
	keys:%s
	body:%s
`
	return fmt.Sprintf(rawStr, m.Topic, m.Tags, m.Keys, m.Body)
}
func (m Message) Hash() string {
	sumBytes := md5.Sum([]byte(m.Topic + m.Tags + m.Keys + m.Body))
	return hex.EncodeToString(sumBytes[:])
}

type Interface interface {
	Listen(topic, expression string, listener func(msg Message) error) error
	Publish(msg Message) error
	StartListen() (err error)
	Close()
}
