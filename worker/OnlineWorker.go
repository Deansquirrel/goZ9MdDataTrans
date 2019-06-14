package worker

import (
	"errors"
	"github.com/Deansquirrel/goZ9MdDataTrans/object"
	"github.com/Deansquirrel/goZ9MdDataTrans/repository"
	"time"
)

import log "github.com/Deansquirrel/goToolLog"

var zxkcsj time.Time

func init() {
	zxkcsj = repository.NewCommon().GetDefaultOprTime()
}

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
		errMsg := "repOnline is nil"
		w.refreshHeartBeatHandleErr(errors.New(errMsg))
		return
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

func (w *onlineWorker) UpdateZxKc() {
	log.Debug("刷新最新库存变动")
	repMd := repository.NewRepMd()
	repOnline, err := repository.NewRepOnline()
	if err != nil {
		w.updateZxKcHandleErr(err)
		return
	}
	kcList, lst, err := repMd.GetZxKcInfo(zxkcsj)
	if err != nil {
		w.updateZxKcHandleErr(err)
		return
	}
	cInfo, err := repMd.GetZlCompany()
	if err != nil {
		w.updateZxKcHandleErr(err)
		return
	}
	err = repOnline.UpdateZxKc(cInfo.FCoId, kcList)
	if err != nil {
		w.updateZxKcHandleErr(err)
		return
	}
	err = repOnline.UpdateKcLastUpdate(cInfo.FCoId)
	if err != nil {
		w.updateZxKcHandleErr(err)
		return
	}
	zxkcsj = lst
}

func (w *onlineWorker) updateZxKcHandleErr(err error) {
	w.comm.HandleErr(object.TaskKeyRefreshZxKc, err)
}
