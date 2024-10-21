package controllers

type ErrorController struct {
	BaseController
}

//自定义404页面
func (c *ErrorController) Error404() {
	if c.IsAjax() {
		c.Abort("404")
	}
	c.Data["content"] = "访问的页面不存在"
	c.TplName = "error.html"
}

func (c *ErrorController) Error403() {
	if c.IsAjax() {
		c.JsonForbidden()
	}
	c.Data["content"] = "无权限访问该页面"
	c.TplName = "error.html"
}

func (c *ErrorController) Error500() {
	if c.IsAjax() {
		c.JsonError()
	}

	c.Data["content"] = "internal server error"
	c.TplName = "error.html"
}
