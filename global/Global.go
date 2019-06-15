package global

import (
	"context"
	"github.com/Deansquirrel/goZ9MdDataTrans/object"
)

const (
	//PreVersion = "0.0.3 Build20190615"
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

var TaskKeyList []object.TaskKey

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
}

func addMdCollectTask() {
	TaskKeyList = append(TaskKeyList, object.TaskKeyRefreshZxKc)
	TaskKeyList = append(TaskKeyList, object.TaskKeyRefreshMdYyInfo)
}

func addBbRestoreTask() {
	TaskKeyList = append(TaskKeyList, object.TaskKeyRestoreZxKc)
	TaskKeyList = append(TaskKeyList, object.TaskKeyRestoreMdYyInfo)
}
