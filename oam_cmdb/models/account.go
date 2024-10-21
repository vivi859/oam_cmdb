package models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sort"

	"strconv"
	"strings"
	"time"

	fn "OAM/util"

	"github.com/astaxie/beego/logs"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/validation"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/exp/slices"
)

// 账号资料
type Account struct {
	AccountId     int    `orm:"PK;auto"`
	AccountName   string `valid:"Required; MaxSize(100)" label:"账号名称"`
	TypeId        int
	CreateBy      string
	UpdateBy      string
	UpdateTime    time.Time       `orm:"auto_now;type(datetime)" json:",omitempty"`
	CreateTime    time.Time       `orm:"auto_now_add;type(datetime)" json:",omitempty"`
	FieldUser     string          `valid:"MaxSize(100)"`
	FieldPwd      string          `valid:"MaxSize(150)" json:",omitempty"`
	FieldRePwd    string          `orm:"-" json:",omitempty"`
	FieldUrl      string          `valid:"MaxSize(255)"`
	FieldRemark   string          `valid:"MaxSize(800)"`
	FieldOther    string          `json:",omitempty"`
	RelAccountIds SqlLikeArrayStr `json:",omitempty"`
	IsDeleted     bool
	HostId        int
	ProjName      string     `orm:"-"`
	TypeName      string     `orm:"-"`
	ProjIds       fn.IntList `orm:"-"`
}

//定义一个方便like查询的多值组合字符串，格式示例：,1,2,4, => like '%,4,%'
type SqlLikeArrayStr string

//RelAccountIds 转json数组
func (str SqlLikeArrayStr) MarshalJSON() ([]byte, error) {
	return json.Marshal(str.ToIntSlice())
}

func (a *SqlLikeArrayStr) UnmarshalJSON(b []byte) error {
	var intArr []int
	jsoniter.Unmarshal(b, &intArr)
	sort.Ints(intArr)
	*a = toSqlLikeArrayStr(intArr)
	return nil
}

func (str SqlLikeArrayStr) ToIntSlice() []int {
	if len(str) < 3 {
		return make([]int, 0)
	}
	tmpStr := string(str[1 : len(str)-1])
	strArr := strings.Split(tmpStr, ",")
	return fn.ToIntSlice(strArr)
}

func toSqlLikeArrayStr(intSlice []int) SqlLikeArrayStr {
	if len(intSlice) > 0 {
		return SqlLikeArrayStr(fmt.Sprintf(",%s,", fn.JoinInteger(",", intSlice...)))
	}
	return ""
}

func (str SqlLikeArrayStr) ToJSONString() string {
	if len(str) < 3 {
		return "[]"
	}
	tmpStr, err := jsoniter.MarshalToString(str)
	if err != nil {
		panic("类型ForSqlLikeArrayStr，值：" + str + "json序列化失败")
	}
	return tmpStr
}

const table_account_name = "account"

func (t *Account) TableName() string {
	return table_account_name
}

func (acct *Account) Valid() string {
	valid := validation.Validation{}
	_, err := valid.Valid(acct)
	if err != nil {
		return err.Error()
	}
	if acct.FieldPwd != acct.FieldRePwd {
		valid.SetError("FieldPwd", "两次密码不相同")
	}
	return ToErrMsg(valid)
}

func (acct *Account) SetRelAccountIds(ids []int) {
	acct.RelAccountIds = toSqlLikeArrayStr(ids)
}

func (acct *Account) GetRelAccountIds() []int {
	ids := acct.RelAccountIds.ToIntSlice()
	if len(ids) > 0 {
		sort.Ints(ids)
	}
	return ids
}

const sql_findaccountbyprojid = `SELECT a.* from account a join rel_proj_account r on a.account_id=r.account_id 
where r.proj_id=? and a.is_deleted=0 `

//查询项目下所有账号
func FindAccountsByProjId(projId int) []*Account {
	var fields []*Account
	_, err := orm.NewOrm().Raw(sql_findaccountbyprojid, projId).QueryRows(&fields)
	if checkQueryErr(err) {
		return nil
	}
	return fields
}

func FindAccountsByHostId(hostId int) []*Account {
	var fields []*Account
	_, err := orm.NewOrm().QueryTable(table_account_name).Filter("host_id", hostId).Filter("is_deleted", 0).All(&fields)
	if checkQueryErr(err) {
		return nil
	}
	return fields
}

func GetAccountById(id int) *Account {
	acct := Account{AccountId: id}
	query := orm.NewOrm()
	err := query.Read(&acct)
	if checkQueryErr(err) {
		return nil
	}
	var ids orm.ParamsList
	num, err := query.Raw("select proj_id from rel_proj_account where account_id=?", id).ValuesFlat(&ids)
	if err == nil && num > 0 {
		acct.ProjIds = ToIntSlice(ids)
	}
	return &acct
}

// 根据id列表查询账号.
// 参数：simpleResult true=只查询"account_id", "account_name", "field_user"字段
func FindAccountByIds(ids []int, simpleResult bool) []*Account {
	var lst []*Account
	query := orm.NewOrm()
	querySeter := query.QueryTable("account").Filter("account_id__in", ids)
	var err error
	if simpleResult {
		_, err = querySeter.All(&lst, "account_id", "account_name", "field_user")
	} else {
		_, err = querySeter.All(&lst)
	}
	if checkQueryErr(err) {
		return nil
	}
	return lst
}

//分页查询账号
func FindAccountForPage(row int, curpage int, condi map[string]interface{}) Page[Account] {
	//querySeter := orm.NewOrm().QueryTable(table_account_name)
	var sql_query = `SELECT a.account_id,a.account_name,a.type_id,a.create_by,a.update_by,a.update_time,a.create_time,a.field_user,
	a.field_pwd,a.field_url,a.field_remark,a.is_deleted,a.host_id from account a `
	var sql_count = "SELECT count(*) from account a "
	var params []interface{}
	var where []string
	if condi != nil {
		pid, exist := condi["projId"]
		if exist {
			sql_query = sql_query + " join rel_proj_account r on a.account_id=r.account_id "
			sql_count = sql_count + " join rel_proj_account r on a.account_id=r.account_id "
			params = append(params, pid.(int))
			where = append(where, "r.proj_id = ?")
		}

		tid, exist := condi["typeId"]
		if exist {
			params = append(params, tid.(int))
			where = append(where, "a.type_id=?")
		}

		aname, exist := condi["accountName"]
		if exist {
			regx, _ := regexp.Compile(`[%&*=<>!'-]`)
			cleanName := regx.ReplaceAllString(aname.(string), "")
			tmpKw := "%" + cleanName + "%"
			params = append(params, tmpKw, tmpKw, tmpKw)
			where = append(where, "(a.account_name like ? or a.field_user like ? or a.field_other like ?)")
		}

		notDeleted, exist := condi["notDeleted"]
		if exist && notDeleted.(string) == "1" {
			params = append(params, 0)
			where = append(where, "a.is_deleted =?")
		}
	}
	if len(where) > 0 {
		wherestr := strings.Join(where, " and ")
		sql_count = sql_count + " where " + wherestr
		sql_query = sql_query + " where " + wherestr
	}
	sql_query = sql_query + " order by a.is_deleted asc,a.account_id desc"
	pageData := Page[Account]{RowPerPage: row, PageNum: curpage}
	pageData.RawQueryPage(sql_count, sql_query, params)
	//查账号关联的项目
	if len(pageData.Rows) > 0 {
		var ids []int
		for _, acct := range pageData.Rows {
			ids = append(ids, acct.AccountId)
		}
		var maps []orm.Params
		var tmpIntMap map[string]int
		num, err := orm.NewOrm().Raw("select proj_id,account_id from rel_proj_account where account_id in (" + fn.JoinInteger(",", ids...) + ")").Values(&maps)
		if err == nil && num > 0 {
			for _, acct := range pageData.Rows {
				tmpProjIds := make([]int, 0)
				for _, rel := range maps {
					tmpIntMap = ToIntValueMap(rel)
					if acct.AccountId == tmpIntMap["account_id"] {
						tmpProjIds = append(tmpProjIds, tmpIntMap["proj_id"])
					}
				}
				if len(tmpProjIds) > 0 {
					acct.ProjIds = tmpProjIds
				}
			}
		}
	}
	return pageData
}

//保存账号
func SaveAccount(a *Account) error {
	// 关系账号是双向关联关系，互相都保存关系
	var delRelIds []int
	var addRelIds []int
	var updateRelIdsAccount []Account
	newRelIds := a.GetRelAccountIds()

	if a.AccountId > 0 {
		oldAccount := GetAccountById(a.AccountId)
		if oldAccount == nil {
			return errors.New("编辑的账号不存在")
		}
		if a.TypeId != oldAccount.TypeId {
			return errors.New("账号类型不可修改")
		}
		oldRelIds := oldAccount.GetRelAccountIds()
		if !slices.Equal[int](newRelIds, oldRelIds) {
			if len(oldRelIds) > 0 {
				//找出被删除的关联账号
				if len(newRelIds) == 0 {
					delRelIds = oldRelIds
				} else {
					//要删除的
					for _, d := range oldRelIds {
						if !slices.Contains[int](newRelIds, d) {
							delRelIds = append(delRelIds, d)
						}
					}
					//新增的
					for _, d := range newRelIds {
						if !slices.Contains[int](oldRelIds, d) {
							addRelIds = append(addRelIds, d)
						}
					}
				}

			} else {
				//新增的
				addRelIds = newRelIds
			}
		}
	} else {
		addRelIds = newRelIds
	}

	if len(delRelIds) > 0 {
		var lst []Account
		var tmpIds []int

		c, err := orm.NewOrm().QueryTable("account").Filter("account_id__in", delRelIds).All(&lst, "account_id", "rel_account_ids")
		if err == nil && c > 0 {
			for _, tmpAccount := range lst {
				tmpIds = tmpAccount.GetRelAccountIds()
				if slices.Contains[int](tmpIds, a.AccountId) {
					tmpAccount.SetRelAccountIds(fn.RemoveSlice[int](tmpIds, a.AccountId))
					updateRelIdsAccount = append(updateRelIdsAccount, tmpAccount)
				}
			}
		}
	}

	var err error
	if a.FieldPwd != "" {
		a.FieldPwd, err = fn.AesEncryptStr(a.FieldPwd)
		if err != nil {
			return err
		}
	}

	if a.TypeId > 0 {
		otherMap := make(map[string]interface{})
		if a.FieldOther != "" {
			jsoniter.UnmarshalFromString(a.FieldOther, &otherMap)
		}
		fields := FindFieldsByTypeId(a.TypeId)
		isNeedReMarshal := false
		//验证动态属性
		for _, field := range fields {
			v, ok := otherMap[field.FieldKey]
			if field.IsRequired && (!ok || v == "") {
				return errors.New("请输入" + field.FieldName)
			}
			if ok {
				if field.ValueType == 0 {
					sv := v.(string)
					if field.MaxLen < len(sv) {
						return errors.New(field.FieldName + "长度应小于" + strconv.Itoa(field.MaxLen))
					}
					//加密字段
					if field.IsCiphertext && v != "" {
						otherMap[field.FieldKey], err = fn.AesEncryptStr(sv)
						if err != nil {
							return err
						}
						isNeedReMarshal = true
					}
				} else if field.ValueType == 2 {
					if field.MaxLen != 0 && field.MaxLen > v.(int) {
						return errors.New(field.FieldName + "应小于" + strconv.Itoa(field.MaxLen))
					}
				}

			}
		}
		if isNeedReMarshal {
			a.FieldOther, err = jsoniter.MarshalToString(otherMap)
			if err != nil {
				return err
			}
		}
	}

	return orm.NewOrm().DoTx(func(ctx context.Context, txOrm orm.TxOrmer) error {
		if a.AccountId > 0 {
			_, err := txOrm.Update(a)
			if err != nil {
				return err
			}
			if len(a.ProjIds) > 0 {
				txOrm.Raw("delete from rel_proj_account where account_id=?", a.AccountId).Exec()
			}
		} else {
			_, err := txOrm.Insert(a)
			if err != nil {
				return err
			}
		}
		if len(a.ProjIds) > 0 {
			prep, err := txOrm.Raw("insert into rel_proj_account(account_id,proj_id) values(?,?)").Prepare()
			if err != nil {
				return err
			}
			for _, id := range a.ProjIds {
				prep.Exec(a.AccountId, id)
			}
			prep.Close()
		}
		//账号保存后再更新关联账号
		if len(addRelIds) > 0 {
			var lst []Account
			var tmpIds []int
			c, err := orm.NewOrm().QueryTable("account").Filter("account_id__in", addRelIds).All(&lst, "account_id", "rel_account_ids")
			if err == nil && c > 0 {
				for _, tmpAccount := range lst {
					tmpIds = tmpAccount.GetRelAccountIds()
					if !slices.Contains[int](tmpIds, a.AccountId) {
						tmpAccount.SetRelAccountIds(append(tmpIds, a.AccountId))
						updateRelIdsAccount = append(updateRelIdsAccount, tmpAccount)
					}
				}
			}
		}

		//更新账号之间互关联关系
		if len(updateRelIdsAccount) > 0 {
			prep, err := txOrm.Raw("update account set rel_account_ids=? where account_id=?").Prepare()
			if err != nil {
				return err
			}
			for _, acct := range updateRelIdsAccount {
				prep.Exec(acct.RelAccountIds, acct.AccountId)
			}
			prep.Close()
		}
		return nil
	})
}

func DeleteAccount(id int, action int) bool {
	var err error
	if action == 1 {
		acct := Account{AccountId: id, IsDeleted: true}
		_, err = orm.NewOrm().Update(&acct, "is_deleted")
	} else if action == 0 {
		acct := Account{AccountId: id}
		db := orm.NewOrm()
		err = db.DoTx(func(ctx context.Context, txOrm orm.TxOrmer) error {
			txOrm.Raw("delete from rel_proj_account where account_id=?", id).Exec()
			_, err1 := txOrm.Delete(&acct)
			return err1
		})
	}
	return err == nil
}

//补全前端要显示的属性,如:密码解密,类型名称,项目名称
func CompletionVOProps(isNeedProjName bool, isNeedPasswd bool, rows []*Account) error {
	if len(rows) == 0 {
		return nil
	}
	var projs map[int]string
	accountTypes := FindAccountTypesForMap()
	if isNeedProjName {
		projs = FindProjectForMap()
	}
	for _, a := range rows {
		if !isNeedPasswd {
			a.FieldPwd = ""
		}
		if a.FieldPwd != "" {
			plainPwd, err := fn.AesDecryptStr(a.FieldPwd)
			if err != nil {
				logs.Error("密码解密异常,账号%s,%s", a.AccountName, a.FieldPwd)
				return errors.New("服务器异常")
			}
			a.FieldPwd, err = fn.RSAEncryptBase64Str(plainPwd, "")
			if err != nil {
				logs.Error("密码rsa加密异常,账号%s,%s", a.AccountName, a.FieldPwd)
				return errors.New("服务器异常")
			}
		}

		t, ok := accountTypes[strconv.Itoa(a.TypeId)]
		if ok {
			a.TypeName = t.(string)
		} else {
			a.TypeName = "无"
		}
		if isNeedProjName {
			if len(a.ProjIds) > 0 {
				var tmpProjName string
				for _, pid := range a.ProjIds {
					pn, ok := projs[pid]
					if ok {
						tmpProjName = tmpProjName + ", " + pn
					}
				}
				a.ProjName = tmpProjName[1:]
			} else {
				a.ProjName = "无"
			}
		}
	}
	return nil
}

//恢复被废除的主机
func RecoverAccount(id int) bool {
	account := Account{AccountId: id, IsDeleted: false}
	_, err := orm.NewOrm().Update(&account, "is_deleted")
	return err == nil
}
