package worker

import (
	"github.com/Deansquirrel/goToolMSSqlHelper"
	"github.com/Deansquirrel/goZ9MdDataTrans/repository"
	"time"
)

import log "github.com/Deansquirrel/goToolLog"

var zxKcSj time.Time

func init() {
	zxKcSj = goToolMSSqlHelper.GetDefaultOprTime()
}

type mdWorker struct {
}

func NewMdWorker() *mdWorker {
	return &mdWorker{}
}

func (w *mdWorker) UpdateMdYyInfo() {
	log.Debug("刷新门店营业信息")
	repMd := repository.NewRepMd()
	repOnline, err := repository.NewRepOnline()
	if err != nil {
		log.Error(err.Error())
		return
	}
	list, err := repMd.GetMdYyInfo()
	if err != nil {
		log.Error(err.Error())
		return
	}

	for _, info := range list {
		err = repOnline.UpdateMdYyInfo(info)
		if err != nil {
			log.Error(err.Error())
		}
	}
}

func (w *mdWorker) UpdateZxKc() {
	log.Debug("刷新最新库存变动")
	repMd := repository.NewRepMd()
	repOnline, err := repository.NewRepOnline()
	if err != nil {
		log.Error(err.Error())
		return
	}
	kcList, lst, err := repMd.GetZxKcInfo(zxKcSj)
	if err != nil {
		log.Error(err.Error())
		return
	}
	cInfo, err := repMd.GetZlCompany()
	if err != nil {
		log.Error(err.Error())
		return
	}
	err = repOnline.UpdateZxKc(cInfo.FCoId, kcList)
	if err != nil {
		log.Error(err.Error())
		return
	}
	err = repOnline.UpdateKcLastUpdate(cInfo.FCoId)
	if err != nil {
		log.Error(err.Error())
		return
	}
	zxKcSj = lst
}
