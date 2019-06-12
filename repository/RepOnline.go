package repository

import (
	"errors"
	"fmt"
	"github.com/Deansquirrel/goToolCommon"
	"github.com/Deansquirrel/goToolMSSql"
	"github.com/Deansquirrel/goZ9MdDataTrans/object"
)

import log "github.com/Deansquirrel/goToolLog"

const (
	sqlGetTaskCronList = "" +
		"SELECT [taskkey],[taskdescription],[cron] " +
		"FROM [taskcron]"

	sqlRefreshHeartBeat = "" +
		"IF EXISTS (SELECT * FROM [heartbeat] WHERE [mdid] = ?) " +
		"BEGIN " +
		"UPDATE [heartbeat] " +
		"SET [heartbeat] = getDate() " +
		"WHERE [mdid] = ? " +
		"END " +
		"ELSE " +
		"BEGIN " +
		"INSERT INTO [heartbeat]([mdid],[mdname],[heartbeat]) " +
		"SELECT ?,?,getDate() " +
		"END"

	sqlUpdateMdYyInfo = "" +
		"INSERT INTO [mdyyinfo]([mdid],[yyr],[tc],[sr],[oprtime]) " +
		"VALUES(?,?,?,?,?)"
)

type repOnline struct {
	dbConfig *goToolMSSql.MSSqlConfig
}

func NewRepOnline() (*repOnline, error) {
	c := common{}
	dbConfig, err := c.GetOnLineDbConfig()
	if err != nil {
		return nil, err
	}
	if dbConfig == nil {
		errMsg := "get rep online err: dbConfig is nil"
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	return &repOnline{
		dbConfig: dbConfig,
	}, nil
}

//获取TaskCron配置
func (r *repOnline) GetTaskCron(key object.TaskKey) (string, error) {
	cronList, err := r.GetTaskCronList()
	if err != nil {
		return "", err
	}
	for _, tc := range cronList {
		if tc.TaskKey == key {
			return tc.Cron, nil
		}
	}
	return "", nil
}

//获取TaskCron配置列表
func (r *repOnline) GetTaskCronList() ([]*object.TaskCron, error) {
	comm := NewCommon()
	rows, err := comm.GetRowsBySQL(r.dbConfig, sqlGetTaskCronList)
	if err != nil {
		return nil, err
	}
	resultList := make([]*object.TaskCron, 0)
	var key object.TaskKey
	var desc, cron string
	for rows.Next() {
		err := rows.Scan(&key, &desc, &cron)
		if err != nil {
			errMsg := fmt.Sprintf("rep online GetTaskCronList read data err: %s", err.Error())
			log.Error(errMsg)
			return nil, errors.New(errMsg)
		}
		resultList = append(resultList, &object.TaskCron{
			TaskKey:         key,
			TaskDescription: desc,
			Cron:            cron,
		})
	}
	if rows.Err() != nil {
		errMsg := fmt.Sprintf("rep online GetTaskCronList read data err: %s", rows.Err().Error())
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	return resultList, nil
}

//刷新心跳
func (r *repOnline) UpdateHeartBeat(company *object.ZlCompany) error {
	if company == nil {
		errMsg := "company is nil"
		log.Error(errMsg)
		return errors.New(errMsg)
	}
	comm := NewCommon()
	return comm.SetRowsBySQL(r.dbConfig, sqlRefreshHeartBeat, company.FCoId, company.FCoId, company.FCoId, company.FCoAb)
}

func (r *repOnline) UpdateMdYyInfo(info *object.MdYyInfo) error {
	comm := NewCommon()
	return comm.SetRowsBySQL(r.dbConfig, sqlUpdateMdYyInfo,
		info.FMdId,
		goToolCommon.GetDateStr(info.FYyr),
		fmt.Sprintf("%v", info.FTc),
		fmt.Sprintf("%v", info.FSr),
		goToolCommon.GetDateTimeStr(info.FOprTime))
}
