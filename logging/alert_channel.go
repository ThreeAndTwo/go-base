package logging

import (
	"encoding/json"
	"fmt"
	"github.com/royeo/dingrobot"
	"os"
)

type AlertChannelI interface {
	SetToken(token string)
	SetService(service string)
	Send(msg string) error
}

var DingDingApi = "https://oapi.dingtalk.com/robot/send?access_token="

type AlertDataField struct {
	Level    string `json:"level"`
	Time     string `json:"time"`
	Line     string `json:"line"`
	Msg      string `json:"msg"`
	Service  string `json:"service"`
	Hostname string `json:"hostname"`
}

func (order *AlertDataField) FromJSON(msg string) error {
	return json.Unmarshal([]byte(msg), order)
}

func (order *AlertDataField) ToJSON() string {
	str, _ := json.Marshal(order)
	return string(str)
}

type AlertChanel struct {
	Token   string
	service string
}

func (c *AlertChanel) SetToken(token string) {
	c.Token = token
}

func (c *AlertChanel) SetService(service string) {
	c.service = service
}

func (c *AlertChanel) Send(msg string) error {
	return nil
}

type DingDingAlertChanel struct {
	AlertChanel
}

func (c *DingDingAlertChanel) Send(msg string) error {
	fields := &AlertDataField{}
	_ = fields.FromJSON(msg)
	hostname, _ := os.Hostname()
	fields.Hostname = hostname
	fields.Service = c.service
	robot := dingrobot.NewRobot(DingDingApi + c.Token)
	text := fmt.Sprintf("* level=%s\n* time=%s\n* line=%s\n* msg=%s\n* service=%s\n* hostname=%s\n", fields.Level, fields.Time, fields.Line, fields.Msg, fields.Service, fields.Hostname)
	md := "### " + c.service + "服务告警 \n " + text
	var atMobiles []string
	err := robot.SendMarkdown(c.service+" Alert", md, atMobiles, false)
	if err != nil {
		println("send dingding alert failed: " + err.Error())
	}
	return err
}

func NewDingDingAlertChanel(token string) *DingDingAlertChanel {
	channel := &DingDingAlertChanel{}
	channel.Token = token
	return channel
}
