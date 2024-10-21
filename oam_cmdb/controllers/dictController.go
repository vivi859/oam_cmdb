package controllers

import (
	"OAM/models"
)

type DictController struct {
	AuthController
}

// func (c *DictController) ToDictList() {
// 	c.TplName = "dict.html"
// }

func (c *DictController) DictList() {
	di := models.FindAllDictItem()
	c.JsonOk(di)
}

func (c *DictController) SaveDict() {
	c.justPost()
	dict := models.DictItem{}
	if err := c.BindForm(&dict); err != nil {
		c.JsonParamError(err.Error())
	}
	errmsg := dict.Valid()
	if errmsg != "" {
		c.JsonParamError(errmsg)
	}

	curUser := c.getLoginUser()
	//dict.ItemGroup = 0
	dict.UpdateBy = curUser.UserName
	models.UpdateDictItem(&dict)
	c.JsonOk(true)
}
