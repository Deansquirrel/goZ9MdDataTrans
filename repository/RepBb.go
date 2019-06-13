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
	sqlGetBBZlCompany = "" +
		"SELECT [coid],[coab],[cocode],[couserab],[cousercode]" +
		",[cofunc],[coaccstartday] " +
		"FROM [zlcompany]"
)

type repBb struct {
	dbConfig *goToolMSSql2000.MSSqlConfig
}

func NewRepBb() *repMd {
	c := common{}
	return &repMd{
		dbConfig: c.ConvertDbConfigTo2000(c.GetMdDbConfig()),
	}
}

//获取zlCompany信息
func (r *repBb) GetZlCompany() (*object.ZlCompany, error) {
	comm := NewCommon()
	rows, err := comm.GetRowsBySQL2000(r.dbConfig, sqlGetBBZlCompany)
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

func (r *repBb) RestoreMdYyInfo(opr *object.MdYyInfoOpr) error {
	//TODO
	return nil
}

func (r *repBb) RestoreZxKc(opr *object.ZxKcOpr) error {
	//TODO
	return nil
}