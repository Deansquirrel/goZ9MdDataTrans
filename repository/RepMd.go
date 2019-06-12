package repository

import (
	"errors"
	"fmt"
	"github.com/Deansquirrel/goToolMSSql2000"
	"github.com/Deansquirrel/goZ9MdDataTrans/object"
	"time"
)

import log "github.com/Deansquirrel/goToolLog"

const (
	sqlGetZlCompany = "" +
		"SELECT [coid],[coab],[cocode],[couserab],[cousercode]" +
		",[cofunc],[coaccstartday] " +
		"FROM [zlcompany]"

	sqlGetMdYyInfo = "" +
		"SELECT [ckmdid],[ckyyr],sum(case when [ckcxbj]=1 then -1 else 1 end) as [num],sum([ckcjje]) as [srmy] " +
		"FROM [z3xsckt] WITH(NOLOCK) " +
		"GROUP BY [ckmdid],[ckyyr]"
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
	rows, err := comm.GetRowsBySQL2000(r.dbConfig, sqlGetZlCompany)
	if err != nil {
		errMsg := fmt.Sprintf("get zlcompany err: %s", err.Error())
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}

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
	var fYyr time.Time
	var fMdId, fTc int
	var fSr float32
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
