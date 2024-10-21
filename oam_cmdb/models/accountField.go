package models

import (
	"path/filepath"
	"strconv"

	"github.com/beego/beego/v2/client/orm"
	jsoniter "github.com/json-iterator/go"
)

/* 账号类型动态属性 */
type AccountField struct {
	FieldId      int `orm:"PK;auto"`
	FieldName    string
	FieldKey     string
	TypeId       int
	IsRequired   bool
	IsCiphertext bool
	MaxLen       int
	ValueType    uint16
	ValueRule    string
	FieldValue   string `orm:"-"`
	Sort         uint16
}

var FieldValueTypes = map[uint16]string{
	0: "文本", 1: "布尔", 2: "数字", 3: "可选值", 4: "文件",
}

const table_field_name = "account_field"

func (t *AccountField) TableName() string {
	return table_field_name
}

//获取账号类型下的所有动态属性
func FindFieldsByTypeId(accountTypeId int) []*AccountField {
	var fields []*AccountField
	orm.NewOrm().QueryTable(table_field_name).Filter("type_id", accountTypeId).OrderBy("sort").All(&fields)
	return fields
}

//根据id批量删除动态属性
func DelFieldByIds(fieldIds []int) (int64, error) {
	return orm.NewOrm().QueryTable(table_field_name).Filter("field_id__in", fieldIds).Delete()
}

//删除账号类型下的所有动态属性
func DelFieldByTypeId(typeId int) (int64, error) {
	return orm.NewOrm().QueryTable(table_field_name).Filter("type_id", typeId).Delete()
}

//值类型为可选值时返回FieldValue对应的Text，其他返回FieldValue
func (t AccountField) FieldTextValue() string {
	if len(t.FieldValue) == 0 {
		return t.FieldValue
	}
	if t.ValueType == 1 {
		b, err := strconv.ParseBool(t.FieldValue)
		if err != nil && b {
			return "是"
		}
		return "否"
	} else if t.ValueType == 3 {
		if len(t.ValueRule) > 0 {
			var rules []JsData
			err := jsoniter.UnmarshalFromString(t.ValueRule, &rules)
			if err == nil {
				var i int
				i, err = strconv.Atoi(t.FieldValue)
				if err == nil {
					for _, r := range rules {
						if r.Id == i {
							return r.Text
						}
					}
				}
			}
		}
	} else if t.ValueType == 4 {
		return filepath.Base(t.FieldValue)
	}
	return t.FieldValue
}
