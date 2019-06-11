package worker

import (
	"errors"
	"fmt"
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
	w.stateCh <- true

	var err error
	isSelfChange := false

	defer func() {
		if !isSelfChange {
			w.stateCh <- false
			if err == nil {
				w.errCh <- nil
			}
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
			ws, err := w.getTaskWorkState(id)
			if err != nil {
				w.errCh <- err
				continue
			}
			if !ws {
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

func (w *commonWorker) getTaskWorkState(key object.TaskKey) (bool, error) {
	s := global.TaskList.GetObject(string(key))
	if s == nil {
		return false, errors.New(fmt.Sprintf("task %s err: task state is empty", key))
	}
	cs := s.(*object.TaskState)
	return cs.Working, nil
}
