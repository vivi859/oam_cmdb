package controllers

import (
	"OAM/util"
)

type MainController struct {
	AuthController
}

func (c *MainController) Get() {
	c.TplName = "index.html"
}

// 直接转向模板目录html页面
func (c *MainController) Forward() {
	page := c.GetString(":page")
	if len(page) > 2 {
		c.TplName = util.ToFirstLetterLower(page) + ".html"
	} else {
		c.Abort("404")
	}
}
