package controllers

import (
	"OAM/models"

	"github.com/beego/beego/v2/core/logs"

	fn "OAM/util"
)

type AppInfoController struct {
	AuthController
}

//转账号列表页
func (c *AppInfoController) ToAppList() {
	projs := models.FindProjectForMap()
	c.Data["projs"] = projs
	c.Data["appTypes"] = models.GetDictValueAsArray(models.DICTKEY_APP_TYPE)
	c.TplName = "appinfo.html"
}

func (c *AppInfoController) AppPage() {
	where := make(map[string]interface{})
	ip := c.GetString("keyword")
	if ip != "" {
		if fn.SafeChars(ip) {
			where["keyword"] = ip
		} else {
			c.JsonParamError("搜索关键字带有不允许的特殊符号")
		}
	}
	pid, err := c.GetInt("projId")
	if err == nil {
		where["projId"] = pid
	}
	appType := c.GetString("appType")
	if appType != "" {
		where["appType"] = appType
	}
	hostId, err := c.GetInt("hostId")
	if err == nil {
		where["hostId"] = hostId
	}

	row, _ := c.GetInt("rows", models.DEFAULT_PAGE_SIZE)
	page, _ := c.GetInt("page", 1)
	pageData := models.FindAppInfoForPage(row, page, where)
	c.JsonOk(pageData)
}

func (c *AppInfoController) SaveApp() {
	c.justPost()
	var appVo models.AppInfoVO
	err := c.BindJSON(&appVo)
	//err := c.ParseForm(&appVo)
	if err != nil {
		logs.Error(err)
		c.JsonParamError("提交的数据错误")
		return
	}
	if errmsg := appVo.Valid(); errmsg != "" {
		c.JsonParamError(errmsg)
	}
	if appVo.AppId == 0 {
		err = models.SaveAppInfo(&appVo)
	} else {
		err = models.UpdateAppInfo(&appVo)
	}
	if err != nil {
		c.JsonFailed(err.Error())
	} else {
		c.JsonOk(appVo.AppId)
	}
}

func (c *AppInfoController) AppDetail() {
	id, err := c.GetInt("id", 0)
	if err != nil {
		c.toErrorPage("参数错误")
		return
	}
	c.Data["srvsoftwares"] = models.GetPresetServiceSoftwares()
	c.Data["devLangs"] = models.GetDictValueAsArray(models.DICTKEY_DEV_LANG)
	c.Data["appTypes"] = models.GetDictValueAsArray(models.DICTKEY_APP_TYPE)
	projs := models.FindProjectForMap()
	c.Data["projs"] = projs
	if id > 0 {
		appInfo := models.GetAppInfoVOById(id)
		if appInfo != nil {
			c.Data["app"] = appInfo
		} else {
			c.toErrorPage("应用不存在")
			return
		}
	} else {
		pid, _ := c.GetInt("projId", 0)
		if pid > 0 {
			c.Data["sepc_projid"] = pid
		}
		hid, _ := c.GetInt("hostId", 0)
		if hid > 0 {
			c.Data["sepc_hostid"] = hid
		}
	}
	c.TplName = "appinfo-edit.html"
}

func (c *AppInfoController) DelApp() {
	appId, err := c.GetInt("id")
	if err != nil || appId <= 0 {
		c.JsonParamError("参数错误")
		return
	}
	isOk := models.DeleteAppInfo(appId)
	c.JsonOk(isOk)
}
