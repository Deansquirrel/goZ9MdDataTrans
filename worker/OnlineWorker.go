package worker

import (
	log "github.com/Deansquirrel/goToolLog"
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
	//TODO 刷新心跳时间
}
