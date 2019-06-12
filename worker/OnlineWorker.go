package worker

import (
	"errors"
	"github.com/Deansquirrel/goZ9MdDataTrans/object"
	"github.com/Deansquirrel/goZ9MdDataTrans/repository"
)

import log "github.com/Deansquirrel/goToolLog"

type onlineWorker struct {
	comm *common
}

func NewOnlineWorker() *onlineWorker {
	return &onlineWorker{
		comm: NewCommon(),
	}
}

//刷新心跳时间
func (w *onlineWorker) RefreshHeartBeat() {
	log.Debug("刷新心跳时间")
	repMd := repository.NewRepMd()
	zlCompany, err := repMd.GetZlCompany()
	if err != nil {
		w.comm.HandleErr(object.TaskKeyRefreshHeartBeat, err)
		return
	}
	if zlCompany == nil {
		errMsg := "ZlCompany is nil"
		w.comm.HandleErr(object.TaskKeyRefreshHeartBeat, errors.New(errMsg))
		return
	}
	repOnline, err := repository.NewRepOnline()
	if err != nil {
		w.comm.HandleErr(object.TaskKeyRefreshHeartBeat, err)
		return
	}
	if repOnline == nil {
		if zlCompany == nil {
			errMsg := "repOnline is nil"
			w.comm.HandleErr(object.TaskKeyRefreshHeartBeat, errors.New(errMsg))
			return
		}
	}
	w.comm.HandleErr(object.TaskKeyRefreshHeartBeat, repOnline.UpdateHeartBeat(zlCompany))
}
