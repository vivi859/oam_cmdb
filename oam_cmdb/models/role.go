package models

import (
	fn "OAM/util"
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/validation"
)

// 超管角色
const ROLE_ROOT = "ROLE_SUPERVISOR"
const table_role_name = "sys_role"
const ROLE_LIST_CACHE_KEY = "role_list_cache"

// 角色
type Role struct {
	RoleCode   string    `orm:"PK" valid:"MaxSize(30)" form:"roleCode"`
	RoleName   string    `valid:"Required;MaxSize(30)"  form:"roleName"`
	RoleStatus int8      `form:"roleStatus"`
	UpdateTime time.Time `orm:"auto_now;type(datetime)" json:",omitempty"`
}

type RoleFun struct {
	RoleCode string
	FunCodes []int
}

func (t *Role) TableName() string {
	return table_role_name
}

func (acct *Role) Valid() string {
	valid := validation.Validation{}
	_, err := valid.Valid(acct)
	if err != nil {
		return err.Error()
	}

	return ToErrMsg(valid)
}

// 保存新增角色
func SaveRole(role *Role) (int64, error) {
	db := orm.NewOrm()
	isExist := db.QueryTable(table_role_name).Filter("role_name", role.RoleName).Exist()
	if isExist {
		return 0, errors.New("相同角色名已存在")
	}
	id, err := FlakeIDGen.NextID()
	if err != nil {
		return 0, err
	}
	role.RoleCode = strconv.FormatUint(id, 16)
	fn.GetPublicCache().Delete(ROLE_LIST_CACHE_KEY)
	return db.Insert(role)
}

func UpdateRole(role *Role) error {
	db := orm.NewOrm()
	if ROLE_ROOT == role.RoleCode {
		return errors.New("系统角色不可修改")
	}
	oldRole := &Role{RoleCode: role.RoleCode}
	err := db.Read(oldRole)
	if checkQueryErr(err) {
		return err
	}
	if role.RoleName != oldRole.RoleName {
		isExist := db.QueryTable(table_role_name).Filter("role_name", role.RoleName).Exist()
		if isExist {
			return errors.New("相同角色名已存在")
		}
	}
	fn.GetPublicCache().Delete(ROLE_LIST_CACHE_KEY)
	_, err = db.Update(&role)

	return err
}

func GetRoleById(roleCode string) *Role {
	query := orm.NewOrm()
	role := &Role{RoleCode: roleCode}
	err := query.Read(role)
	if checkQueryErr(err) {
		return nil
	}

	return role
}

func FindRoles() []*Role {
	rs := fn.QueryCacheFirst2[*Role](fn.CACHE_COMMONS, ROLE_LIST_CACHE_KEY, func() []*Role {
		var roles []*Role
		query := orm.NewOrm()
		query.QueryTable(table_role_name).Limit(-1).All(&roles)
		return roles
	})

	return rs
}

func DeleteRole(roleCode string) error {
	query := orm.NewOrm()
	role := &Role{RoleCode: roleCode}
	err := query.Read(role)
	if checkQueryErr(err) {
		return err
	}
	hasUser := query.QueryTable(table_user_name).Filter("role_code", roleCode).Exist()
	if hasUser {
		return errors.New("角色已分配用户，如需删除请先解除用户角色")
	}
	fn.GetPublicCache().Delete(ROLE_LIST_CACHE_KEY)
	_, err = query.Delete(role)
	return err
}

// 保存角色分配的权限
func SaveRoleFun(roleFun RoleFun) error {
	fn.GetPublicCache().Delete(ROLE_LIST_CACHE_KEY)
	return orm.NewOrm().DoTx(func(ctx context.Context, txOrm orm.TxOrmer) error {
		_, err := txOrm.Raw("delete from rel_role_function where role_code=?", roleFun.RoleCode).Exec()
		if err != nil {
			return err
		}
		if len(roleFun.FunCodes) > 0 {
			p, err1 := txOrm.Raw("insert into rel_role_function (role_code,fun_id) values(?,?)").Prepare()
			if err1 != nil {
				return err1
			}
			defer p.Close()
			for _, fc := range roleFun.FunCodes {
				_, err1 = p.Exec(roleFun.RoleCode, fc)
				if err1 != nil {
					return err1
				}
			}
		}
		return nil
	})
}
