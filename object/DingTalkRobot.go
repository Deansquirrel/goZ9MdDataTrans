package object

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Deansquirrel/goToolCommon"
	log "github.com/Deansquirrel/goToolLog"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	//TODO 钉钉消息发送服务地址待参数化
	address = "http://123.57.70.114:10001"
)

type dingTalkRobot struct {
	config *DingTalkRobotConfigData
}

func NewDingTalkRobot(config *DingTalkRobotConfigData) *dingTalkRobot {
	return &dingTalkRobot{
		config: config,
	}
}

type dingTalkTextMsg struct {
	WebHookKey string   `json:"webhookkey"`
	Content    string   `json:"content"`
	AtMobiles  []string `json:"atmobiles"`
	IsAtAll    bool     `json:"isatall"`
}

func (d *dingTalkRobot) SendMsg(msg string) error {
	log.Debug(fmt.Sprintf("DingTalkRobot send msg [%s]", msg))
	if d.config == nil {
		return errors.New(fmt.Sprintf("dingTalkRobot config data is nil"))
	}
	msg = d.msgFormat(msg)
	if d.config.FIsAtAll == 1 {
		return d.sendTextMsgWithAtAll(d.config.FWebHookKey, msg)
	}
	if strings.Trim(d.config.FAtMobiles, " ") != "" {
		list := strings.Split(strings.Trim(d.config.FAtMobiles, " "), ",")
		list = goToolCommon.ClearBlock(list)
		log.Debug(strconv.Itoa(len(list)))
		if len(list) > 0 {
			log.Debug(d.config.FWebHookKey)
			log.Debug(msg)
			return d.sendTextMsgWithAtList(d.config.FWebHookKey, msg, list)
		}
	}
	return d.sendTextMsg(d.config.FWebHookKey, msg)
}

func (d *dingTalkRobot) msgFormat(msg string) string {
	return goToolCommon.GetDateTimeStr(time.Now()) + "\n" + msg
}

func (d *dingTalkRobot) sendTextMsg(webHookKey string, msg string) error {
	om := dingTalkTextMsg{
		WebHookKey: webHookKey,
		Content:    msg,
	}
	return d.sendMsg(om)
}

func (d *dingTalkRobot) sendTextMsgWithAtList(webHookKey string, msg string, atMobiles []string) error {
	om := dingTalkTextMsg{
		WebHookKey: webHookKey,
		Content:    msg,
		AtMobiles:  atMobiles,
	}
	return d.sendMsg(om)
}

func (d *dingTalkRobot) sendTextMsgWithAtAll(webHookKey string, msg string) error {
	om := dingTalkTextMsg{
		WebHookKey: webHookKey,
		Content:    msg,
		IsAtAll:    true,
	}
	return d.sendMsg(om)
}

//获取Text消息发送地址
func (d *dingTalkRobot) getTextMsgUrl() string {
	return fmt.Sprintf("%s/text", address)
}

//发送消息
func (d *dingTalkRobot) sendMsg(v interface{}) error {
	msg, err := goToolCommon.GetJsonStr(v)
	if err != nil {
		return err
	}
	rData, err := d.sendData([]byte(msg), d.getTextMsgUrl())
	if err != nil {
		return err
	}
	return d.tranRData(rData)
}

//POST发送数据
func (d *dingTalkRobot) sendData(data []byte, url string) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		log.Error(err.Error())
		return nil, errors.New("构造http请求数据时发生错误：" + err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return nil, errors.New("发送http请求时错误：" + err.Error())
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	rData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err.Error())
		return nil, errors.New("读取http返回数据时发生错误：" + err.Error())
	}
	return rData, nil
}

type dingTalkRobotResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

//检查返回数据
func (d *dingTalkRobot) tranRData(data []byte) error {
	var r dingTalkRobotResponse
	err := json.Unmarshal(data, &r)
	if err != nil {
		return errors.New(fmt.Sprintf("返回数据格式化异常，err：[%s]，返回数据：[%s]", err.Error(), string(data)))
	}
	if r.ErrCode == 0 && strings.ToLower(r.ErrMsg) == "ok" {
		return nil
	} else {
		if strings.Trim(r.ErrMsg, " ") != "" {
			return errors.New(r.ErrMsg)
		} else {
			return errors.New(fmt.Sprintf("未知错误，errcode[%d]", r.ErrCode))
		}
	}
}
