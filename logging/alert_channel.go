package logging

import (
	"encoding/json"
	"fmt"
	"github.com/bluele/slack"
	"github.com/royeo/dingrobot"
	"os"
	"strings"
)

type AlertChannelI interface {
	SetToken(token string)
	SetServiceName(service string)
	SetMsg(msg string)
	Send() error
}

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
	token       string
	serviceName string
	channelName string
	mdMsg       string
	fields      *AlertDataField
}

func (c *AlertChanel) SetToken(token string) {
	c.token = token
}

func (c *AlertChanel) SetServiceName(serviceName string) {
	c.serviceName = serviceName
}

func (c *AlertChanel) SetChannelName(channelName string) {
	c.channelName = channelName
}

func (c *AlertChanel) SetMsg(msg string) {
	c.fields = &AlertDataField{}
	_ = c.fields.FromJSON(msg)
	hostname, _ := os.Hostname()
	c.fields.Hostname = hostname
	text := fmt.Sprintf("* level=%s\n* time=%s\n* line=%s\n* msg=%s\n* serviceName=%s\n* hostname=%s\n", c.fields.Level, c.fields.Time, c.fields.Line, c.fields.Msg, c.fields.Service, c.fields.Hostname)
	c.mdMsg = "### " + c.serviceName + "Service Alert \n " + text
}

func (c *AlertChanel) Send() error {
	return nil
}

type DingDingAlertChanel struct {
	AlertChanel
}

func (c *DingDingAlertChanel) Send() error {
	robot := dingrobot.NewRobot("https://oapi.dingtalk.com/robot/send?access_token=" + c.token)
	var atMobiles []string
	err := robot.SendMarkdown(c.serviceName+" Alert", c.mdMsg, atMobiles, false)
	if err != nil {
		println("send dingding alert failed: " + err.Error())
	}
	return err
}

type SlackAlertChanel struct {
	AlertChanel
}

func (c *SlackAlertChanel) Send() error {
	api := slack.New(c.token)
	err := api.ChatPostMessage(c.channelName, c.mdMsg, nil)
	if err != nil {
		println("send slack alert failed: " + err.Error())
	}
	return err
}

func NewDingDingAlertChanel(token string) *DingDingAlertChanel {
	channel := &DingDingAlertChanel{}
	channel.SetToken(token)
	return channel
}

func NewSlackAlertChanel(token string) *SlackAlertChanel {
	params := strings.Split(token, "@")
	if len(params) != 2 {
		return nil
	}
	channel := &SlackAlertChanel{}
	channel.SetToken(params[0])
	channel.SetChannelName(params[1])
	return channel
}
