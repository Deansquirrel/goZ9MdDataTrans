package worker

import (
	"github.com/Deansquirrel/goServiceSupportHelper"
	"github.com/Deansquirrel/goZ9MdDataTrans/repository"
)

import log "github.com/Deansquirrel/goToolLog"

type bbWorker struct {
	comm *common
}

func NewBbWorker() *bbWorker {
	return &bbWorker{
		comm: NewCommon(),
	}
}

//恢复MdYyInfo
func (w *bbWorker) RestoreMdYyInfo(id string) {
	repOnline, err := repository.NewRepOnline()
	if err != nil {
		log.Error(err.Error())
		_ = goServiceSupportHelper.JobErrRecord(id, err.Error())
		return
	}
	if repOnline == nil {
		errMsg := "repOnline is nil"
		log.Error(errMsg)
		_ = goServiceSupportHelper.JobErrRecord(id, errMsg)
		return
	}
	repBb := repository.NewRepBb()
	if repBb == nil {
		errMsg := "repBb is nil"
		log.Error(errMsg)
		_ = goServiceSupportHelper.JobErrRecord(id, errMsg)
		return
	}
	for {
		opr, err := repOnline.GetLstMdYyInfoOpr()
		if err != nil {
			log.Error(err.Error())
			_ = goServiceSupportHelper.JobErrRecord(id, err.Error())
			return
		}
		if opr == nil {
			return
		}
		err = repBb.RestoreMdYyInfo(opr)
		if err != nil {
			log.Error(err.Error())
			_ = goServiceSupportHelper.JobErrRecord(id, err.Error())
			return
		}
		err = repOnline.DelLstMdYyInfoOpr(opr.FOprSn)
		if err != nil {
			log.Error(err.Error())
			_ = goServiceSupportHelper.JobErrRecord(id, err.Error())
			return
		}
	}
}

//恢复ZxKc
func (w *bbWorker) RestoreZxKc(id string) {
	repOnline, err := repository.NewRepOnline()
	if err != nil {
		log.Error(err.Error())
		_ = goServiceSupportHelper.JobErrRecord(id, err.Error())
		return
	}
	if repOnline == nil {
		errMsg := "repOnline is nil"
		log.Error(errMsg)
		_ = goServiceSupportHelper.JobErrRecord(id, errMsg)
		return
	}
	repBb := repository.NewRepBb()
	if repBb == nil {
		errMsg := "repBb is nil"
		log.Error(errMsg)
		_ = goServiceSupportHelper.JobErrRecord(id, errMsg)
		return
	}
	for {
		opr, err := repOnline.GetLstZxKcOpr()
		if err != nil {
			log.Error(err.Error())
			_ = goServiceSupportHelper.JobErrRecord(id, err.Error())
			return
		}
		if opr == nil {
			return
		}
		err = repBb.RestoreZxKc(opr)
		if err != nil {
			log.Error(err.Error())
			_ = goServiceSupportHelper.JobErrRecord(id, err.Error())
			return
		}
		err = repOnline.DelLstZxKcOpr(opr.FOprSn)
		if err != nil {
			log.Error(err.Error())
			_ = goServiceSupportHelper.JobErrRecord(id, err.Error())
			return
		}
	}
}
