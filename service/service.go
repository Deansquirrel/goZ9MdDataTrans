package service

import (
	log "github.com/Deansquirrel/goToolLog"
	"github.com/Deansquirrel/goZ9MdDataTrans/object"
	"github.com/Deansquirrel/goZ9MdDataTrans/worker"
)

//启动服务内容
func StartService() error {
	log.Debug("Start Service")
	defer log.Debug("Start Service Complete")

	comm := worker.NewCommon()
	comm.StartWorker(object.TaskKeyRefreshConfig)

	return nil
}
