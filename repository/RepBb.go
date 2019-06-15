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
	sqlGetBBZlCompany = "" +
		"SELECT [coid],[coab],[cocode],[couserab],[cousercode]" +
		",[cofunc],[coaccstartday] " +
		"FROM [zlcompany]"
	sqlRestoreMdYyInfo = "" +
		"IF EXISTS (SELECT * FROM [mdyyinfo] WHERE [mdid] = ? AND [yyr] = ?) " +
		"	BEGIN " +
		"		UPDATE [mdyyinfo] " +
		"		SET [tc]=?,[sr]=?,[recorddate]=? " +
		"		WHERE [mdid] = ? AND [yyr] = ? " +
		"	END " +
		"ELSE " +
		"	BEGIN " +
		"		INSERT INTO [mdyyinfo]([mdid],[yyr],[tc],[sr],[recorddate]) " +
		"		VALUES(?,?,?,?,?) " +
		"	END"
	sqlRestoreZxKc = "" +
		"IF EXISTS(SELECT * FROM [zxkc] WHERE [mdid]=? AND [hpid]=?) " +
		"	BEGIN " +
		"		UPDATE [zxkc] " +
		"		SET [sl]=?,[lastupdate]=? " +
		"		WHERE [mdid]=? AND [hpid]=? " +
		"	END " +
		"ELSE " +
		"	BEGIN " +
		"		INSERT INTO [zxkc]([mdid],[hpid],[sl],[lastupdate]) " +
		"		VALUES (?,?,?,?) " +
		"	END"
)

type repBb struct {
	dbConfig *goToolMSSql2000.MSSqlConfig
}

func NewRepBb() *repBb {
	c := common{}
	return &repBb{
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
	comm := NewCommon()
	return comm.SetRowsBySQL2000(r.dbConfig, sqlRestoreMdYyInfo,
		opr.FMdId,
		goToolCommon.GetDateStr(opr.FYyr),
		fmt.Sprintf("%v", opr.FTc),
		fmt.Sprintf("%v", opr.FSr),
		goToolCommon.GetDateTimeStrWithMillisecond(opr.FOprTime),
		opr.FMdId,
		goToolCommon.GetDateStr(opr.FYyr),
		opr.FMdId,
		goToolCommon.GetDateStr(opr.FYyr),
		fmt.Sprintf("%v", opr.FTc),
		fmt.Sprintf("%v", opr.FSr),
		goToolCommon.GetDateTimeStrWithMillisecond(opr.FOprTime),
	)
}

func (r *repBb) RestoreZxKc(opr *object.ZxKcOpr) error {
	comm := NewCommon()
	return comm.SetRowsBySQL2000(r.dbConfig, sqlRestoreZxKc,
		opr.FMdId,
		opr.FHpId,
		fmt.Sprintf("%v", opr.FSl),
		goToolCommon.GetDateTimeStrWithMillisecond(opr.FOprTime),
		opr.FMdId,
		opr.FHpId,
		opr.FMdId,
		opr.FHpId,
		fmt.Sprintf("%v", opr.FSl),
		goToolCommon.GetDateTimeStrWithMillisecond(opr.FOprTime))
}
