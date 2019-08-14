package worker

import (
	"errors"
	"fmt"
	"github.com/Deansquirrel/goServiceSupportHelper"
	"github.com/Deansquirrel/goToolCron"
	"github.com/Deansquirrel/goToolMSSqlHelper"
	"github.com/Deansquirrel/goToolSVRV3"
	"github.com/Deansquirrel/goZ9MdDataTrans/global"
	"github.com/Deansquirrel/goZ9MdDataTrans/object"
	"github.com/Deansquirrel/goZ9MdDataTrans/repository"
	"time"
)

import log "github.com/Deansquirrel/goToolLog"

type common struct {
}

func NewCommon() *common {
	return &common{}
}

func (c *common) StartService(sType object.RunMode) {
	log.Debug(fmt.Sprintf("RunMode %s", sType))
	for {
		r := c.checkSysConfig()
		if r {
			break
		} else {
			time.Sleep(time.Minute * 30)
		}
	}
	log.Debug(fmt.Sprintf("dbName: %s", global.SysConfig.MdDB.DbName))
	go func() {
		goServiceSupportHelper.SetOtherInfo(
			repository.NewCommon().GetMdDbConfig(),
			1,
			goServiceSupportHelper.SVRV3)
	}()
	switch sType {
	case object.RunModeMdCollect:
		c.addMdWorker()
	case object.RunModeBbRestore:
		c.addBbWorker()
	default:
		log.Warn(fmt.Sprintf("unknown runmode %v", sType))
		global.Cancel()
	}
}

//系统配置检查
func (c *common) checkSysConfig() bool {
	if global.SysConfig.OnLineConfig.Address == "" {
		log.Error("线上库地址不能为空")
		global.Cancel()
		return false
	}
	err := c.refreshLocalDbConfig()
	if err != nil {
		return false
	}
	return true
}

func (c *common) refreshLocalDbConfig() error {
	port := -1
	appType := ""
	clientType := ""

	switch object.RunMode(global.SysConfig.RunMode.Mode) {
	case object.RunModeMdCollect:
		port = 7083
		appType = "83"
		clientType = "8301"
	case object.RunModeBbRestore:
		port = 7091
		appType = "91"
		clientType = "9101"
	default:
		errMsg := fmt.Sprintf("unexpected runmode %s", global.SysConfig.RunMode.Mode)
		log.Error(errMsg)
		global.Cancel()
		return errors.New(errMsg)
	}

	dbConfig, err := goToolSVRV3.GetSQLConfig(global.SysConfig.SvrV3Info.Address, port, appType, clientType)
	if err != nil {
		errMsg := fmt.Sprintf("get dbConfig from svr v3 err: %s", err.Error())
		log.Error(errMsg)
		return errors.New(errMsg)
	}

	if dbConfig == nil {
		errMsg := fmt.Sprintf("get dbConfig from svr v3 return nil")
		log.Error(errMsg)
		return errors.New(errMsg)
	}

	accList, err := goToolSVRV3.GetAccountList(goToolMSSqlHelper.ConvertDbConfigTo2000(dbConfig), appType)
	if err != nil {
		errMsg := fmt.Sprintf("get acc list err: %s", err.Error())
		log.Error(errMsg)
		return errors.New(errMsg)
	}

	if accList == nil || len(accList) <= 0 {
		errMsg := "acc list is empty"
		log.Error(errMsg)
		return errors.New(errMsg)
	}

	global.SysConfig.MdDB.Server = dbConfig.Server
	global.SysConfig.MdDB.Port = dbConfig.Port
	global.SysConfig.MdDB.User = dbConfig.User
	global.SysConfig.MdDB.Pwd = dbConfig.Pwd

	if global.SysConfig.MdDB.DbName != "" {
		flag := false
		for _, acc := range accList {
			if acc == global.SysConfig.MdDB.DbName {
				flag = true
				break
			}
		}
		if !flag {
			log.Warn(fmt.Sprintf("db [%s] is not a effective acc", global.SysConfig.MdDB.DbName))
			global.SysConfig.MdDB.DbName = ""
		}
	}
	if global.SysConfig.MdDB.DbName == "" {
		global.SysConfig.MdDB.DbName = accList[0]
	}
	if global.SysConfig.MdDB.DbName == "" {
		errMsg := fmt.Sprintf("无可用账套")
		log.Error(errMsg)
		return errors.New(errMsg)
	}
	return nil
}

func (c *common) panicHandle(v interface{}) {
	log.Error(fmt.Sprintf("panicHandle: %s", v))
}

func (c *common) addWorker(key string, cmd func(id string)) {
	go func() {
		for {
			repOnline, err := repository.NewRepOnline()
			if err != nil {
				errMsg := fmt.Sprintf("add job [%s] cron error: %s", key, err.Error())
				log.Error(errMsg)
				time.Sleep(time.Minute)
				continue
			}
			cron, err := repOnline.GetTaskCron(object.TaskKey(key))
			if err != nil {
				errMsg := fmt.Sprintf("add job [%s] cron error: %s", key, err.Error())
				log.Error(errMsg)
				time.Sleep(time.Minute)
				continue
			}
			if cron != "" {
				err = goToolCron.AddFunc(key, cron, c.getWorkerFuncReal(key, cmd), c.panicHandle)
				if err != nil {
					errMsg := fmt.Sprintf("add job [%s] error: %s", key, err.Error())
					log.Error(errMsg)
					time.Sleep(time.Minute)
					continue
				}
			} else {
				log.Warn(fmt.Sprintf("job [%s] cron is empty", key))
			}
			break
		}
	}()
}

func (c *common) addCommonWorker() {
	log.Debug("add common worker")
	cWorker := NewCommonWorker()
	c.addWorker("RefreshHeartBeat", cWorker.RefreshHeartBeat)
	c.addWorker("RefreshConfig", cWorker.RefreshConfig)
}

func (c *common) addMdWorker() {
	log.Debug("add md worker")
	c.addCommonWorker()
	mdWorker := NewMdWorker()
	c.addWorker("RefreshZxKc", mdWorker.UpdateZxKc)
	c.addWorker("RefreshMdYyInfo", mdWorker.UpdateMdYyInfo)
}

func (c *common) addBbWorker() {
	log.Debug("add bb worker")
	c.addCommonWorker()
	bbWorker := NewBbWorker()
	c.addWorker("RestoreZxKc", bbWorker.RestoreZxKc)
	c.addWorker("RestoreMdYyInfo", bbWorker.RestoreMdYyInfo)
}

func (c *common) RestartWorker(key string) {
	if goToolCron.HasTask(key) {
		goToolCron.Stop(key)
		goToolCron.DelFunc(key)
	}
	switch key {
	case "RefreshConfig":
		cWorker := NewCommonWorker()
		c.addWorker("RefreshConfig", cWorker.RefreshConfig)
	case "RefreshHeartBeat":
		cWorker := NewCommonWorker()
		c.addWorker("RefreshHeartBeat", cWorker.RefreshHeartBeat)
	case "RefreshMdYyInfo":
		mdWorker := NewMdWorker()
		c.addWorker("RefreshMdYyInfo", mdWorker.UpdateMdYyInfo)
	case "RefreshZxKc":
		bbWorker := NewBbWorker()
		c.addWorker("RefreshZxKc", bbWorker.RestoreZxKc)
	case "RestoreMdYyInfo":
		bbWorker := NewBbWorker()
		c.addWorker("RestoreMdYyInfo", bbWorker.RestoreMdYyInfo)
	case "RestoreZxKc":
		bbWorker := NewBbWorker()
		c.addWorker("RestoreZxKc", bbWorker.RestoreZxKc)
	default:
		log.Warn(fmt.Sprintf("unknown task key %s", key))
	}
}

func (c *common) getWorkerFuncReal(key string, cmd func(id string)) func() {
	return goServiceSupportHelper.NewJob().FormatSSJob(key, c.formatWorkFunc(key, cmd))
}

func (c *common) formatWorkFunc(key string, cmd func(id string)) func(string) {
	return func(id string) {
		log.Debug(fmt.Sprintf("task %s[%s] start", key, id))
		defer log.Debug(fmt.Sprintf("task %s[%s] complete", key, id))
		cmd(id)
	}
}
