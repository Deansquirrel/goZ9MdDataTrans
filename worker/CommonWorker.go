package worker

import (
	"errors"
	"fmt"
	"github.com/Deansquirrel/goServiceSupportHelper"
	"github.com/Deansquirrel/goToolCron"
	"github.com/Deansquirrel/goZ9MdDataTrans/global"
	"github.com/Deansquirrel/goZ9MdDataTrans/object"
	"github.com/Deansquirrel/goZ9MdDataTrans/repository"
)

import log "github.com/Deansquirrel/goToolLog"

type commonWorker struct {
}

func NewCommonWorker() *commonWorker {
	return &commonWorker{}
}

//刷新并检查任务配置
func (w *commonWorker) RefreshConfig(id string) {
	log.Debug("RefreshConfig")
	checkList := make([]string, 0)
	checkList = append(checkList, w.getCommonTaskKeyList()...)
	switch global.SysConfig.RunMode.Mode {
	case object.RunModeMdCollect:
		checkList = append(checkList, w.getMdTaskKeyList()...)
	case object.RunModeBbRestore:
		checkList = append(checkList, w.getBbTaskKeyList()...)
	default:
		errMsg := fmt.Sprintf("unknown runmode %s", global.SysConfig.RunMode.Mode)
		log.Warn(errMsg)
		_ = goServiceSupportHelper.JobErrRecord(id, errMsg)
		return
	}
	repOnline, err := repository.NewRepOnline()
	if err != nil {
		errMsg := fmt.Sprintf("get rep online error: %s", err.Error())
		log.Error(errMsg)
		_ = goServiceSupportHelper.JobErrRecord(id, errMsg)
		return
	}
	for _, id := range checkList {
		if !goToolCron.HasTask(id) {
			continue
		}
		configStr, err := repOnline.GetTaskCron(object.TaskKey(id))
		if err != nil {
			errMsg := fmt.Sprintf("get [%s] config cron error: %s", id, err.Error())
			log.Error(errMsg)
			_ = goServiceSupportHelper.JobErrRecord(id, errMsg)
			return
		}
		currCron := goToolCron.CronStr(id)
		if configStr != currCron {
			NewCommon().RestartWorker(id)
		}
	}
}

func (w *commonWorker) getCommonTaskKeyList() []string {
	rList := make([]string, 0)
	rList = append(rList, "RefreshHeartBeat")
	rList = append(rList, "RefreshConfig")
	return rList
}

func (w *commonWorker) getMdTaskKeyList() []string {
	rList := make([]string, 0)
	rList = append(rList, "RefreshZxKc")
	rList = append(rList, "RefreshMdYyInfo")
	return rList
}

func (w *commonWorker) getBbTaskKeyList() []string {
	rList := make([]string, 0)
	rList = append(rList, "RestoreZxKc")
	rList = append(rList, "RestoreMdYyInfo")
	return rList
}

//刷新心跳时间
func (w *commonWorker) RefreshHeartBeat(id string) {
	log.Debug("刷新心跳时间")
	repMd := repository.NewRepMd()
	zlCompany, err := repMd.GetZlCompany()
	if err != nil {
		w.refreshHeartBeatHandleErr(id, err)
		return
	}
	if zlCompany == nil {
		errMsg := "ZlCompany is nil"
		w.refreshHeartBeatHandleErr(id, errors.New(errMsg))
		return
	}
	repOnline, err := repository.NewRepOnline()
	if err != nil {
		w.refreshHeartBeatHandleErr(id, err)
		return
	}
	if repOnline == nil {
		errMsg := "repOnline is nil"
		w.refreshHeartBeatHandleErr(id, errors.New(errMsg))
		return
	}
	err = repOnline.UpdateHeartBeat(zlCompany)
	if err != nil {
		w.refreshHeartBeatHandleErr(id, err)
		return
	}
}

func (w *commonWorker) refreshHeartBeatHandleErr(id string, err error) {
	log.Error(err.Error())
	_ = goServiceSupportHelper.JobErrRecord(id, err.Error())
}
