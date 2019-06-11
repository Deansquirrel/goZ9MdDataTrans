package worker

import (
	log "github.com/Deansquirrel/goToolLog"
	"github.com/Deansquirrel/goZ9MdDataTrans/repository"
	"github.com/kataras/iris/core/errors"
)

type onlineWorker struct {
	errChan   chan<- error //错误通知通道
	stateChan chan<- bool  //运行状态变更通道
}

func NewOnlineWorker(errChan chan<- error, stateChan chan<- bool) *onlineWorker {
	return &onlineWorker{
		errChan:   errChan,
		stateChan: stateChan,
	}
}

//刷新心跳时间
func (w *onlineWorker) RefreshHeartBeat() {
	log.Debug("刷新心跳时间")
	repMd := repository.NewRepMd()
	zlCompany, err := repMd.GetZlCompany()
	if err != nil {
		w.errChan <- err
		return
	}
	if zlCompany == nil {
		errMsg := "ZlCompany is nil"
		log.Error(errMsg)
		w.errChan <- errors.New(errMsg)
		return
	}
	repOnline, err := repository.NewRepOnline()
	if err != nil {
		w.errChan <- err
		return
	}
	if repOnline == nil {
		if zlCompany == nil {
			errMsg := "repOnline is nil"
			log.Error(errMsg)
			w.errChan <- errors.New(errMsg)
			return
		}
	}
	w.errChan <- repOnline.UpdateHeartBeat(zlCompany)
}
