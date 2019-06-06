package object

import (
	"context"
	"github.com/robfig/cron"
)

type TaskState struct {
	Key     TaskKey
	Cron    *cron.Cron
	CronStr string
	Running bool
	Working bool
	Err     error

	Ctx    context.Context    `json:"-"`
	Cancel context.CancelFunc `json:"-"`
}

//任务执行时间配置表
type TaskCron struct {
	TaskKey         TaskKey //任务标识
	TaskDescription string  //任务描述
	Cron            string  //任务执行cron公式
}

//钉钉消息发送配置
type DingTalkRobotConfigData struct {
	FWebHookKey string
	FAtMobiles  string
	FIsAtAll    int
}
