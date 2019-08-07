package repository

import (
	"errors"
	"fmt"
	"github.com/Deansquirrel/goToolMSSql"
	"github.com/Deansquirrel/goToolMSSqlHelper"
	"github.com/Deansquirrel/goToolSecret"
	"github.com/Deansquirrel/goZ9MdDataTrans/global"
)

import log "github.com/Deansquirrel/goToolLog"

type common struct {
}

func NewCommon() *common {
	return &common{}
}

//获取线上库连接配置
func (c *common) GetOnLineDbConfig() (*goToolMSSql.MSSqlConfig, error) {
	if global.SysConfig.OnLineConfig.Address == "" {
		errMsg := fmt.Sprintf("online db config is empty")
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	rAddress, err := goToolSecret.DecryptFromBase64Format(global.SysConfig.OnLineConfig.Address, global.SecretKey)
	if err != nil {
		errMsg := fmt.Sprintf("online db config decrypt err: %s", err.Error())
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	return goToolMSSqlHelper.GetDBConfigByStr(rAddress)
}

//获取门店库连接配置
func (c *common) GetMdDbConfig() *goToolMSSql.MSSqlConfig {
	return &goToolMSSql.MSSqlConfig{
		Server: global.SysConfig.MdDB.Server,
		Port:   global.SysConfig.MdDB.Port,
		DbName: global.SysConfig.MdDB.DbName,
		User:   global.SysConfig.MdDB.User,
		Pwd:    global.SysConfig.MdDB.Pwd,
	}
}
