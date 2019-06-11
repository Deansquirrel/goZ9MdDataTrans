package worker

import (
	"context"
	"errors"
	"fmt"
	"github.com/Deansquirrel/goToolCommon"
	"github.com/Deansquirrel/goZ9MdDataTrans/global"
	"github.com/Deansquirrel/goZ9MdDataTrans/object"
	"github.com/Deansquirrel/goZ9MdDataTrans/repository"
	"github.com/robfig/cron"
	"runtime/debug"
	"strings"
	"sync"
)

import log "github.com/Deansquirrel/goToolLog"

const (
	//TODO 钉钉机器人消息Key待参数化
	webHook = "7685d3a4b1630c8cb0028ee9cd17beae1578869f8ddf0b7871ecbb8b0743f8ce"
)

var syncLock sync.Mutex

type common struct {
}

func NewCommon() *common {
	return &common{}
}

//启动工作线程
func (c *common) StartWorker(key object.TaskKey) {
	s := &object.TaskState{
		Key:     key,
		Cron:    nil,
		CronStr: "",
		Running: false,
		Working: false,
		Err:     nil,
	}

	syncLock.Lock()
	defer syncLock.Unlock()

	task := global.TaskList.GetObject(string(s.Key))
	if task != nil {
		return
	}
	s.Ctx, s.Cancel = context.WithCancel(context.Background())
	global.TaskList.Register() <- goToolCommon.NewObject(string(s.Key), s)

	errCh := make(chan error)
	stateCh := make(chan bool)

	go func() {
		defer func() {
			err := recover()
			if err != nil {
				log.Error(fmt.Sprintf("task %s recover err: %s", key, err))
				log.Error(string(debug.Stack()))
			}
		}()
		defer func() {
			close(errCh)
			close(stateCh)
		}()
		for {
			select {
			case err := <-errCh:
				s := global.TaskList.GetObject(string(key))
				if s == nil {
					errCh <- errors.New(fmt.Sprintf("task %s err: task state is empty", key))
					return
				}
				cs := s.(*object.TaskState)
				if err != nil {
					//c.errHandle(err)
					c.sendMsg(err.Error())
				} else {
					if cs != nil && cs.Err != nil {
						msg := fmt.Sprintf("Task resume %s", key)
						log.Warn(msg)
						c.sendMsg(msg)
					}
				}
				if cs != nil {
					cs.Err = err
				}
			case b := <-stateCh:
				s := global.TaskList.GetObject(string(key))
				if s == nil {
					errCh <- errors.New(fmt.Sprintf("task %s err: task state is empty", key))
					return
				}
				cs := s.(*object.TaskState)
				if cs != nil {
					cs.Working = b
				}
			case <-s.Ctx.Done():
				return
			case <-global.Ctx.Done():
				return
			}
		}
	}()

	repOnline, err := repository.NewRepOnline()
	if err != nil {
		errCh <- err
		return
	}
	cronStr, err := repOnline.GetTaskCron(s.Key)
	if err != nil {
		s.CronStr = ""
		errCh <- err
		return
	}
	cronStr = strings.Trim(cronStr, " ")
	if cronStr == "" {
		errCh <- errors.New(fmt.Sprintf("task %s start err: cron is empty", key))
		return
	}
	s.CronStr = cronStr
	cr := cron.New()
	err = cr.AddFunc(s.CronStr, func() {
		guid := goToolCommon.Guid()
		log.Debug(fmt.Sprintf("task %s[%s] start", key, guid))

		defer func() {
			//错误处理（panic）
			err := recover()
			if err != nil {
				errMsg := fmt.Sprintf("task recover get err: %s", err)
				log.Error(errMsg)
				log.Error(string(debug.Stack()))
				//errCh <- errors.New(errMsg)
				c.sendMsg(errMsg)
			}
			log.Debug(fmt.Sprintf("task %s[%s] complete", key, guid))
		}()

		global.TaskSyncLockList[key].Lock()
		defer global.TaskSyncLockList[key].Unlock()

		f := c.getWorkerFunc(s.Key, errCh, stateCh)
		if f != nil {
			f()
		} else {
			errMsg := fmt.Sprintf("task %s[%s] start err: func is nil", key, guid)
			log.Error(errMsg)
			errCh <- errors.New(errMsg)
			return
		}
	})
	if err != nil {
		errCh <- err
		return
	}

	log.Debug(fmt.Sprintf("start worker %s cron %s", key, cronStr))

	cr.Start()
	s.Running = true
	s.Cron = cr
}

//停止工作线程
func (c *common) StopWorker(key object.TaskKey) {
	t := global.TaskList.GetObject(string(key))
	if t == nil {
		return
	}
	global.TaskList.Unregister() <- string(key)
	log.Debug(fmt.Sprintf("stop worker %s", key))
	ts := t.(*object.TaskState)
	if ts.Running {
		ts.Running = false
		ts.Cron.Stop()
		ts.Cancel()
	}
}

//获取任务执行函数
func (c *common) getWorkerFunc(key object.TaskKey, errCh chan<- error, stateCh chan<- bool) func() {
	//TODO 获取任务执行函数
	switch key {
	case object.TaskKeyRefreshConfig:
		return NewCommonWorker(errCh, stateCh).RefreshConfig
	case object.TaskKeyRefreshHeartBeat:
		return NewOnlineWorker(errCh, stateCh).RefreshHeartBeat
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
func (c *common) errHandle(err error) {
	if err != nil {
		c.sendMsg(err.Error())
	}
}

//发送钉钉消息
func (c *common) sendMsg(msg string) {
	dt := object.NewDingTalkRobot(&object.DingTalkRobotConfigData{
		FWebHookKey: webHook,
		FAtMobiles:  "",
		FIsAtAll:    0,
	})
	sendErr := dt.SendMsg(msg)
	if sendErr != nil {
		log.Error(fmt.Sprintf("send msg error,msg: %s, error: %s", msg, sendErr.Error()))
	}
}

func (c *common) HandlePanic() {
	err := recover()
	if err != nil {
		log.Error(fmt.Sprintf("recover get err: %s", err))
		log.Error(string(debug.Stack()))
	}
}
