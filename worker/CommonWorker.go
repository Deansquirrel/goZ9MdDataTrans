package worker

import (
	"fmt"
	"github.com/Deansquirrel/goZ9MdDataTrans/global"
	"github.com/Deansquirrel/goZ9MdDataTrans/object"
	"github.com/Deansquirrel/goZ9MdDataTrans/repository"
)

type commonWorker struct {
	errCh chan<- error //错误通知通道
}

func NewCommonWorker(errCh chan<- error) *commonWorker {
	return &commonWorker{
		errCh: errCh,
	}
}

//刷新并检查任务配置
func (w *commonWorker) RefreshConfig() {
	var err error
	isSelfChange := false

	defer func() {
		if !isSelfChange && err == nil {
			w.errCh <- nil
		}
	}()

	repOnline, err := repository.NewRepOnline()
	if err != nil {
		w.errCh <- err
		return
	}

	comm := NewCommon()
	idList := global.TaskKeyList

	for _, id := range idList {
		if id == object.TaskKeyRefreshConfig {
			isSelfChange = true
		} else {
			isSelfChange = false
		}
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
			if isSelfChange {
				w.restartWorker(id)
				return
			}
			if !w.isTaskRunning(id) {
				w.restartWorker(id)
			}
		}
	}
}

func (w *commonWorker) restartWorker(key object.TaskKey) {
	comm := NewCommon()
	comm.StopWorker(key)
	comm.StartWorker(key)
}

func (w *commonWorker) isTaskRunning(key object.TaskKey) bool {
	for k, ch := range global.TaskTicket {
		if k == key {
			if len(ch) > 0 {
				return false
			} else {
				return true
			}
		}
	}
	fmt.Println("fsssssssssssssssssssssssssssssssssssssssssssssssssssss")
	return false

	//s := global.TaskList.GetObject(string(key))
	//if s == nil {
	//	return false, errors.New(fmt.Sprintf("task %s err: task state is empty", key))
	//}
	//cs := s.(*object.TaskState)
	//return cs.Working, nil
}
