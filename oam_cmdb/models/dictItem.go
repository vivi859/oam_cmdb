package models

import (
	"errors"
	"strconv"
	"strings"
	"time"

	fn "OAM/util"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
	jsoniter "github.com/json-iterator/go"
	"github.com/sony/sonyflake"
)

const (
	DICTKEY_SYMMETRIC         = "symmetric_key"
	DICTKEY_OS_NAMES          = "os_names"
	DICTKEY_SERVICE_SOFTWARES = "service_softwares"
	DICTKEY_DEV_LANG          = "dev_lang"
	DICTKEY_APP_TYPE          = "app_type"
)

type DictItem struct {
	ItemId    string `orm:"PK" form:"itemId"`
	ItemName  string `form:"itemName" valid:"Required; MaxSize(100)" label:"参数名"`
	ItemValue string `form:"itemValue" valid:"MaxSize(5000)" label:"参数值"`
	ItemType  int    `form:"itemType"`
	//ItemGroup  int
	UpdateBy   string
	UpdateTime time.Time `orm:"auto_now;type(datetime)"`
}

func (t *DictItem) TableName() string {
	return "dict_item"
}

func (acct *DictItem) Valid() string {
	valid := validation.Validation{}
	_, err := valid.Valid(acct)
	if err != nil {
		return err.Error()
	}

	return ToErrMsg(valid)
}

func GetDictItemById(id string) *DictItem {
	item := DictItem{ItemId: id}
	query := orm.NewOrm()
	err := query.Read(&item)
	if checkQueryErr(err) {
		return nil
	}

	return &item
}

// 读取aes密钥,如果不存在初始化一个
func InitSymmetrickey() string {
	item := DictItem{ItemId: DICTKEY_SYMMETRIC}
	db := orm.NewOrm()
	err := db.Read(&item)
	if err != nil && err != orm.ErrNoRows {
		panic("读取初始密钥失败")
	}
	var key string
	if err == orm.ErrNoRows || item.ItemValue == "" {
		key = fn.RandomNumStr(32)
		var err1 error
		encryptKey, err1 := fn.RSAEncryptBase64Str(key, "")
		if err1 != nil {
			logs.Error("生成初始密钥失败", err)
			panic("生成初始密钥失败")
		}
		item.ItemValue = encryptKey

		if err == orm.ErrNoRows {
			item.UpdateBy = "sys"
			item.ItemName = "初始密钥"
			_, err1 = db.Insert(&item)
		} else {
			_, err1 = db.Update(&item)
		}

		if err1 != nil {
			logs.Error("初始密钥保存失败", err)
			panic("初始密钥保存失败")
		}
	} else {
		key, err = fn.RSADecryptBase64Str(item.ItemValue, "")
		if err != nil {
			logs.Error("解密初始密钥失败", err)
			panic("解密初始密钥失败")
		}
	}

	return key
}

var startime = time.Date(2023, time.November, 1, 1, 1, 1, 0, time.Local)

var FlakeIDGen = sonyflake.NewSonyflake(sonyflake.Settings{StartTime: startime})

func UpdateDictItem(item *DictItem) (int64, error) {
	if len(item.ItemId) == 0 {
		id, err := FlakeIDGen.NextID()
		if err != nil {
			return 0, err
		}
		item.ItemId = strconv.FormatUint(id, 10)
		item.fixValue()
		return orm.NewOrm().Insert(item)
	} else {
		old := GetDictItemById(item.ItemId)
		if old == nil {
			return 0, errors.New("参数不存在")
		}
		item.ItemType = old.ItemType
		item.fixValue()
		cache().Delete(item.ItemId)
		return orm.NewOrm().Update(item)
	}
}

func (item *DictItem) fixValue() {
	if len(item.ItemValue) > 0 {
		if item.ItemType == 3 {
			arr := strings.Split(strings.TrimSpace(item.ItemValue), "\r\n")
			for i := 0; i < len(arr); i++ {
				arr[i] = strings.TrimSpace(arr[i])
			}
			item.ItemValue = fn.Join(arr)
		}
	}
}

func FindAllDictItem() []*DictItem {
	var datums []*DictItem
	orm.NewOrm().QueryTable("dict_item").All(&datums)
	return datums
}

//查询默认设置的操作系统名
func GetPresetOS() []string {
	item := getDictCacheFirst(DICTKEY_OS_NAMES)
	return fn.Split(item.ItemValue)
}

// 查询数据字典,值用逗号拆分作为数组返回
func GetDictValueAsArray(itemId string) []string {
	item := getDictCacheFirst(itemId)
	return fn.Split(item.ItemValue)
}

// 获取预定义的常用软件列表
func GetPresetServiceSoftwares() []AppInfo {
	var apps []AppInfo
	item := getDictCacheFirst(DICTKEY_SERVICE_SOFTWARES)
	if item.ItemValue != "" {
		err := jsoniter.UnmarshalFromString(item.ItemValue, &apps)
		if err != nil {
			logs.Error("解析预定义软件列表异常:{}", err)
		}
	}
	return apps
}

func getDictCacheFirst(itemId string) DictItem {
	item, err := getCache(itemId)
	if err != nil {
		item = *GetDictItemById(itemId)
		putCache(item)
	}
	return item
}

func cache() fn.MyCache {
	return fn.GetCache(fn.CACHE_DICT)
}

func putCache(item DictItem) {
	cache().PutWithExpireTime(item.ItemId, item, time.Minute*60*72)
}

func getCache(itemId string) (DictItem, error) {
	v, err := cache().Get(itemId)
	if err == nil {
		return v.(DictItem), nil
	}
	return DictItem{}, err
}
