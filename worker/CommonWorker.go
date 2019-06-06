package worker

import (
	"github.com/Deansquirrel/goZ9MdDataTrans/global"
	"github.com/Deansquirrel/goZ9MdDataTrans/object"
	"github.com/Deansquirrel/goZ9MdDataTrans/repository"
)

type commonWorker struct {
	errCh   chan<- error //错误通知通道
	stateCh chan<- bool  //运行状态变更通道
}

func NewCommonWorker(errCh chan<- error, stateCh chan<- bool) *commonWorker {
	return &commonWorker{
		errCh:   errCh,
		stateCh: stateCh,
	}
}

//刷新并检查任务配置
func (w *commonWorker) RefreshConfig() {
	repOnline, err := repository.NewRepOnline()
	if err != nil {
		w.errCh <- err
		return
	}

	comm := NewCommon()
	idList := global.TaskKeyList
	for _, id := range idList {
		t := global.TaskList.GetObject(string(id))
		if t == nil {
			comm.StartWorker(id)
			continue
		}
		configCron, err := repOnline.GetTaskCron(id)
		if err != nil {
			w.errCh <- err
		}
		ts := t.(*object.TaskState)
		if configCron != ts.CronStr {
			comm.StopWorker(id)
			comm.StartWorker(id)
		}
	}
	if err == nil {
		w.errCh <- nil
	}
}
