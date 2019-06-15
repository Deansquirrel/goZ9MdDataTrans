package worker

import (
	"github.com/Deansquirrel/goToolCron"
	"github.com/Deansquirrel/goZ9MdDataTrans/global"
	"github.com/Deansquirrel/goZ9MdDataTrans/object"
	"github.com/Deansquirrel/goZ9MdDataTrans/repository"
)

type commonWorker struct {
	comm *common
}

func NewCommonWorker() *commonWorker {
	return &commonWorker{
		comm: NewCommon(),
	}
}

//刷新并检查任务配置
func (w *commonWorker) RefreshConfig() {

	w.comm.StartDelay()

	repOnline, err := repository.NewRepOnline()
	if err != nil {
		w.refreshConfigHandleErr(err)
		return
	}

	for _, id := range global.TaskKeyList {
		if !goToolCron.HasTask(string(id)) {
			w.comm.StartWorker(id)
			continue
		}
		configCron, err := repOnline.GetTaskCron(id)
		if err != nil {
			w.refreshConfigHandleErr(err)
			continue
		}
		if configCron != goToolCron.CronStr(string(id)) {
			w.restartWorker(id)
			continue
		}
	}
}

func (w *commonWorker) refreshConfigHandleErr(err error) {
	w.comm.HandleErr(object.TaskKeyRefreshConfig, err)
}

func (w *commonWorker) restartWorker(key object.TaskKey) {
	comm := NewCommon()
	comm.StopWorker(key)
	comm.StartWorker(key)
}

func (w *commonWorker) isTaskRunning(key object.TaskKey) bool {
	return goToolCron.IsRunning(string(key))
}
