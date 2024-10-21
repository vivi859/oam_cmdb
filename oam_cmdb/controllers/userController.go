package controllers

import (
	"OAM/models"
	fn "OAM/util"
	"strconv"
	"time"

	"github.com/beego/beego/v2/core/logs"
)

type UserController struct {
	AuthController
}

func (c *UserController) Profile() {
	user := c.getLoginUser()
	c.Data["user"] = models.GetUserById(user.UserId)
	c.TplName = "user-profile.html"
}

func (c *UserController) Detail() {
	id, err := c.GetInt("id", 0)
	if err != nil || id == 0 {
		c.JsonParamError("请选择用户")
	}
	user := models.GetUserById(id)
	if user != nil {
		user.Passwd = ""
		c.JsonOk(user)
	} else {
		c.JsonParamError("用户不存在")
	}
}

func (c *UserController) ToUserList() {
	c.Data["roleList"] = models.FindRoles()
	c.TplName = "user.html"
}

// 查找用户
func (c *UserController) Find() {
	condiMap := make(map[string]interface{})
	status, err := c.GetInt("status", -1)
	if err == nil && status > -1 {
		condiMap["status"] = status
	}
	uname := c.GetString("keyword")
	if uname != "" {
		condiMap["keyword"] = uname
	}

	role := c.GetString("role")
	if role != "" {
		condiMap["roleCode"] = role
	}
	users := models.FindUserByCondi(condiMap)
	c.JsonOk(users)
}

func (c *UserController) Save() {
	c.justPost()
	u := models.UserInfo{}
	if err := c.ParseForm(&u); err != nil {
		c.JsonParamError(err.Error())
	}
	curUser := c.getLoginUser()
	if !curUser.IsSupervisor() && u.RoleCode == models.ROLE_ROOT {
		c.JsonFailed("对不起!您没有建立超级管理员账号的权限")
		return
	}
	var id int64
	var err error
	if u.UserId == 0 {
		u.CreateBy = curUser.UserName
		u.UpdateBy = curUser.UserName
		id, err = models.SaveUser(&u)
	} else {
		u.UpdateBy = curUser.UserName
		id, err = models.UpdateUser(&u)
	}

	if err == nil {
		c.JsonOk(id)
	} else {
		c.JsonFailed(err.Error())
	}
}

// 登录用户修改密码
func (c *UserController) UpdateLoginUserPasswd() {
	c.justPost()
	oldPwd := c.GetString("oldPasswd")
	pwd := c.GetString("passwd")
	rePwd := c.GetString("rePasswd")
	if oldPwd == "" || pwd == "" || rePwd == "" {
		c.JsonParamError("参数错误")
	}
	if pwd != rePwd {
		c.JsonParamError("两次密码不相同")
	}
	curUser := c.getLoginUser()
	user := models.GetUserById(curUser.UserId)
	if user == nil {
		c.JsonParamError("用户不存在")
		return
	}
	if user.UserStatus == models.Diabled {
		c.DestroySession()
		c.JsonFailed("用户已禁用,不允许操作")
	}
	if !user.CheckPasswd(oldPwd) {
		c.JsonFailed("旧密码不正确")
	}
	user.UpdateBy = curUser.UserName
	user.Passwd = pwd
	user.RePasswd = rePwd
	err := models.UpdateUserPasswd(user)
	if err == nil {
		logs.Info("%s修改了个人密码", curUser.UserName)
		c.JsonOk(true)
	} else {
		c.JsonFailed(err.Error())
	}
}

// 重置用户密码
func (c *UserController) ResetUserPasswd() {
	id, err := c.GetInt("id", 0)
	if err != nil || id == 0 {
		c.JsonParamError("请选择用户")
	}
	curUser := c.getLoginUser()
	if id == curUser.UserId {
		c.JsonParamError("修改当前登录用户密码,请使用修改个人密码功能")
	}

	newpwd := generateNewPwd()
	u := models.UserInfo{UserId: id, Passwd: newpwd, RePasswd: newpwd, UpdateBy: curUser.UserName}
	err = models.UpdateUserPasswd(&u)
	logs.Info("%s重置了[userid=%d]的密码,结果:%v", curUser.UserName, id, err == nil)
	if err == nil {
		c.JsonOk(newpwd)
	} else {
		c.JsonFailed(err.Error())
	}
}

// 生成一个密码,规则:年份 + 随机符号 + 随机字母
func generateNewPwd() string {
	now := time.Now()
	newpwd := strconv.Itoa(now.Year())
	newpwd = newpwd + fn.Random(1, "@#*!%") + fn.RandomStr(3)
	return newpwd
}
