package global

import (
	"context"
	"github.com/Deansquirrel/goToolCommon"
	"github.com/Deansquirrel/goZ9MdDataTrans/object"
)

const (
	//PreVersion = "0.0.0 Build20190101"
	//TestVersion = "0.0.0 Build20190101"
	Version = "0.0.0 Build20190101"
)

const SecretKey = "Z9MdDataTrans"

var Ctx context.Context
var Cancel func()

//程序启动参数
var Args *object.ProgramArgs

//系统参数
var SysConfig *object.SystemConfig

//TaskList
var TaskList goToolCommon.IObjectManager

var TaskKeyList []object.TaskKey
var TaskTicket map[object.TaskKey]chan struct{}

func TaskKeyListInit() {
	TaskKeyList = make([]object.TaskKey, 0)

	TaskKeyList = append(TaskKeyList, object.TaskKeyRefreshConfig)
	TaskKeyList = append(TaskKeyList, object.TaskKeyRefreshHeartBeat)

	switch SysConfig.RunMode.Mode {
	case object.RunModeMdCollect:
		addMdCollectTask()
	case object.RunModeBbRestore:
		addBbRestoreTask()
	default:
		addMdCollectTask()
	}

	//任务锁（同一任务不可并行）
	TaskTicket = make(map[object.TaskKey]chan struct{})
	for _, k := range TaskKeyList {
		ch := make(chan struct{}, 1)
		TaskTicket[k] = ch
		ch <- struct{}{}
	}
}

func addMdCollectTask() {
	TaskKeyList = append(TaskKeyList, object.TaskKeyRefreshZxKc)
	TaskKeyList = append(TaskKeyList, object.TaskKeyRefreshMdYyInfo)
}

func addBbRestoreTask() {

}
