package controllers

import (
	"OAM/models"
	fn "OAM/util"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	"github.com/emirpasic/gods/sets/hashset"
)

type FunController struct {
	AuthController
}

func (c *FunController) Detail() {
	id, err := c.GetInt("id")
	if err != nil {
		c.JsonParamError("参数错误")
	}
	fun := models.GetFunById(id)
	if fun != nil {
		c.JsonOk(fun)
	} else {
		c.JsonParamError("功能不存在")
	}
}

func (c *FunController) Save() {
	c.justPost()
	fun := models.Fun{}
	var err error
	if err = c.ParseForm(&fun); err != nil {
		c.JsonParamError(err.Error())
		return
	}

	if errmsg := fun.Valid(); errmsg != "" {
		c.JsonParamError(errmsg)
	}
	curUser := c.getLoginUser()
	fun.UpdateBy = curUser.UserName
	//var id int64
	if fun.FunId <= 0 {
		fun.CreateBy = fun.UpdateBy
		_, err = models.SaveFun(&fun)
	} else {
		err = models.UpdateFun(&fun)
	}
	if err == nil {
		c.JsonOk(fun.ToTreeNode())
	} else {
		c.JsonFailed(err.Error())
	}
}

func (c *FunController) DeleteFun() {
	id, err := c.GetInt("id")
	if err != nil {
		c.JsonParamError("参数错误")
	}
	err = models.DeleteFunById(id)
	logs.Info("删除功能菜单：%d, err:%s", id, err)
	if err == nil {
		c.JsonOk(id)
	} else {
		c.JsonFailed(err.Error())
	}
}

// 权限树
func (c *FunController) FunTree() {
	funs := models.FindAllFuns()
	var roleFunIds []int
	tree := models.BuildFunTree(funs, 0, roleFunIds)
	c.JsonOk(tree)
}

// 角色权限树
func (c *FunController) RoleFunTree() {
	id := c.GetString("id")
	if id == "" {
		c.JsonParamError("请选择角色")
	}
	funs := models.FindAllFuns()
	roleFunIds := models.FindFunIdsByRoleCode(id)
	tree := models.BuildFunTree(funs, 0, roleFunIds)
	c.JsonOk(tree)
}

func (c *FunController) GetLoginUserPerms() {
	permStr := c.GetString("perms")
	if len(permStr) < 2 {
		c.JsonParamError("请求的权限标识错误")
		return
	}
	funCodes := strings.Split(permStr, ",")
	funCodeSet := hashset.New()

	for _, f := range funCodes {
		if strings.HasSuffix(f, ":*") {
			childFunCodes := models.FindFunCodesByLike(f[:len(f)-2])
			if len(childFunCodes) > 0 {
				for _, cf := range childFunCodes {
					funCodeSet.Add(cf)
				}
			}
		} else {
			funCodeSet.Add(f)
		}
	}

	curUser := c.getLoginUser()
	permMap := curUser.HasPermByFunCodes(fn.ToStrSlice(funCodeSet.Values())...)
	c.JsonOk(permMap)
}

func (c *FunController) SaveRoleFun() {
	c.justPost()
	var rf models.RoleFun
	err := c.BindJSON(&rf)
	if err != nil {
		logs.Error(err)
		c.JsonParamError("提交的数据错误")
		return
	}
	err = models.SaveRoleFun(rf)
	if err != nil {
		logs.Error("保存角色[%s]分配的权限失败,ex:%s", rf.RoleCode, err)
		c.JsonFailed("保存失败")
		return
	}
	c.JsonOk(nil)
}

func (c *FunController) NavMenu() {
	curUser := c.getLoginUser()
	menu := c.GetSession("nav_menu")
	if menu == nil {
		funs := models.FindFunByRoleCode(curUser.RoleCode, 1)
		if len(funs) > 0 {
			menu = models.BuildMenuTree(funs, 0)
			c.SetSession("nav_menu", menu)
		}
	}

	c.JsonOk(menu)
}
