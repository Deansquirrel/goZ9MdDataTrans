package worker

import (
	"errors"
	"fmt"
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
func (w *commonWorker) RefreshConfig() {
	log.Debug("RefreshConfig")
	checkList := make([]string, 0)
	checkList = append(checkList, w.getCommonTaskKeyList()...)
	switch global.SysConfig.RunMode.Mode {
	case object.RunModeMdCollect:
		checkList = append(checkList, w.getMdTaskKeyList()...)
	case object.RunModeBbRestore:
		checkList = append(checkList, w.getBbTaskKeyList()...)
	default:
		log.Warn(fmt.Sprintf("unknown runmode %s", global.SysConfig.RunMode.Mode))
	}
	repOnline, err := repository.NewRepOnline()
	if err != nil {
		log.Error(fmt.Sprintf("get rep online error: %s", err.Error()))
		return
	}
	for _, id := range checkList {
		if !goToolCron.HasTask(id) {
			continue
		}
		configStr, err := repOnline.GetTaskCron(object.TaskKey(id))
		if err != nil {
			log.Error(fmt.Sprintf("get [%s] config cron error: %s", id, err.Error()))
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
func (w *commonWorker) RefreshHeartBeat() {
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

func (w *commonWorker) refreshHeartBeatHandleErr(err error) {
	log.Error(err.Error())
}
