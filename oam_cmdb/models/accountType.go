package models

import (
	fn "OAM/util"
	"context"
	"errors"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/validation"
)

const table_type_name = "account_type"

/* 账号类型 */
type AccountType struct {
	TypeId     int `orm:"PK;auto"`
	TypeName   string
	UpdateTime time.Time       `orm:"auto_now;type(datetime)"`
	CreateTime time.Time       `orm:"auto_now_add;type(datetime)"`
	Fields     []*AccountField `orm:"-"`
}

func (t *AccountType) TableName() string {
	return table_type_name
}

func (dtype *AccountType) Valid() string {
	valid := validation.Validation{}
	valid.Required(dtype.TypeName, "TypeName").Message("请输入类型名称")
	valid.MaxSize(dtype.TypeName, 30, "TypeName").Message("名称最大长度30个字符")
	//valid.MinSize(dtype.Fields, 1, "Fields").Message("至少定义一个属性")
	return ToErrMsg(valid)
}

//根据ID查询账号类型和它的动态属性
func GetAccountTypeById(accountTypeId int, isFetchFields bool) *AccountType {
	atype := AccountType{TypeId: accountTypeId}
	query := orm.NewOrm()
	err := query.Read(&atype)
	if checkQueryErr(err) {
		return nil
	}
	if isFetchFields {
		fields := FindFieldsByTypeId(accountTypeId)
		if fields != nil {
			atype.Fields = fields
		}
	}
	return &atype
}

//查询所有账号类型
func FindAllAccountTypes() []*AccountType {
	var datums []*AccountType
	orm.NewOrm().QueryTable(table_type_name).All(&datums)
	return datums
}

const type_map_cache_key = "account_types_map"

func FindAccountTypesForMap() orm.Params {
	cache := fn.GetPublicCache()
	typeMap, err := cache.Get(type_map_cache_key)
	if err == nil {
		return typeMap.(orm.Params)
	} else {
		r := make(orm.Params)
		_, err := orm.NewOrm().Raw("select type_id,type_name from account_type").RowsToMap(&r, "type_id", "type_name")
		if err != nil {
			logs.Error(err)
		}
		if len(r) > 0 {
			cache.Put(type_map_cache_key, r)
		}
		return r
	}
}

//新增账号类型
func SaveOrUpdateAccountType(dtype *AccountType) error {
	db := orm.NewOrm()
	var err error
	if dtype.TypeId > 0 {
		//update
		oldType := GetAccountTypeById(dtype.TypeId, false)
		if oldType == nil {
			return errors.New("类型不存在")
		}

		err = db.DoTx(func(ctx context.Context, tx orm.TxOrmer) error {
			_, err1 := tx.Update(dtype)
			for _, f := range dtype.Fields {
				f.TypeId = dtype.TypeId
				f.FieldId = 0
			}
			if err1 == nil {
				fn.GetPublicCache().Delete(type_map_cache_key)
				_, err1 = tx.Raw("delete from account_field where type_id=?", dtype.TypeId).Exec()
			}
			if err1 == nil && len(dtype.Fields) > 0 {
				_, err1 = tx.InsertMulti(30, dtype.Fields)
			}
			return err1
		})

	} else {
		//insert
		isExistSameName := db.QueryTable(table_type_name).Filter("type_name", dtype.TypeName).Exist()
		if isExistSameName {
			return errors.New("已存在相同名称的类型")
		}

		err = db.DoTx(func(ctx context.Context, tx orm.TxOrmer) error {
			_, err1 := tx.Insert(dtype)
			if err1 == nil && len(dtype.Fields) > 0 {
				for _, f := range dtype.Fields {
					f.TypeId = dtype.TypeId
				}
				if err1 == nil {
					_, err1 = tx.InsertMulti(30, dtype.Fields)
				}
			}
			return err1
		})
	}
	if err == nil {
		cache := fn.GetPublicCache()
		cache.Delete(type_map_cache_key)
	}
	return err
}

//查询类型是否已经有使用的账号
func HasAccountInType(accountTypeId int) bool {
	return orm.NewOrm().QueryTable("account").Filter("type_id", accountTypeId).Filter("is_deleted", 0).Exist()
}

func DelAccountType(accountTypeId int) error {
	if HasAccountInType(accountTypeId) {
		return errors.New("不可删除，类型已有关联账号")
	}
	_, err := orm.NewOrm().Delete(&AccountType{TypeId: accountTypeId})
	cache := fn.GetPublicCache()
	cache.Delete(type_map_cache_key)
	return err
}
