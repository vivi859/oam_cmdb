package controllers

import (
	"OAM/models"

	fn "OAM/util"

	"github.com/beego/beego/v2/core/logs"
	jsoniter "github.com/json-iterator/go"
)

type AccountController struct {
	AuthController
}

//转账号列表页
func (c *AccountController) ToAccountList() {
	projs := models.FindProjectForMap()
	c.Data["projs"] = projs
	accountTypes := models.FindAccountTypesForMap()
	c.Data["accountTypes"] = accountTypes
	c.TplName = "account.html"
}

func (c *AccountController) TypeList() {
	dts := models.FindAllAccountTypes()
	c.JsonOk(dts)
}

//查询类型及字段
func (c *AccountController) GetType() {
	typeId, err := c.GetInt("typeid")
	if err != nil {
		c.JsonParamError("")
	}
	dt := models.GetAccountTypeById(typeId, true)
	if dt == nil {
		c.JsonParamError("未找到该类型")
	}
	isUsed := models.HasAccountInType(typeId)
	result := make(map[string]interface{})
	result["type"] = dt
	result["isUsed"] = isUsed
	c.JsonOk(result)
}

//保存类型
func (c *AccountController) SaveType() {
	c.justPost()
	var dtype models.AccountType
	err := c.BindJSON(&dtype)
	if err != nil {
		logs.Error(err)
		c.JsonParamError("提交的数据错误")
		return
	}
	//logs.Debug(jsoniter.MarshalToString(dtype))
	errmsg := dtype.Valid()
	if errmsg != "" {
		c.JsonParamError(errmsg)
	}
	for _, ele := range dtype.Fields {
		if ele.FieldId < 0 {
			ele.FieldId = 0
		}
	}
	err = models.SaveOrUpdateAccountType(&dtype)
	if err != nil {
		c.JsonFailed(err.Error())
	}
	c.JsonOk(dtype.TypeId)
}

func (c *AccountController) DelType() {
	typeId, err := c.GetInt("id")
	if err != nil {
		c.JsonParamError("参数错误")
	}
	err = models.DelAccountType(typeId)
	if err != nil {
		c.JsonFailed(err.Error())
	}
	c.JsonOk(true)
}

//获取账号类型的动态属性
func (c *AccountController) GetTypeFields() {
	typeId, err := c.GetInt("id")
	if err != nil {
		c.JsonParamError("参数错误")
	}
	fields := models.FindFieldsByTypeId(typeId)
	c.JsonOk(fields)
}

func (c *AccountController) GetAccountDetail() {
	actId, err := c.GetInt("id")
	if err != nil {
		c.JsonParamError("参数错误")
	}
	account := models.GetAccountById(actId)
	c.JsonOk(account)
}

//查询项目下账号
func (c *AccountController) ProjAccounts() {
	projId, err := c.GetInt("projId")
	if err != nil {
		c.JsonParamError("请选择项目")
	}
	accounts := models.FindAccountsByProjId(projId)
	err = models.CompletionVOProps(false, true, accounts)
	if err != nil {
		c.JsonFailed(err.Error())
	} else {
		c.JsonOk(accounts)
	}
}

//账号分页数据
func (c *AccountController) AccountPage() {
	projId, err := c.GetInt("projId")
	where := make(map[string]interface{})
	if err == nil && projId > 0 {
		where["projId"] = projId
	}
	typeId, err := c.GetInt("typeId")
	if err == nil && typeId >= 0 {
		where["typeId"] = typeId
	}
	name := c.GetString("accountName")
	if name != "" {
		where["accountName"] = fn.Trim(name)
	}

	notDeleted := c.GetString("notDeleted")
	if notDeleted != "" {
		where["notDeleted"] = notDeleted
	}
	baseQuery, _ := c.GetBool("base", false)
	row, _ := c.GetInt("rows", models.DEFAULT_PAGE_SIZE)
	page, _ := c.GetInt("page", 1)
	pageData := models.FindAccountForPage(row, page, where)
	err = models.CompletionVOProps(true, !baseQuery, pageData.Rows)
	//没有密码查看权限,清掉密码
	if pageData.TotalRow > 0 {
		curUser := c.getLoginUser()
		hasCopyPwdPerm := curUser.HasPermByFunCode("ACCOUNT_MGT:COPYPWD")
		if !hasCopyPwdPerm {
			for _, r := range pageData.Rows {
				r.FieldPwd = ""
			}
		}
	}
	if err == nil {
		c.JsonOk(pageData)
	} else {
		c.JsonFailed(err.Error())
	}
}

//转账号新增页
func (c *AccountController) ToAddAccount() {
	projId, err := c.GetInt("projId") //指定了项目
	if err == nil && projId > 0 {
		c.Data["projId"] = projId
	}
	projs := models.FindProjectForMap()
	c.Data["projs"] = projs
	accountTypes := models.FindAccountTypesForMap()
	c.Data["accountTypes"] = accountTypes
	c.TplName = "account-form.html"
}

func (c *AccountController) ToEditAccount() {
	toAccount(c, "edit")
}

func (c *AccountController) ToViewAccount() {
	toAccount(c, "view")
}

//转账号编辑页
func toAccount(c *AccountController, action string) {
	accountId, err := c.GetInt("id")
	if err != nil || accountId <= 0 {
		c.toErrorPage("错误!请选择账号")
		return
	}
	c.Data["act"] = action
	account := models.GetAccountById(accountId)
	if account == nil {
		c.toErrorPage("错误!账号不存在")
		return
	}

	if account.FieldPwd != "" {
		account.FieldPwd, err = fn.AesDecryptStr(account.FieldPwd)
		if err != nil {
			c.toErrorPage("错误!密码解密异常")
			return
		}
		if action == "view" {
			account.FieldPwd, err = fn.RSAEncryptBase64Str(account.FieldPwd, "")
			if err != nil {
				logs.Error("密码rsa加密异常,账号%s,%s", account.AccountName, account.FieldPwd)
				c.toErrorPage("错误!密码加密异常")
				return
			}
		}
	}
	c.Data["account"] = account
	curUser := c.getLoginUser()
	hasCopyPwdPerm := curUser.HasPermByFunCode("ACCOUNT_MGT:COPYPWD")
	//账号类型
	if account.TypeId > 0 {
		atype := models.GetAccountTypeById(account.TypeId, true)
		c.Data["type"] = atype
		//组装动态字段
		if account.FieldOther != "" {
			otherMap := make(map[string]string)
			jsoniter.UnmarshalFromString(account.FieldOther, &otherMap)
			for _, f := range atype.Fields {
				v, ok := otherMap[f.FieldKey]
				if ok {
					if f.ValueType == 0 && f.IsCiphertext && v != "" {
						f.FieldValue, err = fn.AesDecryptStr(v)
						if err != nil {
							c.toErrorPage("错误!密码解密异常")
							return
						}
						if action == "view" {
							if hasCopyPwdPerm {
								f.FieldValue, err = fn.RSAEncryptBase64Str(f.FieldValue, "")
								if err != nil {
									logs.Error("rsa加密字段异常,账号%s,fieldName:%s", account.AccountName, f.FieldName)
									c.toErrorPage("错误!加密异常")
									return
								}
							} else {
								f.FieldValue = ""
							}
						}
					} else {
						f.FieldValue = v
					}
				}
			}
		}

	}

	if account.HostId > 0 {
		host := models.GetHostById(account.HostId)
		if host != nil {
			hostName := host.HostName
			if host.PublicIp != "" {
				hostName = hostName + "(" + host.PublicIp + ")"
			} else {
				if host.InternalIp != "" {
					hostName = hostName + "(" + host.InternalIp + ")"
				}
			}
			c.Data["hostName"] = hostName
		}
	}

	if action == "view" {
		//关联账号
		if len(account.RelAccountIds) > 2 {
			relAccounts := models.FindAccountByIds(account.RelAccountIds.ToIntSlice(), true)
			relAccountNames := fn.SliceConvert[*models.Account, string](relAccounts, func(a *models.Account) string {
				return a.AccountName
			})
			c.Data["relAccounts"] = fn.Join(relAccountNames)
		}
		if len(account.ProjIds) > 0 {
			projs := models.FindProjectForMap()
			if len(projs) > 0 {
				var pnames []string
				for _, v := range account.ProjIds {
					pn, ok := projs[v]
					if ok {
						pnames = append(pnames, pn)
					}
				}
				if len(pnames) > 0 {
					c.Data["projs"] = fn.Join(pnames)
				}
			}
		}
		if !hasCopyPwdPerm && account.FieldPwd != "" {
			account.FieldPwd = ""
		}
		c.TplName = "account-detail.html"
	} else {
		//项目列表,项目可修改
		projs := models.FindProjectForMap()
		c.Data["projs"] = projs
		if len(account.ProjIds) == 0 {
			account.ProjIds = []int{0}
		}
		accountTypes := models.FindAccountTypesForMap()
		c.Data["accountTypes"] = accountTypes
		//关联账号
		if len(account.RelAccountIds) > 2 {
			relAccounts := models.FindAccountByIds(account.RelAccountIds.ToIntSlice(), true)
			var relAccountJs []models.JsData
			for _, v := range relAccounts {
				relAccountJs = append(relAccountJs, models.JsData{Id: v.AccountId, Text: v.AccountName})
			}
			c.Data["relAccounts"] = relAccountJs
		}
		c.TplName = "account-edit-form.html"
	}
}

//保存新增账号
func (c *AccountController) SaveAccount() {
	c.justPost()
	var acct models.Account
	err := c.BindJSON(&acct)
	if err != nil {
		logs.Error(err)
		c.JsonParamError("提交的数据错误")
		return
	}
	//logs.Debug(jsoniter.MarshalToString(acct))
	errmsg := acct.Valid()
	if errmsg != "" {
		c.JsonParamError(errmsg)
	}
	err = models.SaveAccount(&acct)
	if err != nil {
		c.JsonFailed(err.Error())
	}
	c.JsonOk(acct.AccountId)
}

func (c *AccountController) DelAccount() {
	accountId, err := c.GetInt("id")
	if err != nil || accountId <= 0 {
		c.JsonParamError("参数错误")
		return
	}
	var action int
	action, err = c.GetInt("action", 1)
	if err != nil || action < 0 {
		c.JsonParamError("参数错误")
		return
	}
	isOk := models.DeleteAccount(accountId, action)
	c.JsonOk(isOk)
}

func (c *AccountController) RecoverAccount() {
	accountId, err := c.GetInt("id")
	if err != nil || accountId <= 0 {
		c.JsonParamError("参数错误")
		return
	}
	acct := models.GetAccountById(accountId)
	if !acct.IsDeleted {
		c.JsonParamError("当前账号未废除")
		return
	}
	isOk := models.RecoverAccount(accountId)
	c.JsonOk(isOk)
}
