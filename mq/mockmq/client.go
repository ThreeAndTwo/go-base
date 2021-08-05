package mockmq

import (
	"github.com/deng00/go-base/mq"
	"strings"
	"time"
)

var _ mq.Interface = &Mock{}

type Mock struct {
	msgChan chan mq.Message
	exitCh  chan struct{}
}

func (m *Mock) StartListen() (err error) {
	return nil
}

func (m *Mock) Listen(topic, expression string, listener func(msg mq.Message) error) error {
	go func() {
		for {
			select {
			case <-m.exitCh:
				return
			case msg := <-m.msgChan:
				if msg.Topic != topic || strings.Contains(expression, msg.Tags) {
					break
				}
				err := listener(msg)
				if err != nil {
					time.AfterFunc(1*time.Second, func() {
						m.msgChan <- msg
					})
				}
			}
		}
	}()
	return nil
}

func (m *Mock) Publish(msg mq.Message) error {
	m.msgChan <- msg
	return nil
}

func (m *Mock) Start() (err error) {
	return nil
}

func (m *Mock) Close() {
}

func New() *Mock {
	return &Mock{
		msgChan: make(chan mq.Message, 100),
		exitCh:  make(chan struct{}),
	}
}
