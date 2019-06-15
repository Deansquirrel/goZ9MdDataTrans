package object

import (
	"encoding/json"
	"fmt"
	"github.com/Deansquirrel/goZ9MdDataTrans/global"
	"strings"
)

import log "github.com/Deansquirrel/goToolLog"

//系统配置（Server|Client）
type SystemConfig struct {
	Total        systemConfigTotal        `toml:"total"`
	OnLineConfig systemConfigOnLineConfig `toml:"onLineConfig"`
	MdDB         systemConfigMdDB         `toml:"mdDb"`
	RunMode      systemConfigRunMode      `toml:"runMode"`
	Service      systemConfigService      `toml:"service"`
}

func (sc *SystemConfig) FormatConfig() {
	sc.Total.FormatConfig()
	sc.OnLineConfig.FormatConfig()
	sc.MdDB.FormatConfig()
	sc.Service.FormatConfig()
}

func (sc *SystemConfig) ToString() string {
	d, err := json.Marshal(sc)
	if err != nil {
		log.Warn(fmt.Sprintf("SystemConfig转换为字符串时遇到错误：%s", err.Error()))
		return ""
	}
	return string(d)
}

//通用配置
type systemConfigTotal struct {
	StdOut   bool   `toml:"stdOut"`
	LogLevel string `toml:"logLevel"`
}

func (t *systemConfigTotal) FormatConfig() {
	//去除首尾空格
	t.LogLevel = strings.Trim(t.LogLevel, " ")
	//设置默认日志级别
	if t.LogLevel == "" {
		t.LogLevel = "warn"
	}
	//设置字符串转换为小写
	t.LogLevel = strings.ToLower(t.LogLevel)
	t.LogLevel = t.checkLogLevel(t.LogLevel)
}

//校验SysConfig中iris日志级别设置
func (t *systemConfigTotal) checkLogLevel(level string) string {
	switch level {
	case "debug", "info", "warn", "error":
		return level
	default:
		return "warn"
	}
}

//线上连接信息
type systemConfigOnLineConfig struct {
	Address string `toml:"address"`
}

func (c *systemConfigOnLineConfig) FormatConfig() {
	c.Address = strings.Trim(c.Address, " ")
}

//RunMode
type systemConfigRunMode struct {
	Mode RunMode `toml:"mode"`
}

func (c *systemConfigRunMode) FormatConfig() {
	if string(c.Mode) == "" {
		c.Mode = global.DefaultRunMode
	}
}

//门店库配置库
type systemConfigMdDB struct {
	Server string `toml:"server"`
	Port   int    `toml:"port"`
	DbName string `toml:"dbName"`
	User   string `toml:"user"`
	Pwd    string `toml:"pwd"`
}

func (c *systemConfigMdDB) FormatConfig() {
	c.Server = strings.Trim(c.Server, " ")
	if c.Port == 0 {
		c.Port = 1433
	}
	c.DbName = strings.Trim(c.DbName, " ")
	c.User = strings.Trim(c.User, " ")
	c.Pwd = strings.Trim(c.Pwd, " ")
}

//服务配置
type systemConfigService struct {
	Name        string `toml:"name"`
	DisplayName string `toml:"displayName"`
	Description string `toml:"description"`
}

//格式化
func (sc *systemConfigService) FormatConfig() {
	sc.Name = strings.Trim(sc.Name, " ")
	sc.DisplayName = strings.Trim(sc.DisplayName, " ")
	sc.Description = strings.Trim(sc.Description, " ")
	if sc.Name == "" {
		sc.Name = "Z9MdDataTrans"
	}
	if sc.DisplayName == "" {
		sc.DisplayName = "Z9MdDataTrans"
	}
	if sc.Description == "" {
		sc.Description = sc.Name
	}
}
