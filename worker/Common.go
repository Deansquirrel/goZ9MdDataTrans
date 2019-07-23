package worker

import (
	"errors"
	"fmt"
	"github.com/Deansquirrel/goToolCommon"
	"github.com/Deansquirrel/goToolCron"
	"github.com/Deansquirrel/goZ9MdDataTrans/object"
	"github.com/Deansquirrel/goZ9MdDataTrans/repository"
	"runtime/debug"
	"strings"
)

import log "github.com/Deansquirrel/goToolLog"

const (
	//TODO 钉钉机器人消息Key待参数化
	webHook = "7685d3a4b1630c8cb0028ee9cd17beae1578869f8ddf0b7871ecbb8b0743f8ce"
)

type common struct {
}

func NewCommon() *common {
	return &common{}
}

//启动工作线程
func (c *common) StartWorker(key object.TaskKey) {
	var err error
	var cronStr string
	defer func() {
		if err != nil {
			c.HandleErr(key, err)
		} else {
			log.Debug(fmt.Sprintf("start worker %s cron %s", key, cronStr))
		}
	}()
	repOnline, err := repository.NewRepOnline()
	if err != nil {
		return
	}
	cronStr, err = repOnline.GetTaskCron(key)
	if err != nil {
		return
	}
	cronStr = strings.Trim(cronStr, " ")
	if cronStr == "" {
		err = errors.New(fmt.Sprintf("task %s start err: cron is empty", key))
		return
	}

	err = goToolCron.AddFunc(string(key), cronStr, c.getWorkerFuncReal(key), c.GetHandlePanic(key))
}

//停止工作线程
func (c *common) StopWorker(key object.TaskKey) {
	defer log.Debug(fmt.Sprintf("stop worker %s", key))
	goToolCron.Stop(string(key))
}

func (c *common) getWorkerFuncReal(key object.TaskKey) func() {
	return func() {
		guid := goToolCommon.Guid()

		log.Debug(fmt.Sprintf("task %s[%s] start", key, guid))
		defer log.Debug(fmt.Sprintf("task %s[%s] complete", key, guid))

		f := c.getWorkerFunc(key)
		if f != nil {
			f()
		} else {
			errMsg := fmt.Sprintf("task %s[%s] start err: func is nil", key, guid)
			c.HandleErr(key, errors.New(errMsg))
			return
		}
	}
}

//获取任务执行函数
func (c *common) getWorkerFunc(key object.TaskKey) func() {
	switch key {
	case object.TaskKeyRefreshConfig:
		return NewCommonWorker().RefreshConfig
	case object.TaskKeyRefreshHeartBeat:
		return NewOnlineWorker().RefreshHeartBeat
	case object.TaskKeyRefreshMdYyInfo:
		return NewOnlineWorker().UpdateMdYyInfo
	case object.TaskKeyRefreshZxKc:
		return NewOnlineWorker().UpdateZxKc
	case object.TaskKeyRestoreMdYyInfo:
		return NewBbWorker().RestoreMdYyInfo
	case object.TaskKeyRestoreZxKc:
		return NewBbWorker().RestoreZxKc
	default:
		return nil
	}
}

//	switch key {
//	case object.TaskKeyHeartBeat:
//		c.startHeartBeat(errHandle)
//	case object.TaskKeyRefreshMdDataTransState:
//		c.startRefreshMdDataTransState(errHandle)
//	case object.TaskKeyRestoreMdYyStateTransTime:
//		c.startRestoreMdYyStateTransTime(errHandle)
//	case object.TaskKeyRefreshWaitRestoreDataCount:
//		c.startRefreshWaitRestoreDataCount(errHandle)
//	case object.TaskKeyRestoreMdYyStateRestoreTime:
//		c.startRestoreMdYyStateRestoreTime(errHandle)
//	case object.TaskKeyRestoreMdYyState:
//		c.startRestoreMdYyState(errHandle)
//	case object.TaskKeyRestoreMdSet:
//		c.startRestoreMdSet(errHandle)
//	case object.TaskKeyRestoreCwGsSet:
//		c.startRestoreCwGsSet(errHandle)
//	case object.TaskKeyRestoreMdCwGsRef:
//		c.startRestoreMdCwGsRef(errHandle)
//	case object.TaskKeyRefreshConfig:
//		c.StartRefreshConfig(errHandle)
//	default:
//		errMsg := fmt.Sprintf("unknow task key: %s", key)
//		log.Error(errMsg)
//		c.errChan <- errors.New(errMsg)
//	}

//task错误处理
func (c *common) HandleErr(key object.TaskKey, err error) {
	if err != nil {
		log.Error(err.Error())
		c.sendMsg(err.Error())
	}
}

//发送钉钉消息
func (c *common) sendMsg(msg string) {
	log.Warn(msg)
	//dt := object.NewDingTalkRobot(&object.DingTalkRobotConfigData{
	//	FWebHookKey: webHook,
	//	FAtMobiles:  "",
	//	FIsAtAll:    0,
	//})
	//sendErr := dt.SendMsg(msg)
	//if sendErr != nil {
	//	log.Error(fmt.Sprintf("send msg error,msg: %s, error: %s", msg, sendErr.Error()))
	//}
}

func (c *common) GetHandlePanic(key object.TaskKey) func(interface{}) {
	return func(err interface{}) {
		c.handlePanic(key, err)
	}
}

func (c *common) handlePanic(key object.TaskKey, err interface{}) {
	errMsg := fmt.Sprintf("task [%s] recover get err: %s", key, err)
	c.HandleErr(key, errors.New(errMsg))
	log.Error(string(debug.Stack()))
}

func (c *common) StartDelay() {
	//启动延迟，错峰运行
	//delaySecond := goToolCommon.RandInt(1, 60)
	//log.Debug(fmt.Sprintf("启动延迟[%d]秒", delaySecond))
	//time.Sleep(time.Duration(1000 * 1000 * 1000 * delaySecond))
}
