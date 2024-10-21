package models

import (
	fn "OAM/util"
	"errors"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/validation"
)

type Fun struct {
	FunId      int    `orm:"PK;auto" form:"funId"`
	FunName    string `form:"funName" valid:"Required;MaxSize(30)"`
	FunCode    string `form:"funCode" valid:"Required;MaxSize(100)"`
	FunType    int    `form:"funType" valid:"Required"`
	ParentId   int    `form:"parentId" valid:"Required"`
	FunOrder   uint16 `form:"funOrder" valid:"Required"`
	FunLevel   uint8  `form:"funLevel"`
	FunUrl     string `form:"funUrl"`
	MenuClass  string `form:"menuClass"`
	CreateBy   string
	UpdateBy   string
	UpdateTime time.Time `orm:"auto_now;type(datetime)" json:",omitempty"`
	CreateTime time.Time `orm:"auto_now_add;type(datetime)" json:",omitempty"`
}

const table_fun_name = "sys_fun"

func (t *Fun) TableName() string {
	return table_fun_name
}

func (acct *Fun) Valid() string {
	valid := validation.Validation{}
	_, err := valid.Valid(acct)
	if err != nil {
		return err.Error()
	}
	if len(acct.FunUrl) > 0 {
		urls := strings.Split(acct.FunUrl, ",")
		for _, url := range urls {
			if !ValidUrl(url) {
				return "接口地址格式错误"
			}
		}
	}

	return ToErrMsg(valid)
}

func GetFunById(id int) *Fun {
	fun := &Fun{FunId: id}
	query := orm.NewOrm()
	err := query.Read(fun)
	if checkQueryErr(err) {
		return nil
	}

	return fun
}

// 保存新增功能
func SaveFun(fun *Fun) (int64, error) {
	db := orm.NewOrm()
	isExist := db.QueryTable(table_fun_name).Filter("fun_code", fun.FunCode).Exist()
	if isExist {
		return 0, errors.New("功能标识冲突")
	}
	fun.FunCode = strings.TrimSpace(fun.FunCode)
	fun.FunUrl = strings.TrimSpace(fun.FunUrl)
	if fun.ParentId == 0 {
		fun.FunLevel = 1
	} else {
		parentFun := &Fun{FunId: fun.ParentId}
		err := db.Read(parentFun)
		if checkQueryErr(err) {
			return 0, err
		}
		fun.FunLevel = parentFun.FunLevel + 1
	}
	fn.GetFunGroupCache().ClearAll()
	return db.Insert(fun)
}

// 修改功能
func UpdateFun(fun *Fun) error {
	db := orm.NewOrm()
	oldFun := &Fun{FunId: fun.FunId}
	err := db.Read(oldFun)
	if checkQueryErr(err) {
		return err
	}
	fun.FunCode = oldFun.FunCode
	fun.ParentId = oldFun.ParentId
	fun.FunLevel = oldFun.FunLevel
	_, err = db.Update(fun)
	return err
}

// 删除功能
func DeleteFunById(funId int) error {
	query := orm.NewOrm()
	oldFun := &Fun{FunId: funId}
	err := query.Read(oldFun)
	if checkQueryErr(err) {
		return err
	}
	hasChild := query.QueryTable(table_fun_name).Filter("parent_id", funId).Exist()
	if hasChild {
		return errors.New("有子菜单或功能，不允许删除")
	}
	var c int
	err = query.Raw("select count(*) from rel_role_function where fun_id=?", funId).QueryRow(&c)
	if err == nil && c > 0 {
		return errors.New("已分配，不允许删除！ 如需要删除请先取消相关角色权限")
	}
	fn.GetFunGroupCache().ClearAll()
	_, err = query.Delete(oldFun)
	return err
}

func FindAllFuns() []*Fun {
	var funs []*Fun
	db := orm.NewOrm()
	db.QueryTable(table_fun_name).Limit(-1).OrderBy("fun_order").All(&funs)
	return funs
}

// 查询角色权限
func FindFunByRoleCode(code string, funType int) []*Fun {
	var funs []*Fun
	db := orm.NewOrm()
	sql := "SELECT f.* from sys_fun f"
	params := make([]interface{}, 2)
	if code != ROLE_ROOT {
		sql += " join rel_role_function rel on f.fun_id=rel.fun_id where rel.role_code=?"
		params = append(params, code)
	} else {
		sql += " where 1=1"
	}
	if funType > 0 {
		sql += " and f.fun_type=?"
		params = append(params, funType)
	}

	sql += " order by f.fun_order"
	db.Raw(sql).SetArgs(params).QueryRows(&funs)
	return funs
}

func FindFunIdsByRoleCode(code string) []int {
	var funs []int
	db := orm.NewOrm()
	db.Raw("SELECT f.fun_id from sys_fun f join rel_role_function rel on f.fun_id=rel.fun_id where rel.role_code=?", code).QueryRows(&funs)
	return funs
}

func FindFunCodesByLike(funCodePrefix string) []string {
	return fn.QueryCacheFirst[[]string](fn.GetFunGroupCache(), funCodePrefix, func() []string {
		var codes []string
		db := orm.NewOrm()
		db.Raw("select fun_code from sys_fun where fun_code like ?", funCodePrefix+"%").QueryRows(&codes)
		return codes
	})
}
