package models

import (
	"context"
	"strings"
	"time"

	fn "OAM/util"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/validation"
)

/* 应用信息 */
type AppInfo struct {
	AppId          int    `orm:"PK;auto"`
	AppName        string `valid:"Required;MaxSize(50)"`
	AppUrl         string `valid:"MaxSize(255)"`
	AppDir         string `valid:"MaxSize(255)"`
	AppPort        int
	AppType        string
	ProjId         int
	SourcecodeRepo string
	Desc           string    `valid:"MaxSize(600)"`
	CreateTime     time.Time `orm:"auto_now_add;type(datetime)" json:",omitempty"`
	DevLang        string
}

type AppInfoVO struct {
	AppInfo
	HostIds fn.IntList
}

type AppInfoRow struct {
	AppInfo
	HostNames string
}

const table_app_name = "app_info"

func (t *AppInfo) TableName() string {
	return table_app_name
}

func (acct *AppInfo) Valid() string {
	valid := validation.Validation{}
	_, err := valid.Valid(acct)
	if err != nil {
		return err.Error()
	}
	return ToErrMsg(valid)
}

func GetAppInfoById(id int) *AppInfo {
	app := AppInfo{AppId: id}
	query := orm.NewOrm()
	err := query.Read(&app)
	if checkQueryErr(err) {
		return nil
	}
	return &app
}

func GetAppInfoVOById(id int) *AppInfoVO {
	app := AppInfo{AppId: id}
	query := orm.NewOrm()
	err := query.Read(&app)
	if checkQueryErr(err) {
		return nil
	}
	if len(app.AppType) > 0 {
		app.AppType = strings.TrimSpace(app.AppType)
	}
	appVO := AppInfoVO{AppInfo: app}

	var ids orm.ParamsList
	num, err := query.Raw("select host_id from rel_host_app where app_id=?", id).ValuesFlat(&ids)
	if err == nil && num > 0 {
		appVO.HostIds = ToIntSlice(ids)
	}
	return &appVO
}

func FindAppInfoByHostId(hostId int) []*AppInfo {
	var apps []*AppInfo
	_, err := orm.NewOrm().Raw("select a.* from app_info a join rel_host_app rel on a.app_id=rel.app_id where rel.host_id=?", hostId).QueryRows(&apps)
	if checkQueryErr(err) {
		return nil
	}
	return apps
}

// 通过主机id查找部署的应用名称列表
func FindAppNameByHostId(hostId int) []string {
	var names []string
	var ids orm.ParamsList
	num, err := orm.NewOrm().Raw("select a.app_name from app_info a join rel_host_app rel on a.app_id=rel.app_id where rel.host_id=?", hostId).ValuesFlat(&ids)
	if err == nil && num > 0 {
		names = fn.ToStrSlice(ids)
	}
	return names
}

func FindAppInfoByProjId(projId int) []*AppInfo {
	var apps []*AppInfo
	query := orm.NewOrm()
	selector := query.QueryTable(table_app_name)
	_, err := selector.Filter("proj_id", projId).All(&apps)
	if checkQueryErr(err) {
		return nil
	}
	return apps
}

func FindAppInfoForPage(row int, curpage int, condi map[string]interface{}) Page[AppInfoRow] {
	var sql_query = `SELECT app.*,GROUP_CONCAT(DISTINCT h.host_name) as host_names
	FROM app_info app LEFT JOIN rel_host_app ra on app.app_id=ra.app_id 
	LEFT JOIN host_info h on ra.host_id=h.host_id`
	var sql_count = "SELECT count(*) from app_info app "

	var params []interface{}
	var where []string
	if condi != nil {
		hid, exist := condi["hostId"]
		if exist {
			where = append(where, "ra.host_id = ?")
			params = append(params, hid.(int))
			sql_count = sql_count + " LEFT JOIN rel_host_app ra on app.app_id=ra.app_id"
		}
		kw, exist := condi["keyword"]
		if exist {
			word := "%" + kw.(string) + "%"
			where = append(where, "(app.app_name like ?)")
			params = append(params, word)
		}
		appType, exist := condi["appType"]
		if exist {
			where = append(where, "app.app_type = ?")
			params = append(params, appType.(string))
		}
		// justNormal, exist := condi["justNormal"]
		// if exist {
		// 	where = append(where, "h.is_deleted = ?")
		// 	params = append(params, justNormal)
		// }
		projId, exist := condi["projId"]
		if exist {
			where = append(where, "app.proj_id =?")
			params = append(params, projId.(int))
		}
	}
	if len(where) > 0 {
		wherestr := strings.Join(where, " and ")
		sql_count = sql_count + " where " + wherestr
		sql_query = sql_query + " where " + wherestr
	}
	sql_query = sql_query + " GROUP BY app.app_id"
	pageData := Page[AppInfoRow]{RowPerPage: row, PageNum: curpage}
	pageData.RawQueryPage(sql_count, sql_query, params)
	return pageData
}

func SaveAppInfo(appVO *AppInfoVO) error {
	return orm.NewOrm().DoTx(func(ctx context.Context, txOrm orm.TxOrmer) error {
		_, err := txOrm.Insert(&appVO.AppInfo)
		if err != nil {
			return err
		}
		if len(appVO.HostIds) > 0 {
			p, err1 := txOrm.Raw("insert into rel_host_app (host_id,app_id) values(?,?)").Prepare()
			if err1 != nil {
				return err1
			}
			defer p.Close()
			for _, hostId := range appVO.HostIds {
				_, err1 = p.Exec(hostId, appVO.AppId)
				if err1 != nil {
					return err1
				}
			}
		}
		return nil
	})
}

func UpdateAppInfo(appVO *AppInfoVO) error {
	return orm.NewOrm().DoTx(func(ctx context.Context, txOrm orm.TxOrmer) error {
		_, err := txOrm.Update(&appVO.AppInfo)
		if err != nil {
			return err
		}
		if len(appVO.HostIds) > 0 {
			txOrm.Raw("delete from rel_host_app where app_id=?", appVO.AppId).Exec()
			p, err1 := txOrm.Raw("insert into rel_host_app (host_id,app_id) values(?,?)").Prepare()
			if err1 != nil {
				return err1
			}
			defer p.Close()
			for _, hostId := range appVO.HostIds {
				_, err1 = p.Exec(hostId, appVO.AppId)
				if err1 != nil {
					return err1
				}
			}
		}
		return nil
	})
}

func DeleteAppInfo(id int) bool {
	app := AppInfo{AppId: id}
	err := orm.NewOrm().DoTx(func(ctx context.Context, txOrm orm.TxOrmer) error {
		i, err := txOrm.Delete(&app)
		if err == nil && i > 0 {
			txOrm.Raw("delete from rel_host_app where app_id=?", id).Exec()
		}
		return err
	})
	return err == nil
}
