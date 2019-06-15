package repository

import (
	"errors"
	"fmt"
	"github.com/Deansquirrel/goToolCommon"
	"github.com/Deansquirrel/goToolMSSql2000"
	"github.com/Deansquirrel/goZ9MdDataTrans/object"
	"time"
)

import log "github.com/Deansquirrel/goToolLog"

const (
	sqlGetMdZlCompany = "" +
		"SELECT [coid],[coab],[cocode],[couserab],[cousercode]" +
		",[cofunc],[coaccstartday] " +
		"FROM [zlcompany]"

	sqlGetMdYyInfo = "" +
		"SELECT [ckmdid],[ckyyr],sum(case when [ckcxbj]=1 then -1 else 1 end) as [num],sum([ckcjje]) as [srmy] " +
		"FROM [z3xsckt] WITH(NOLOCK) " +
		"GROUP BY [ckmdid],[ckyyr]"

	sqlGetZxKcInfo = "" +
		"SELECT [tzhpid],[tzsl],[tzbdsj] " +
		"FROM [z3xttz] " +
		"WHERE [tzckid] = 0 AND [tzbdsj] > ?"
)

type repMd struct {
	dbConfig *goToolMSSql2000.MSSqlConfig
}

func NewRepMd() *repMd {
	c := common{}
	return &repMd{
		dbConfig: c.ConvertDbConfigTo2000(c.GetMdDbConfig()),
	}
}

//获取zlCompany信息
func (r *repMd) GetZlCompany() (*object.ZlCompany, error) {
	comm := NewCommon()
	rows, err := comm.GetRowsBySQL2000(r.dbConfig, sqlGetMdZlCompany)
	if err != nil {
		errMsg := fmt.Sprintf("get zlcompany err: %s", err.Error())
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}

	defer func() {
		_ = rows.Close()
	}()

	var fCoId, fCoFunc int
	var fCoAb, fCoCode, fCoUserAb, fCoUserCode string
	var fCoAccStartDay time.Time

	rList := make([]*object.ZlCompany, 0)
	for rows.Next() {
		err = rows.Scan(&fCoId, &fCoAb, &fCoCode, &fCoUserAb, &fCoUserCode, &fCoFunc, &fCoAccStartDay)
		if err != nil {
			errMsg := fmt.Sprintf("read zlcompany data err: %s", err.Error())
			log.Error(errMsg)
			return nil, errors.New(errMsg)
		}
		rList = append(rList, &object.ZlCompany{
			FCoId:          fCoId,
			FCoAb:          fCoAb,
			FCoCode:        fCoCode,
			FCoUserAb:      fCoUserAb,
			FCoUserCode:    fCoUserCode,
			FCoFunc:        fCoFunc,
			FCoAccStartDay: fCoAccStartDay,
		})
	}
	if rows.Err() != nil {
		errMsg := fmt.Sprintf("read zlcompany data err: %s", rows.Err().Error())
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	if len(rList) < 1 {
		errMsg := fmt.Sprintf("get zlcompany err: %s", "zlcompany is empty")
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	return rList[0], nil
}

func (r *repMd) GetMdYyInfo() ([]*object.MdYyInfo, error) {
	comm := NewCommon()
	rows, err := comm.GetRowsBySQL2000(r.dbConfig, sqlGetMdYyInfo)
	if err != nil {
		errMsg := fmt.Sprintf("get tc info err: %s", err.Error())
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	defer func() {
		_ = rows.Close()
	}()
	var fYyr time.Time
	var fMdId int
	var fTc, fSr float32
	rList := make([]*object.MdYyInfo, 0)
	for rows.Next() {
		err = rows.Scan(&fMdId, &fYyr, &fTc, &fSr)
		if err != nil {
			errMsg := fmt.Sprintf("read tc info data err: %s", err.Error())
			log.Error(errMsg)
			return nil, errors.New(errMsg)
		}
		rList = append(rList, &object.MdYyInfo{
			FMdId:    fMdId,
			FYyr:     fYyr,
			FTc:      fTc,
			FSr:      fSr,
			FOprTime: time.Now(),
		})
	}
	if rows.Err() != nil {
		errMsg := fmt.Sprintf("read tc info data err: %s", rows.Err().Error())
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	return rList, nil
}

func (r *repMd) GetZxKcInfo(lst time.Time) ([]*object.ZxKc, time.Time, error) {
	comm := NewCommon()
	log.Debug(goToolCommon.GetDateTimeStrWithMillisecond(lst))
	rows, err := comm.GetRowsBySQL2000(r.dbConfig, sqlGetZxKcInfo,
		goToolCommon.GetDateTimeStrWithMillisecond(lst))
	if err != nil {
		errMsg := fmt.Sprintf("get zxkc err: %s", err.Error())
		log.Error(errMsg)
		return nil, NewCommon().GetDefaultOprTime(), errors.New(errMsg)
	}
	defer func() {
		_ = rows.Close()
	}()
	newLst := lst
	var fHpId int
	var fSl float32
	var fTime time.Time
	rList := make([]*object.ZxKc, 0)
	for rows.Next() {
		err = rows.Scan(&fHpId, &fSl, &fTime)
		if err != nil {
			errMsg := fmt.Sprintf("read zxkc err: %s", err.Error())
			log.Error(errMsg)
			return nil, NewCommon().GetDefaultOprTime(), errors.New(errMsg)
		}
		rList = append(rList, &object.ZxKc{
			FHpId:    fHpId,
			FSl:      fSl,
			FOprTime: fTime,
		})
		if fTime.After(newLst) {
			newLst = fTime
		}
	}
	if rows.Err() != nil {
		errMsg := fmt.Sprintf("read zxkc err: %s", rows.Err().Error())
		log.Error(errMsg)
		return nil, NewCommon().GetDefaultOprTime(), errors.New(errMsg)
	}
	return rList, newLst, nil
}
