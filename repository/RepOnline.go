package repository

import (
	"errors"
	"fmt"
	"github.com/Deansquirrel/goToolCommon"
	"github.com/Deansquirrel/goToolMSSql"
	"github.com/Deansquirrel/goToolMSSqlHelper"
	"github.com/Deansquirrel/goZ9MdDataTrans/global"
	"github.com/Deansquirrel/goZ9MdDataTrans/object"
	"time"
)

import log "github.com/Deansquirrel/goToolLog"

const (
	sqlGetTaskCronList = "" +
		"SELECT [taskkey],[taskdescription],[cron] " +
		"FROM [taskcron]"

	sqlRefreshHeartBeat = "" +
		"IF EXISTS (SELECT * FROM [heartbeat] WHERE [mdid] = ?) " +
		"	BEGIN " +
		"		UPDATE [heartbeat] " +
		"		SET [heartbeat] = getDate(),clientVersion = ? " +
		"		WHERE [mdid] = ? " +
		"	END " +
		"ELSE " +
		"	BEGIN " +
		"		INSERT INTO [heartbeat]([mdid],[mdname],[clientVersion],[heartbeat]) " +
		"		SELECT ?,?,?,getDate() " +
		"	END"

	sqlUpdateMdYyInfo = "" +
		"INSERT INTO [mdyyinfo]([mdid],[yyr],[tc],[sr],[oprtime]) " +
		"VALUES(?,?,?,?,?)"

	sqlUpdateZxKc = "" +
		"INSERT INTO [zxkc]([mdid],[hpid],[sl],[oprtime]) " +
		"VALUES (?,?,?,?)"

	sqlUpdateKcLastUpdate = "" +
		"IF EXISTS (SELECT * FROM [kclastupdate] WHERE [mdid]=?)  " +
		"	BEGIN " +
		"		UPDATE [kclastupdate] " +
		"		SET [lastupdate] = getDate() " +
		"		WHERE [mdid] = ? " +
		"	END " +
		"ELSE " +
		"	BEGIN " +
		"		INSERT INTO [kclastupdate]([mdid],[lastupdate]) " +
		"		VALUES (?,getDate()) " +
		"	END"

	sqlGetLstMdYyInfoOpr = "" +
		"SELECT top 1 [oprsn],[mdid],[yyr],[tc],[sr],[oprtime] " +
		"FROM [mdyyinfo] " +
		"ORDER BY [oprsn] ASC"

	sqlDelLstMdYyInfoOpr = "" +
		"DELETE FROM [mdyyinfo] " +
		"WHERE [oprsn]=?"

	sqlGetLstZxKcOpr = "" +
		"SELECT TOP 1 [oprsn],[mdid],[hpid],[sl],[oprtime] " +
		"FROM [zxkc] " +
		"ORDER BY [oprsn] ASC"

	sqlDelLstZxKcOpr = "" +
		"DELETE FROM [zxkc] " +
		"WHERE [oprsn]=?"
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
	rows, err := goToolMSSqlHelper.GetRowsBySQL(r.dbConfig, sqlGetTaskCronList)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
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
	return goToolMSSqlHelper.SetRowsBySQL(r.dbConfig, sqlRefreshHeartBeat,
		company.FCoId, global.Version, company.FCoId, company.FCoId, company.FCoAb, global.Version)
}

func (r *repOnline) UpdateMdYyInfo(info *object.MdYyInfo) error {
	return goToolMSSqlHelper.SetRowsBySQL(r.dbConfig, sqlUpdateMdYyInfo,
		info.FMdId,
		goToolCommon.GetDateStr(info.FYyr),
		info.FTc,
		info.FSr,
		goToolCommon.GetDateTimeStr(info.FOprTime))
}

func (r *repOnline) UpdateZxKc(fMdId int, kcList []*object.ZxKc) error {
	var err error
	for _, kc := range kcList {
		err = goToolMSSqlHelper.SetRowsBySQL(r.dbConfig, sqlUpdateZxKc,
			fMdId,
			kc.FHpId,
			kc.FSl,
			goToolCommon.GetDateTimeStrWithMillisecond(kc.FOprTime))
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *repOnline) UpdateKcLastUpdate(mdId int) error {
	return goToolMSSqlHelper.SetRowsBySQL(r.dbConfig, sqlUpdateKcLastUpdate,
		mdId, mdId, mdId)
}

func (r *repOnline) GetLstMdYyInfoOpr() (*object.MdYyInfoOpr, error) {
	rows, err := goToolMSSqlHelper.GetRowsBySQL(r.dbConfig, sqlGetLstMdYyInfoOpr)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	var oprSn, mdId, tc int
	var yyr, oprTime time.Time
	var sr float64

	rList := make([]*object.MdYyInfoOpr, 0)
	for rows.Next() {
		err = rows.Scan(&oprSn, &mdId, &yyr, &tc, &sr, &oprTime)
		if err != nil {
			errMsg := fmt.Sprintf("Get Lst MdYyInfoOpr, read data err: %s", err.Error())
			log.Error(errMsg)
			return nil, errors.New(errMsg)
		}
		rList = append(rList, &object.MdYyInfoOpr{
			FOprSn:   oprSn,
			FMdId:    mdId,
			FYyr:     yyr,
			FTc:      tc,
			FSr:      sr,
			FOprTime: oprTime,
		})
	}
	if rows.Err() != nil {
		errMsg := fmt.Sprintf("Get Lst MdYyInfoOpr, read data err: %s", rows.Err().Error())
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	if len(rList) > 0 {
		return rList[0], nil
	} else {
		return nil, nil
	}
}

func (r *repOnline) DelLstMdYyInfoOpr(sn int) error {
	return goToolMSSqlHelper.SetRowsBySQL(r.dbConfig, sqlDelLstMdYyInfoOpr, sn)
}

func (r *repOnline) GetLstZxKcOpr() (*object.ZxKcOpr, error) {
	rows, err := goToolMSSqlHelper.GetRowsBySQL(r.dbConfig, sqlGetLstZxKcOpr)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	var oprSn, mdId, hpId int
	var sl float32
	var oprTime time.Time
	rList := make([]*object.ZxKcOpr, 0)
	for rows.Next() {
		err = rows.Scan(&oprSn, &mdId, &hpId, &sl, &oprTime)
		if err != nil {
			errMsg := fmt.Sprintf("get lst zxkc opr,read data err: %s", err.Error())
			log.Error(errMsg)
			return nil, errors.New(errMsg)
		}
		rList = append(rList, &object.ZxKcOpr{
			FOprSn:   oprSn,
			FMdId:    mdId,
			FHpId:    hpId,
			FSl:      sl,
			FOprTime: oprTime,
		})
	}
	if rows.Err() != nil {
		errMsg := fmt.Sprintf("get lst zxkc opr,read data err: %s", rows.Err().Error())
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	if len(rList) > 0 {
		return rList[0], nil
	} else {
		return nil, nil
	}
}

func (r *repOnline) DelLstZxKcOpr(sn int) error {
	return goToolMSSqlHelper.SetRowsBySQL(r.dbConfig, sqlDelLstZxKcOpr, sn)
}
