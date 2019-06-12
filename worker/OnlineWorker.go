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
		w.refreshHeartBeatHandleErr(err)
		return
	}
	if zlCompany == nil {
		errMsg := "ZlCompany is nil"
		w.refreshHeartBeatHandleErr(errors.New(errMsg))
		return
	}
	repOnline, err := repository.NewRepOnline()
	if err != nil {
		w.refreshHeartBeatHandleErr(err)
		return
	}
	if repOnline == nil {
		if zlCompany == nil {
			errMsg := "repOnline is nil"
			w.refreshHeartBeatHandleErr(errors.New(errMsg))
			return
		}
	}
	err = repOnline.UpdateHeartBeat(zlCompany)
	if err != nil {
		w.refreshHeartBeatHandleErr(err)
		return
	}
}

func (w *onlineWorker) refreshHeartBeatHandleErr(err error) {
	w.comm.HandleErr(object.TaskKeyRefreshHeartBeat, err)
}

func (w *onlineWorker) UpdateMdYyInfo() {
	log.Debug("刷新门店营业信息")
	repMd := repository.NewRepMd()
	repOnline, err := repository.NewRepOnline()
	if err != nil {
		w.updateMdYyInfoHandleErr(err)
		return
	}
	list, err := repMd.GetMdYyInfo()
	if err != nil {
		w.updateMdYyInfoHandleErr(err)
		return
	}

	for _, info := range list {
		err = repOnline.UpdateMdYyInfo(info)
		if err != nil {
			w.updateMdYyInfoHandleErr(err)
		}
	}
}

func (w *onlineWorker) updateMdYyInfoHandleErr(err error) {
	w.comm.HandleErr(object.TaskKeyRefreshMdYyInfo, err)
}
