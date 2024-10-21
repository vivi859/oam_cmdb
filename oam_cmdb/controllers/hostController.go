package controllers

import (
	"OAM/models"
	fn "OAM/util"
	"regexp"

	"github.com/beego/beego/v2/core/logs"
)

type HostController struct {
	AuthController
}

//转主机列表页
func (c *HostController) ToHostList() {
	c.Data["osnames"] = models.GetPresetOS()
	projs := models.FindProjectForMap()
	/*curUser := c.getLoginUser()
	 if len(projs) > 0 && !curUser.IsSupervisor() {
		myProjIds := models.FindAuthorizedProjectIds(curUser.UserId)
		maps.DeleteFunc(projs, func(k int, v string) bool {
			return !slices.Contains[int](myProjIds, k)
		})
	} */
	c.Data["projs"] = projs
	c.TplName = "host.html"
}

func (c *HostController) HostPage() {
	where := make(map[string]interface{})
	ip := c.GetString("keyword")
	if ip != "" {
		if safe, _ := regexp.MatchString("[~#\\^\\$><%!*=`]", ip); !safe {
			where["keyword"] = ip
		} else {
			c.JsonParamError("查询条件不能带有特殊符号~#^$><%!*=`")
			return
		}
	}
	os := c.GetString("os")
	if os != "" {
		where["os"] = os
	}
	htype, _ := c.GetInt("htype", -1)
	if htype >= 0 {
		where["htype"] = htype
	}
	hid, _ := c.GetInt("hostId", 0)
	if hid > 0 {
		where["hostId"] = hid
	}
	//只查找正常状态主机
	justNormal, _ := c.GetBool("justNormal", false)
	if justNormal {
		where["justNormal"] = 0
	}
	projId, _ := c.GetInt("projId", 0)
	if projId > 0 {
		where["projId"] = projId
	}
	row, _ := c.GetInt("rows", models.DEFAULT_PAGE_SIZE)
	page, _ := c.GetInt("page", 1)
	pageData := models.FindHostForPage(row, page, where)
	c.JsonOk(pageData)
}

func (c *HostController) BaseHostPage() {
	where := make(map[string]interface{})
	ip := c.GetString("keyword")
	if ip != "" {
		if safe, _ := regexp.MatchString("[~#\\^\\$><%!*=`]", ip); !safe {
			where["keyword"] = ip
		} else {
			c.JsonParamError("查询条件不能带有特殊符号~#^$><%!*=`")
			return
		}
	}
	hid, _ := c.GetInt("hostId", 0)
	if hid > 0 {
		where["hostId"] = hid
	}
	projId, _ := c.GetInt("projId", 0)
	if projId > 0 {
		where["projId"] = projId
	}
	//只查找正常状态主机
	justNormal, _ := c.GetBool("justNormal", false)
	if justNormal {
		where["justNormal"] = 0
	}

	row, _ := c.GetInt("rows", models.DEFAULT_PAGE_SIZE)
	page, _ := c.GetInt("page", 1)
	pageData := models.FindBaseHostForPage(row, page, where)
	c.JsonOk(pageData)
}

func (c *HostController) SaveHost() {
	c.justPost()
	var hostvo models.HostVO
	err := c.BindJSON(&hostvo)
	if err != nil {
		logs.Error(err)
		c.JsonParamError("提交的数据错误")
		return
	}
	if errmsg := hostvo.Valid(); errmsg != "" {
		c.JsonParamError(errmsg)
	}
	//curUser, _ := c.getLoginUser()
	if hostvo.HostId == 0 {
		err = models.SaveHost(&hostvo)
	} else {
		err = models.UpdateHost(&hostvo)
	}
	if err != nil {
		c.JsonFailed(err.Error())
	} else {
		c.JsonOk(hostvo.HostId)
	}
}

func (c *HostController) HostDetail() {
	id, err := c.GetInt("id", 0)
	if err != nil {
		//c.JsonParamError("参数错误")
		c.toErrorPage("参数错误")
		return
	}

	projs := models.FindProjectForMap()
	c.Data["projs"] = projs

	if id == 0 {
		host := models.HostVO{}
		projId, err := c.GetInt("projId") //指定了项目
		if err == nil && projId > 0 {
			host.ProjIds = fn.IntList{projId}
		}
		c.Data["host"] = host
	} else {
		host := models.GetHostVOById(id)
		if host != nil {
			if len(host.Accounts) > 0 {
				for _, a := range host.Accounts {
					if a.FieldPwd != "" {
						a.FieldPwd, err = fn.AesDecryptStr(a.FieldPwd)
						if err != nil {
							c.toErrorPage("账号解密异常")
							return
						}
					}
				}
			}
			c.Data["host"] = host
		} else {
			c.toErrorPage("主机不存在")
			return
		}
	}
	c.Data["osnames"] = models.GetPresetOS()
	c.TplName = "host-edit.html"
}

func (c *HostController) DelHost() {
	hostId, err := c.GetInt("id")
	if err != nil || hostId <= 0 {
		c.JsonParamError("参数错误")
		return
	}
	var action int
	action, err = c.GetInt("action", 1)
	if err != nil || action < 0 {
		c.JsonParamError("参数错误")
		return
	}
	isOk := models.DeleteHost(hostId, action)
	logs.Info("[%s]删除主机id:%d", c.getLoginUser().UserName, hostId)
	c.JsonOk(isOk)
}

func (c *HostController) RecoverHost() {
	hostId, err := c.GetInt("id")
	if err != nil || hostId <= 0 {
		c.JsonParamError("参数错误")
		return
	}
	host := models.GetHostById(hostId)
	if !host.IsDeleted {
		c.JsonParamError("当前主机未废除")
		return
	}
	isOk := models.RecoverHost(hostId)
	logs.Info("[%s]恢复废除主机id:%d", c.getLoginUser().UserName, hostId)
	c.JsonOk(isOk)
}

func (c *HostController) QueryAppForHost() {
	hostId, err := c.GetInt("id")
	if err != nil || hostId <= 0 {
		c.Ctx.WriteString("参数错误")
		return
	}
	names := models.FindAppNameByHostId(hostId)
	var html string
	if len(names) == 0 {
		html = "<p>未部署应用</p>"
	} else {
		html = "<ul>"
		for _, n := range names {
			html = html + "<li>" + n + "</li>"
		}
		html += "</ul>"
	}
	c.Ctx.WriteString(html)
}
