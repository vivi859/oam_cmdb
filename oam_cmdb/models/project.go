package models

import (
	fn "OAM/util"
	"context"
	"errors"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
)

/* 项目 */
type Project struct {
	ProjId      int       `orm:"PK;auto" form:"projId"`
	ProjName    string    `form:"projName" valid:"Required; MaxSize(50)" label:"项目名称"`
	ProjDesc    string    `form:"projDesc" valid:"MaxSize(500)"  label:"项目描述"`
	CreateTime  time.Time `orm:"auto_now_add"`
	CreatorId   int
	CreatorName string `orm:"-"`
	//MemberIds   util.IntList `orm:"-"`
	UpdaterId int         `orm:"-" json:"-"`
	Members   []*UserInfo `orm:"-"`
	Docs      []*Document `orm:"-"`
	//Accounts  []*Account  `orm:"-"`
}

const table_proj_name = "project"

func (dgrp *Project) Valid(isAdd bool) string {
	valid := validation.Validation{}
	_, err := valid.Valid(dgrp)
	if err != nil {
		return err.Error()
	}
	//新增时
	if isAdd {
		if isExistSameName(dgrp.ProjName) {
			valid.SetError("GroupName", "分组名称重复")
		}
	}

	return ToErrMsg(valid)
}

func isExistSameName(projName string) bool {
	return orm.NewOrm().QueryTable(table_proj_name).Filter("proj_name", projName).Exist()
}

//根据id查项目信息
func GetProjectById(projId int) *Project {
	query := orm.NewOrm()
	g := Project{ProjId: projId}
	err := query.Read(&g)
	if checkQueryErr(err) {
		return nil
	}
	//查创建者名称
	creator := GetUserById(g.CreatorId)
	if creator != nil {
		if creator.RealName != "" {
			g.CreatorName = creator.RealName
		} else {
			g.CreatorName = creator.UserName
		}
	}
	return &g
}

//查询所有项目
func FindAllProjects() []*Project {
	var datums []*Project
	orm.NewOrm().QueryTable(table_proj_name).All(&datums)
	return datums
}

const proj_map_cache_key = "project_map"

func FindProjectForMap() map[int]string {
	/* cache := fn.GetPublicCache()
	projMap, err := cache.Get(proj_map_cache_key)
	if err == nil {
		return projMap.(map[int]string)
	} else {
		r := make(orm.Params)
		_, err := orm.NewOrm().Raw("select proj_id,proj_name from project").RowsToMap(&r, "proj_id", "proj_name")
		if err != nil {
			logs.Error(err)
			return nil
		}
		newProjMap := ToIntKeyMap(r)
		cache.Put(proj_map_cache_key, newProjMap)
		return newProjMap
	} */
	return fn.QueryCacheFirst(fn.GetPublicCache(), proj_map_cache_key, func() map[int]string {
		r := make(orm.Params)
		_, err := orm.NewOrm().Raw("select proj_id,proj_name from project").RowsToMap(&r, "proj_id", "proj_name")
		if err != nil {
			logs.Error(err)
			return nil
		}
		return ToIntKeyMap(r)
	})
}

// 查询某用户有权限的项目
func FindAuthorizedProjects(userId int) []*Project {
	var datums []*Project
	orm.NewOrm().Raw("select p.* from project p join rel_proj_user r on p.proj_id=r.proj_id where r.user_id=?", userId).QueryRows(&datums)
	return datums
}

func FindAuthorizedProjectIds(userId int) []int {
	var ids orm.ParamsList
	num, err := orm.NewOrm().Raw("select proj_id from rel_proj_user where user_id=?", userId).ValuesFlat(&ids)
	if err == nil && num > 0 {
		return ToIntSlice(ids)
	}
	return nil
}

//判断用户是否有查看项目的权限
func (proj Project) HasPermission(user LoginUser) bool {
	if user.IsSupervisor() || user.UserId == proj.CreatorId {
		return true
	}
	if len(proj.Members) > 0 {
		for _, m := range proj.Members {
			if m.UserId == user.UserId {
				return true
			}
		}
	}
	return false
}

//新增项目,成功会返回新的id
func SaveProject(project *Project) (int, error) {
	db := orm.NewOrm()
	err := db.DoTx(func(ctx context.Context, txOrm orm.TxOrmer) error {
		_, err1 := txOrm.Insert(project)
		if err1 == nil {
			_, err1 = txOrm.Raw("insert into rel_proj_user(proj_id,user_id) values (?,?)", project.ProjId, project.CreatorId).Exec()
		}
		return err1
	})
	if err != nil {
		return 0, err
	}
	fn.GetPublicCache().Delete(proj_map_cache_key)
	return project.ProjId, nil
}

//修改项目
func UpdateProject(project *Project) (int64, error) {
	db := orm.NewOrm()
	oldProj := GetProjectById(project.ProjId)
	if oldProj == nil {
		return 0, errors.New("项目不存在")
	}
	project.CreatorId = oldProj.CreatorId
	_, err := db.Update(project)
	if err != nil {
		return 0, err
	}
	fn.GetPublicCache().Delete(proj_map_cache_key)
	return 1, nil
}

//删除项目
func DeleteProject(projId int) bool {
	db := orm.NewOrm()
	proj := GetProjectById(projId)
	if proj == nil {
		return false
	}
	err := db.DoTx(func(ctx context.Context, txOrm orm.TxOrmer) error {
		//删除成员
		txOrm.Raw("delete from rel_proj_user where proj_id=?", proj.ProjId).Exec()
		//删除账号关系
		txOrm.Raw("delete from rel_proj_account where proj_id=?", proj.ProjId).Exec()
		//删除文档关系
		txOrm.Raw("update document set proj_id=0 where proj_id=?", proj.ProjId).Exec()
		//删除项目
		_, err1 := db.Delete(&Project{ProjId: proj.ProjId})
		fn.GetPublicCache().Delete(proj_map_cache_key)
		return err1
	})
	return err == nil
}

//删除项目成员
func DeleteProjectMember(projId int, userIds ...int) bool {
	var err error
	if len(userIds) == 1 {
		_, err = orm.NewOrm().Raw("delete from rel_proj_user where proj_id=? and user_id=?", projId, userIds[0]).Exec()
	} else {
		_, err = orm.NewOrm().Raw("delete from rel_proj_user where proj_id=? and user_id in (?)", projId, fn.JoinInteger(",", userIds...)).Exec()
	}
	if err != nil {
		logs.Error("删除项目成员失败", err)
		return false
	}
	return true
}

//添加项目成员
func AddProjectMember(projId int, userIds ...int) error {
	var idsLen = len(userIds)
	if idsLen > 0 {
		pstmt, err := orm.NewOrm().Raw("replace into rel_proj_user(proj_id,user_id) values (?,?)").Prepare()
		if err != nil {
			return err
		}
		defer pstmt.Close()
		for _, id := range userIds {
			pstmt.Exec(projId, id)
		}
	}
	return nil
}
