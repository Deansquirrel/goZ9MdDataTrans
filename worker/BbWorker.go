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
func (w *bbWorker) RestoreRestoreMdYyInfo() {
	repOnline, err := repository.NewRepOnline()
	if err != nil {
		w.restoreRestoreMdYyInfoHandleErr(err)
		return
	}
	if repOnline == nil {
		errMsg := "repOnline is nil"
		w.restoreRestoreMdYyInfoHandleErr(errors.New(errMsg))
		return
	}
	repBb := repository.NewRepBb()
	if repBb == nil {
		errMsg := "repBb is nil"
		w.restoreRestoreMdYyInfoHandleErr(errors.New(errMsg))
		return
	}
	for {
		opr, err := repOnline.GetLstMdYyInfoOpr()
		if err != nil {
			w.restoreRestoreMdYyInfoHandleErr(err)
			return
		}
		if opr == nil {
			return
		}
		err = repBb.RestoreMdYyInfo(opr)
		if err != nil {
			w.restoreRestoreMdYyInfoHandleErr(err)
			return
		}
		err = repOnline.DelLstMdYyInfoOpr(opr.FOprSn)
		if err != nil {
			w.restoreRestoreMdYyInfoHandleErr(err)
			return
		}
	}
}

func (w *bbWorker) restoreRestoreMdYyInfoHandleErr(err error) {
	w.comm.HandleErr(object.TaskKeyRestoreMdYyInfo, err)
}
