package worker

import (
	"errors"
	"github.com/Deansquirrel/goZ9MdDataTrans/object"
	"github.com/Deansquirrel/goZ9MdDataTrans/repository"
)

type bbWorker struct {
	comm *common
}

func NewBbWorker() *bbWorker {
	return &bbWorker{
		comm: NewCommon(),
	}
}

//恢复MdYyInfo
func (w *bbWorker) RestoreMdYyInfo() {
	repOnline, err := repository.NewRepOnline()
	if err != nil {
		w.restoreMdYyInfoHandleErr(err)
		return
	}
	if repOnline == nil {
		errMsg := "repOnline is nil"
		w.restoreMdYyInfoHandleErr(errors.New(errMsg))
		return
	}
	repBb := repository.NewRepBb()
	if repBb == nil {
		errMsg := "repBb is nil"
		w.restoreMdYyInfoHandleErr(errors.New(errMsg))
		return
	}
	for {
		opr, err := repOnline.GetLstMdYyInfoOpr()
		if err != nil {
			w.restoreMdYyInfoHandleErr(err)
			return
		}
		if opr == nil {
			return
		}
		err = repBb.RestoreMdYyInfo(opr)
		if err != nil {
			w.restoreMdYyInfoHandleErr(err)
			return
		}
		err = repOnline.DelLstMdYyInfoOpr(opr.FOprSn)
		if err != nil {
			w.restoreMdYyInfoHandleErr(err)
			return
		}
	}
}

func (w *bbWorker) restoreMdYyInfoHandleErr(err error) {
	w.comm.HandleErr(object.TaskKeyRestoreMdYyInfo, err)
}

//恢复ZxKc
func (w *bbWorker) RestoreZxKc() {
	repOnline, err := repository.NewRepOnline()
	if err != nil {
		w.restoreZxKcHandleErr(err)
		return
	}
	if repOnline == nil {
		errMsg := "repOnline is nil"
		w.restoreZxKcHandleErr(errors.New(errMsg))
		return
	}
	repBb := repository.NewRepBb()
	if repBb == nil {
		errMsg := "repBb is nil"
		w.restoreZxKcHandleErr(errors.New(errMsg))
		return
	}
	for {
		opr, err := repOnline.GetLstZxKcOpr()
		if err != nil {
			w.restoreZxKcHandleErr(err)
			return
		}
		if opr == nil {
			return
		}
		err = repBb.RestoreZxKc(opr)
		if err != nil {
			w.restoreZxKcHandleErr(err)
			return
		}
		err = repOnline.DelLstZxKcOpr(opr.FOprSn)
		if err != nil {
			w.restoreZxKcHandleErr(err)
			return
		}
	}
}

func (w *bbWorker) restoreZxKcHandleErr(err error) {
	w.comm.HandleErr(object.TaskKeyRestoreZxKc, err)
}
