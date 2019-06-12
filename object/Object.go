package object

import (
	"time"
)

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

//ZlCompany
type ZlCompany struct {
	FCoId          int       //分支机构ID
	FCoAb          string    //名称
	FCoCode        string    //通道码
	FCoUserAb      string    //主名称
	FCoUserCode    string    //主代码
	FCoFunc        int       //系统类型
	FCoAccStartDay time.Time //建账营业日
}
