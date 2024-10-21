package controllers

import (
	"OAM/models"
	"fmt"
	"os"
	"strings"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
)

type InstallController struct {
	BaseController
}

type CfgParam struct {
	DbType     string `valid:"Required;" label:"请选择数据库类型"`
	DbHost     string `valid:"Required;" label:"请输入数据库地址"`
	DbName     string
	DbUsername string `valid:"Required;" label:"请输入数据库账号"`
	Dbpasswd   string `valid:"Required;" label:"请输入数据库密码"`
}

func (cfg *CfgParam) Valid() string {
	valid := validation.Validation{}
	_, err := valid.Valid(cfg)
	if err != nil {
		return err.Error()
	}

	return models.ToErrMsg(valid)
}
func (c *InstallController) install() {
	c.justPost()
	var cfg CfgParam
	err := c.BindJSON(&cfg)
	if err != nil {
		logs.Error(err)
		c.JsonParamError("提交的数据错误")
		return
	}
	if errmsg := cfg.Valid(); errmsg != "" {
		c.JsonParamError(errmsg)
	}

	if len(cfg.DbName) == 0 {
		cfg.DbName = "oam"
	}

	if cfg.DbType == "mysql" {
		dbUrlTpl := "%s:%s@tcp(%s)/%s?charset=utf8&loc=Local"
		var dbUrl = fmt.Sprintf(dbUrlTpl, cfg.DbUsername, cfg.Dbpasswd, cfg.DbHost, cfg.DbName)
		orm.RegisterDataBase("default", "mysql", dbUrl)
		bs, err := os.ReadFile("data/oam_mysql.sql")
		if err != nil {
			logs.Error(err)
			c.JsonFailed("初始数据库失败" + err.Error())
		}
		sqlText := string(bs)
		sqlScripts := strings.Split(sqlText, ";")

		for _, sql := range sqlScripts {
			_, err := orm.NewOrm().Raw(sql).Exec()
			if err != nil {
				logs.Error(err)
				c.JsonFailed("初始数据库失败" + err.Error())
				break
			}
		}
	} else if cfg.DbType == "sqlite" {
		orm.RegisterDriver("sqlite", orm.DRSqlite)
		orm.RegisterDataBase("default", "sqlite3", "data/oam.db")
	}
}
