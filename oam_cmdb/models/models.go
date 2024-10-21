package models

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego/logs"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/validation"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

type CallResult struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

//js组件用数据结构
type JsData struct {
	Id   int
	Text string
}

//构建一个成功结果封闭
func (r *CallResult) Ok(data interface{}) {
	r.Status = 200
	r.Message = "success"
	r.Data = data
}

func (r *CallResult) ParamError(message string) {
	r.Status = 400
	if message == "" {
		message = "参数错误"
	}
	r.Message = message
}

func (r *CallResult) Failed(message string) {
	r.Status = 0
	r.Message = message
}

func (r *CallResult) Forbidden() {
	r.Status = 403
	r.Message = "无权限操作"
}

func (r *CallResult) Error() {
	r.Status = 500
	r.Message = "执行错误"
}

//检查查询结果并记录错误日志. 如果有error返回true
func checkQueryErr(err error) bool {
	if err != nil {
		if orm.ErrNoRows != err {
			logs.Error(err)
		}
		return true
	}
	return false
}

//把所有验证错误拼接为一个错误消息
func ToErrMsg(valid validation.Validation) string {
	var errmsgs []string
	if valid.HasErrors() {
		for _, er := range valid.Errors {
			errmsgs = append(errmsgs, er.Message)
		}
	}

	return strings.Join(errmsgs, ";")
}

func init() {
	var err error
	orm.Debug, err = beego.AppConfig.Bool("debug")
	if err != nil {
		fmt.Println(err)
	}
	orm.RegisterModel(new(UserInfo), new(Account), new(AccountField), new(AccountType), new(Project), new(Document),
		new(DictItem), new(Host), new(AppInfo), new(Role), new(Fun))

	dbtype := beego.AppConfig.DefaultString("dbtype", "mysql")
	if dbtype == "mysql" {
		dbuser, _ := beego.AppConfig.String("dbuser")
		dbpwd, _ := beego.AppConfig.String("dbpasswd")
		dbaddr, _ := beego.AppConfig.String("dbaddr")
		dbname, _ := beego.AppConfig.String("dbname")
		dbUrlTpl := "%s:%s@tcp(%s)/%s?charset=utf8&loc=Local"
		var dbUrl = fmt.Sprintf(dbUrlTpl, dbuser, dbpwd, dbaddr, dbname)
		orm.RegisterDataBase("default", "mysql", dbUrl)
	} else if dbtype == "sqlite" {
		orm.RegisterDriver("sqlite", orm.DRSqlite)
		orm.RegisterDataBase("default", "sqlite3", "data/oam.db")
	} else {
		panic("不支持的数据库类型")
	}
	//表单验证信息重置为中文(框架有点Low)
	SetDefaultMessage()
}

var MessageTmpls = map[string]string{
	"Required":     "不能为空",
	"Min":          "最小为 %d",
	"Max":          "最大为 %d",
	"Range":        "范围在 %d 至 %d",
	"MinSize":      "最小长度为 %d",
	"MaxSize":      "最大长度为 %d",
	"Length":       "长度必须是 %d",
	"Alpha":        "必须是有效的字母字符",
	"Numeric":      "必须是有效的数字字符",
	"AlphaNumeric": "必须是有效的字母或数字字符",
	"Match":        "格式不正确 %s",
	"NoMatch":      "必须不匹配格式 %s",
	"AlphaDash":    "必须是有效的字母或数字或破折号(-_)字符",
	"Email":        "不是有效的邮件地址",
	"IP":           "必须是有效的IP地址",
	"Base64":       "必须是有效的base64字符",
	"Mobile":       "不是有效的手机号码",
	"Tel":          "不是有效电话号码",
	"Phone":        "必须是有效的电话号码或者手机号码",
	"ZipCode":      "不是有效的邮政编码",
}

//默认设置通用的错误验证和提示项
func SetDefaultMessage() {
	if len(MessageTmpls) == 0 {
		return
	}
	//将默认的提示信息转为自定义
	for k, v := range MessageTmpls {
		validation.MessageTmpls[k] = v
	}

	//增加默认的自定义验证方法
	//_ = validation.AddCustomFunc("Unique", Unique)
}

func ToIntSlice(list orm.ParamsList) []int {
	var ids []int
	for _, id := range list {
		elem := id.(string) //ParamsList内容是string
		elemInt, _ := strconv.Atoi(elem)
		ids = append(ids, elemInt)
	}
	return ids
}

// orm.Params转为map[string]int
func ToIntValueMap(paramMap orm.Params) map[string]int {
	if len(paramMap) == 0 {
		return nil
	}
	intValueMap := make(map[string]int, len(paramMap))
	var tmpVal int
	var err error
	for key, value := range paramMap {
		tmpVal, err = strconv.Atoi(value.(string))
		if err == nil {
			intValueMap[key] = tmpVal
		} else {
			panic(err)
		}
	}
	return intValueMap
}

// orm.Params转为map[int]string
func ToIntKeyMap(paramMap orm.Params) map[int]string {
	if len(paramMap) == 0 {
		return nil
	}
	intValueMap := make(map[int]string, len(paramMap))
	var tmpKey int
	var err error
	for key, value := range paramMap {
		tmpKey, err = strconv.Atoi(key)
		if err == nil {
			intValueMap[tmpKey] = value.(string)
		} else {
			panic(err)
		}
	}
	return intValueMap
}
