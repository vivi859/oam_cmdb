package controllers

import (
	"OAM/models"

	"github.com/beego/beego/v2/core/logs"
)

type RoleController struct {
	AuthController
}

//角色列表页
func (c *RoleController) RoleList() {
	roles := models.FindRoles()
	c.JsonOk(roles)
}

func (c *RoleController) Detail() {
	id := c.GetString("id")
	if id == "" {
		c.JsonParamError("请选择角色")
	}
	r := models.GetRoleById(id)
	if r != nil {
		c.JsonOk(r)
	} else {
		c.JsonParamError("角色不存在")
	}
}

func (c *RoleController) Save() {
	c.justPost()
	r := models.Role{}
	if err := c.ParseForm(&r); err != nil {
		c.JsonParamError(err.Error())
	}
	if errmsg := r.Valid(); errmsg != "" {
		c.JsonParamError(errmsg)
	}
	var id int64
	var err error
	if len(r.RoleCode) <= 1 {
		id, err = models.SaveRole(&r)
		logs.Info("新增角色：%s", r)
	} else {
		err = models.UpdateRole(&r)
		logs.Info("修改角色：%s", r)
	}
	if err == nil {
		c.JsonOk(id)
	} else {
		c.JsonFailed(err.Error())
	}
}

func (c *RoleController) DeleteRole() {
	id := c.GetString("id")
	if id == "" {
		c.JsonParamError("请选择角色")
		return
	}
	err := models.DeleteRole(id)
	logs.Info("删除角色：%d, err:%s", id, err)
	if err == nil {
		c.JsonOk(id)
	} else {
		c.JsonFailed(err.Error())
	}
}
