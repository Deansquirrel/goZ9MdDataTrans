package repository

import "github.com/Deansquirrel/goToolMSSql2000"

type repMd struct {
	dbConfig *goToolMSSql2000.MSSqlConfig
}

func NewRepMd() *repMd {
	c := common{}
	return &repMd{
		dbConfig: c.ConvertDbConfigTo2000(c.GetMdDbConfig()),
	}
}
