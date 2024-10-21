package controllers

import (
	"OAM/conf"
	"OAM/models"
	fn "OAM/util"
	"encoding/base64"
	"net/url"
	"regexp"
	"strings"

	"github.com/beego/beego/v2/core/logs"
)

type LoginController struct {
	BaseController
}

func (c *LoginController) Get() {
	c.TplName = "login.html"
}

func (c *LoginController) Post() {
	//参数验证
	uname := c.GetString("uname")
	upwd := c.GetString("upwd")
	if uname == "" || upwd == "" {
		c.JsonParamError("用户名或密码错误")
	}
	uname = strings.TrimSpace(uname)
	upwd = strings.TrimSpace(upwd)
	if isValid, err := regexp.MatchString(models.USER_NAME_PATTERN, uname); err != nil && !isValid {
		logs.Warn("登录失败,用户名%s格式不正确", uname)
		c.JsonParamError("用户名不正确")
	}
	if isValid, err := regexp.MatchString(models.USER_PWD_PATTERN, upwd); err != nil && !isValid {
		logs.Warn("登录失败,密码%s格式不正确", upwd)
		c.JsonParamError("密码不正确")
		return
	}

	var failedQty = getLoginFailedQty(uname)
	if failedQty > 5 {
		logs.Error("账号登录次数超限锁定,用户名:%s,ip:%s,尝试次数:%d", uname, c.Ctx.Input.IP(), failedQty)
		c.JsonFailed("账号已被锁定,请稍候再试")
	}
	//解密密码
	plainPwd, err := pwdDecode(upwd)
	if err != nil {
		cacheFailedQty(uname, failedQty == 0)
		c.JsonParamError("用户名或密码错误")
	}

	// 密码和用户状态验证
	user := models.GetUserByName(uname)
	if user != nil {
		if user.UserStatus == 0 {
			c.JsonParamError("用户已被禁用")
		}
		if user.CheckPasswd(plainPwd) {
			if user.RoleCode == "" {
				c.JsonParamError("用户角色未配置")
				return
			}
			loginUser := models.CreateLoginUser(*user)
			c.SessionRegenerateID()
			c.SetSession(LOGIN_USER_KEY, loginUser)
			clearFailedQty(uname)
			userMap := make(map[string]interface{})
			userMap["userName"] = loginUser.UserName
			userMap["realName"] = loginUser.RealName
			userMap["roleCode"] = loginUser.RoleCode
			c.JsonOk(userMap)
			return
		}
	}
	cacheFailedQty(uname, failedQty == 0)
	c.JsonParamError("用户名或密码错误")
}

func (c *LoginController) Logout() {
	c.DestroySession()
	c.Redirect("/login", 302)
}

// 返回客户端公钥
func (c *BaseController) Pubkey() {
	c.justPost()
	key, err := fn.PublicKeyToBase64(conf.GlobalCfg.RSA_DEFAULT_PUBLIC_KEY)
	if err != nil {
		c.JsonFailed("公钥加载失败")
		return
	}
	c.JsonOk(key)
}

func (c *BaseController) Prikey() {
	c.justPost()
	key, err := fn.PrivateKeyToBase64(conf.GlobalCfg.RSA_DEFAULT_PRIVATE_KEY)
	if err != nil {
		c.JsonFailed("客户端私钥加载失败")
		return
	}
	c.JsonOk(key)
}

func getCache() fn.MyCache {
	return fn.GetCache(fn.CACHE_LOGIN_FAILED)
}

// 获取用户已登录失败次数
func getLoginFailedQty(userName string) int {
	return getCache().GetInt(userName)
}

// 缓存登录错误次数
func cacheFailedQty(userName string, isFirst bool) {
	if isFirst {
		getCache().Put(userName, 1)
	} else {
		getCache().Incr(userName)
	}
}

func clearFailedQty(userName string) {
	getCache().Delete(userName)
}

// 解码密码
func pwdDecode(pwd string) (string, error) {
	tmp, err := base64.StdEncoding.DecodeString(pwd)
	if err == nil {
		tmpStr, err := url.QueryUnescape(string(tmp))
		if err == nil {
			tmpStr = fn.Reverse(tmpStr[0 : len(tmpStr)-5])
			tmp, err := base64.StdEncoding.DecodeString(tmpStr)
			if err == nil {
				tmpStr, _ = url.QueryUnescape(string(tmp))
				return url.QueryUnescape(tmpStr)
			}
		}
	}
	return "", err
}
