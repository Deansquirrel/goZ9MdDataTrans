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
	//var err error
	//isSelfChange := false
	//
	//defer func() {
	//	if !isSelfChange && err == nil {
	//		w.comm.HandleErr(err)
	//	}
	//}()

	repOnline, err := repository.NewRepOnline()
	if err != nil {
		w.comm.HandleErr(object.TaskKeyRefreshConfig, err)
		return
	}

	for _, id := range global.TaskKeyList {
		if !goToolCron.HasTask(string(id)) {
			w.comm.StartWorker(id)
			continue
		}
		configCron, err := repOnline.GetTaskCron(id)
		if err != nil {
			w.comm.HandleErr(object.TaskKeyRefreshConfig, err)
			continue
		}
		if configCron != goToolCron.CronStr(string(id)) {
			w.restartWorker(id)
			continue
		}
	}

	//comm := NewCommon()
	//idList := global.TaskKeyList
	//
	//for _, id := range idList {
	//	if id == object.TaskKeyRefreshConfig {
	//		isSelfChange = true
	//	} else {
	//		isSelfChange = false
	//	}
	//	t := global.TaskList.GetObject(string(id))
	//	if t == nil {
	//		comm.StartWorker(id)
	//		continue
	//	}
	//	configCron, err := repOnline.GetTaskCron(id)
	//	if err != nil {
	//		w.errCh <- err
	//	}
	//	ts := t.(*object.TaskState)
	//	if configCron != ts.CronStr {
	//		if isSelfChange {
	//			w.restartWorker(id)
	//			return
	//		}
	//		if !w.isTaskRunning(id) {
	//			w.restartWorker(id)
	//		}
	//	}
	//}
}

func (w *commonWorker) restartWorker(key object.TaskKey) {
	comm := NewCommon()
	comm.StopWorker(key)
	comm.StartWorker(key)
}

func (w *commonWorker) isTaskRunning(key object.TaskKey) bool {
	return goToolCron.IsRunning(string(key))
}
