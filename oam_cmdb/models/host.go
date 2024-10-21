package models

import (
	fn "OAM/util"
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/validation"
	"github.com/emirpasic/gods/sets/hashset"
)

/* 主机 */
type Host struct {
	HostId     int `orm:"PK;auto"`
	HostType   uint8
	HostName   string `valid:"Required; MaxSize(50)"`
	OsName     string
	PublicIp   string `valid:"MaxSize(32)"`
	InternalIp string `valid:"MaxSize(32)"`
	SshPort    int    `json:",omitempty"`
	IsDeleted  bool
	Desc       string `valid:"MaxSize(500)"`
	//ServiceSoftwares string
	CreateTime time.Time `orm:"auto_now_add;type(datetime)" json:",omitempty"`
}

type HostVO struct {
	Host
	Apps     []*AppInfo
	Accounts []*Account
	ProjIds  fn.IntList `orm:"-"`
	//AppNames  string
	ProjNames string
}

type HostExt struct {
	Host
	AppId int
}

const table_host_name = "host_info"

func (t *Host) TableName() string {
	return table_host_name
}

func (host *HostVO) Valid() string {
	valid := validation.Validation{}
	_, err := valid.Valid(host)
	if err != nil {
		return err.Error()
	}
	if len(host.Accounts) > 0 {
		for _, a := range host.Accounts {
			if a.FieldUser == "" {
				return "账号名不能为空"
			}
		}
	}

	return ToErrMsg(valid)
}

func GetHostById(id int) *Host {
	host := Host{HostId: id}
	query := orm.NewOrm()
	err := query.Read(&host)
	if checkQueryErr(err) {
		return nil
	}
	return &host
}

// 根据id查询host及关联的应用和账号信息
func GetHostVOById(id int) *HostVO {
	host := Host{HostId: id}
	query := orm.NewOrm()
	err := query.Read(&host)
	if checkQueryErr(err) {
		return nil
	}
	hostVo := HostVO{Host: host}
	hostVo.Apps = FindAppInfoByHostId(id)
	hostVo.Accounts = FindAccountsByHostId(id)
	var ids orm.ParamsList
	num, err := query.Raw("select proj_id from rel_proj_host where host_id=?", id).ValuesFlat(&ids)
	if err == nil && num > 0 {
		hostVo.ProjIds = ToIntSlice(ids)
	}
	return &hostVo
}

// 生成主机账号表单html,前端模板语法太弱后台实现
func (host HostVO) ToBuildAccountHTML() string {
	tpl := `<tr><td><input type="hidden" id="hostAccountId%d" value="%s">
	 <input type="text" id="hostAccountName%d" maxlength="50" class="textbox" style="width:99%%" value="%s"></td>
	   <td><input data-options="validType:'maxLength[50]'" class="easyui-passwordbox" id="hostAccountPwd%d" value="%s" style="width:99%%"></td>
	   </tr>`
	acctLen := len(host.Accounts)
	var html, tmpName, tmpPwd, tmpId string
	j := 0
	for i := 0; i < 3; i++ {
		if acctLen > 0 && j < acctLen {
			tmpName = host.Accounts[j].FieldUser
			tmpPwd = host.Accounts[j].FieldPwd
			tmpId = strconv.Itoa(host.Accounts[j].AccountId)
			j = j + 1
		} else {
			tmpName = ""
			tmpPwd = ""
			tmpId = ""
		}
		html = html + fmt.Sprintf(tpl, i, tmpId, i, tmpName, i, tmpPwd)

	}
	return html
}

func FindRelHostIdByAppIds(appIds []int) map[int]*hashset.Set {
	result := make(map[int]*hashset.Set)
	query := orm.NewOrm()
	var maps []orm.Params
	sql := `SELECT r.host_id,GROUP_CONCAT(app_id) appids from rel_host_app r join host_info h on r.host_id=h.host_id 
	where r.app_id in (%s) and h.is_deleted=0  GROUP BY r.host_id`
	num, err := query.Raw(fmt.Sprintf(sql, fn.JoinInteger(",", appIds...))).Values(&maps)
	if !checkQueryErr(err) && num > 0 {
		for _, rs := range maps {
			hid, _ := strconv.Atoi(rs["host_id"].(string))
			appIds := fn.Split(rs["appids"].(string))
			tmpSet := hashset.New()
			for _, id := range appIds {
				aid, _ := strconv.Atoi(id)
				tmpSet.Add(aid)
			}
			result[hid] = tmpSet
		}
	}
	return result
}

func FindHostForPage(row int, curpage int, condi map[string]interface{}) Page[HostVO] {
	var sql_query = `SELECT h.*,GROUP_CONCAT(DISTINCT proj.proj_name) proj_names
				FROM host_info h LEFT JOIN rel_proj_host rp on rp.host_id=h.host_id
				LEFT JOIN project proj on proj.proj_id=rp.proj_id`
	var sql_count = "SELECT count(*) from host_info h "
	var params []interface{}
	var where []string
	if condi != nil {
		hid, exist := condi["hostId"]
		if exist {
			where = append(where, "h.host_id = ?")
			params = append(params, hid.(int))
		}
		kw, exist := condi["keyword"]
		if exist {
			word := kw.(string)
			isip, _ := regexp.MatchString("^\\d{1,3}\\.[\\d|\\.]*", word)
			if isip {
				where = append(where, "(h.public_ip like ? or h.internal_ip like ?)")
				word = word + "%"
				params = append(params, word, word)
			} else {
				where = append(where, "(h.host_name like ?)")
				word = "%" + word + "%"
				params = append(params, word)
			}
		}
		htype, exist := condi["htype"]
		if exist {
			where = append(where, "h.host_type = ?")
			params = append(params, htype.(int))
		}
		os, exist := condi["os"]
		if exist {
			where = append(where, "h.os_name = ?")
			params = append(params, os.(string))
		}
		justNormal, exist := condi["justNormal"]
		if exist {
			where = append(where, "h.is_deleted = ?")
			params = append(params, justNormal)
		}
		projId, exist := condi["projId"]
		if exist {
			where = append(where, "rp.proj_id =?")
			params = append(params, projId.(int))
			sql_count = sql_count + " join rel_proj_host rp on rp.host_id=h.host_id"
		}
	}
	if len(where) > 0 {
		wherestr := strings.Join(where, " and ")
		sql_count = sql_count + " where " + wherestr
		sql_query = sql_query + " where " + wherestr
	}
	sql_query = sql_query + " GROUP BY h.host_id ORDER BY h.host_id desc"
	pageData := Page[HostVO]{RowPerPage: row, PageNum: curpage}
	pageData.RawQueryPage(sql_count, sql_query, params)
	return pageData

}

func FindBaseHostForPage(row int, curpage int, condi map[string]interface{}) Page[Host] {
	var sql_query = "SELECT h.* FROM host_info h "
	var sql_count = "SELECT count(*) from host_info h "
	var params []interface{}
	var where []string
	if condi != nil {
		hid, exist := condi["hostId"]
		if exist {
			where = append(where, "h.host_id = ?")
			params = append(params, hid.(int))
		}
		kw, exist := condi["keyword"]
		if exist {
			word := kw.(string)
			isip, _ := regexp.MatchString("^\\d{1,3}\\.[\\d|\\.]*", word)
			if isip {
				where = append(where, "(h.public_ip like ? or h.internal_ip like ?)")
				word = word + "%"
				params = append(params, word, word)
			} else {
				where = append(where, "(h.host_name like ?)")
				word = "%" + word + "%"
				params = append(params, word)
			}
		}
		htype, exist := condi["htype"]
		if exist {
			where = append(where, "h.host_type = ?")
			params = append(params, htype.(int))
		}
		os, exist := condi["os"]
		if exist {
			where = append(where, "h.os_name = ?")
			params = append(params, os.(string))
		}
		justNormal, exist := condi["justNormal"]
		if exist {
			where = append(where, "h.is_deleted = ?")
			params = append(params, justNormal)
		}
		projId, exist := condi["projId"]
		if exist {
			where = append(where, "rp.proj_id =?")
			params = append(params, projId.(int))
			sql_query = sql_query + " join rel_proj_host rp on rp.host_id=h.host_id"
			sql_count = sql_count + " join rel_proj_host rp on rp.host_id=h.host_id"
		}
	}
	if len(where) > 0 {
		wherestr := strings.Join(where, " and ")
		sql_count = sql_count + " where " + wherestr
		sql_query = sql_query + " where " + wherestr
	}
	sql_query = sql_query + " ORDER BY h.is_deleted asc,h.host_id desc"
	pageData := Page[Host]{RowPerPage: row, PageNum: curpage}
	pageData.RawQueryPage(sql_count, sql_query, params)
	return pageData

}

func SaveHost(hostvo *HostVO) error {
	err := orm.NewOrm().DoTx(func(ctx context.Context, txOrm orm.TxOrmer) error {
		_, err := txOrm.Insert(&hostvo.Host)
		if err != nil {
			return err
		}
		if len(hostvo.Apps) > 0 {
			p, err1 := txOrm.Raw("insert into rel_host_app (host_id,app_id) values(?,?)").Prepare()
			if err1 != nil {
				return err1
			}
			for _, app := range hostvo.Apps {
				_, err1 = p.Exec(hostvo.HostId, app.AppId)
				if err1 != nil {
					return err1
				}
			}
			p.Close()
		}
		if len(hostvo.Accounts) > 0 {
			for _, acct := range hostvo.Accounts {
				acct.TypeId = 0
				acct.HostId = hostvo.HostId
				if hostvo.PublicIp != "" {
					acct.FieldUrl = hostvo.PublicIp
				} else {
					if hostvo.InternalIp != "" {
						acct.FieldUrl = hostvo.InternalIp
					}
				}
				if acct.FieldPwd != "" {
					acct.FieldPwd, err = fn.AesEncryptStr(acct.FieldPwd)
					if err != nil {
						return err
					}
				}
				_, err = txOrm.Insert(acct)
				if err != nil {
					return err
				}
				if len(hostvo.ProjIds) > 0 {
					p, err1 := txOrm.Raw("INSERT INTO rel_proj_account (proj_id, account_id) VALUES(?,?)").Prepare()
					if err1 != nil {
						return err1
					}
					for _, projId := range hostvo.ProjIds {
						_, err1 = p.Exec(projId, acct.AccountId)
						if err1 != nil {
							return err1
						}
					}
					p.Close()
				}
			}
		}
		if len(hostvo.ProjIds) > 0 {
			p, err1 := txOrm.Raw("INSERT INTO rel_proj_host (proj_id, host_id) VALUES(?,?)").Prepare()
			if err1 != nil {
				return err1
			}
			for _, projId := range hostvo.ProjIds {
				_, err1 = p.Exec(projId, hostvo.HostId)
				if err1 != nil {
					return err1
				}
			}
			p.Close()
		}
		return nil
	})
	return err
}

func UpdateHost(hostvo *HostVO) error {
	err := orm.NewOrm().DoTx(func(ctx context.Context, txOrm orm.TxOrmer) error {
		_, err := txOrm.Update(&hostvo.Host)
		if err != nil {
			return err
		}
		if len(hostvo.Apps) > 0 {
			txOrm.Raw("delete from rel_host_app where host_id=?", hostvo.HostId).Exec()
			p, err1 := txOrm.Raw("insert into rel_host_app (host_id,app_id) values(?,?)").Prepare()
			if err1 != nil {
				return err1
			}
			defer p.Close()
			for _, app := range hostvo.Apps {
				_, err1 = p.Exec(hostvo.HostId, app.AppId)
				if err1 != nil {
					return err1
				}
			}

		}
		if len(hostvo.Accounts) > 0 {
			oldAccounts := FindAccountsByHostId(hostvo.HostId)
			var needUpdate = false
			for _, acct := range hostvo.Accounts {
				if acct.FieldPwd != "" {
					acct.FieldPwd, err = fn.AesEncryptStr(acct.FieldPwd)
					if err != nil {
						return err
					}
				}
				if acct.AccountId == 0 {
					acct.TypeId = 0
					acct.HostId = hostvo.HostId
					if hostvo.PublicIp != "" {
						acct.FieldUrl = hostvo.PublicIp
					} else {
						if hostvo.InternalIp != "" {
							acct.FieldUrl = hostvo.InternalIp
						}
					}
					txOrm.Insert(acct)
					if len(hostvo.ProjIds) > 0 {
						p, err1 := txOrm.Raw("INSERT INTO rel_proj_account (proj_id, account_id) VALUES(?,?)").Prepare()
						if err1 != nil {
							return err1
						}
						for _, projId := range hostvo.ProjIds {
							_, err1 = p.Exec(projId, acct.AccountId)
							if err1 != nil {
								return err1
							}
						}
						p.Close()
					}
					continue
				}
				needUpdate = false
				for _, oldacct := range oldAccounts {
					if acct.AccountId == oldacct.AccountId {
						if acct.FieldUser != oldacct.FieldUser || acct.FieldPwd != oldacct.FieldPwd {
							needUpdate = true
							break
						}
					}
				}
				if needUpdate {
					_, err1 := txOrm.Update(acct, "field_user", "field_pwd")
					if err1 != nil {
						return err1
					}
				}
			}
		}

		if len(hostvo.ProjIds) > 0 {
			txOrm.Raw("delete from rel_proj_host where host_id=?", hostvo.HostId).Exec()
			p, err1 := txOrm.Raw("INSERT INTO rel_proj_host (proj_id, host_id) VALUES(?,?)").Prepare()
			if err1 != nil {
				return err1
			}
			for _, projId := range hostvo.ProjIds {
				_, err1 = p.Exec(projId, hostvo.HostId)
				if err1 != nil {
					return err1
				}
			}
			p.Close()
		}
		return nil
	})
	return err
}

func DeleteHost(id int, action int) bool {
	host := Host{HostId: id}
	err := orm.NewOrm().DoTx(func(ctx context.Context, txOrm orm.TxOrmer) error {
		var err1 error
		if action == 0 {
			_, err1 = txOrm.Delete(&host)
		} else {
			host.IsDeleted = true
			_, err1 = txOrm.Update(&host, "is_deleted")
		}

		if err1 == nil {
			txOrm.Raw("delete from rel_host_app where host_id=?", id).Exec()
			txOrm.Raw("delete from rel_proj_host where host_id=?", id).Exec()
			txOrm.Raw("update account set host_id=null where host_id=?", id).Exec()
		}
		return err1
	})

	return err == nil
}

func FindHostsByProjId(projId int) []*Host {
	var hosts []*Host
	query := orm.NewOrm()
	sql := "SELECT h.* from host_info h join rel_proj_host r on r.host_id=h.host_id where r.proj_id=? and h.is_deleted=0"
	_, err := query.Raw(sql, projId).QueryRows(&hosts)
	if checkQueryErr(err) {
		return nil
	}
	return hosts
}

//恢复被废除的主机
func RecoverHost(id int) bool {
	host := Host{HostId: id, IsDeleted: false}
	_, err := orm.NewOrm().Update(&host, "is_deleted")
	return err == nil
}
