package models

import (
	fn "OAM/util"
	"errors"
	"regexp"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
)

// const SUPER_ADMIN string = "root"
const (
	Diabled = iota
	Enabled
	Locked
)

// 用户名正则
const USER_NAME_PATTERN = `^[\w@.]{3,30}$`
const USER_PWD_PATTERN = `^[\S]{6,30}$`

// 用户信息
type UserInfo struct {
	UserId     int    `orm:"PK;auto" form:"userId"`
	UserName   string `valid:"Required" form:"userName"`
	Passwd     string `valid:"Required" form:"passwd" json:"-"`
	RealName   string `form:"realName" valid:"MaxSize(30)"`
	UserStatus int8   `form:"userStatus"`
	CreateBy   string
	CreateTime time.Time `orm:"auto_now_add;type(datetime)" json:",omitempty"`
	UpdateBy   string
	UpdateTime time.Time `orm:"auto_now;type(datetime)" json:",omitempty"`
	RoleCode   string    `form:"roleCode"`
	RePasswd   string    `orm:"-" form:"rePasswd"`
}

func (u UserInfo) ToLoginUser() LoginUser {
	luser := LoginUser{
		UserId:     u.UserId,
		UserName:   u.UserName,
		RealName:   u.RealName,
		RoleCode:   u.RoleCode,
		UserStatus: u.UserStatus,
	}
	return luser
}

const table_user_name = "user_info"

func (t *UserInfo) TableName() string {
	return table_user_name
}

func (user *UserInfo) Valid() string {
	valid := validation.Validation{}
	_, err := valid.Valid(user)
	if err != nil {
		return err.Error()
	}
	if isValid, err := regexp.MatchString(USER_NAME_PATTERN, user.UserName); err != nil && !isValid {
		valid.SetError("UserName", "用户名不正确")
	}
	if isValid, err := regexp.MatchString(USER_PWD_PATTERN, user.Passwd); err != nil && !isValid {
		valid.SetError("Passwd", "密码不正确")
	}
	if user.Passwd != user.RePasswd {
		valid.SetError("Passwd", "两次密码不相同")
	}
	return ToErrMsg(valid)
}

// 验证密码
func (u UserInfo) CheckPasswd(validpasswd string) bool {
	return u.Passwd == encryptPasswd(validpasswd)
}

func (u *UserInfo) EncryptPasswd() {
	u.Passwd = encryptPasswd(u.Passwd)
}

// 密码加密
func encryptPasswd(pwd string) string {
	return fn.SHA256Hex(pwd + "&&")
}

/*func (u UserInfo) IsSuperAdmin() bool {
	return u.UserName == SUPER_ADMIN
}*/

// 根据ID查用户,如果不存在返回nil
func GetUserById(userId int) *UserInfo {
	query := orm.NewOrm()
	user := &UserInfo{UserId: userId}
	err := query.Read(user)
	if checkQueryErr(err) {
		return nil
	}

	return user
}

// 根据用户名查询用户
func GetUserByName(userName string) *UserInfo {
	u := new(UserInfo)
	err := orm.NewOrm().QueryTable("user_info").Filter("user_name", userName).One(u)
	if checkQueryErr(err) {
		return nil
	}
	return u
}

// 按条件查询用户
//
// 参数condi查询条件可用key说明: status状态, userName用户名(模糊匹配)
func FindUserByCondi(condi map[string]interface{}) []*UserInfo {
	var users []*UserInfo
	querySeter := orm.NewOrm().QueryTable("user_info")
	where := orm.NewCondition()
	if condi != nil {
		status, exist := condi["status"]
		if exist {
			where = where.And("user_status", status.(int))
		}
		name, exist := condi["keyword"]
		if exist {
			//where.And("user_name__icontains", name).Or("real_name__icontains",name)
			where = where.AndCond(where.And("user_name__icontains", name).Or("real_name__icontains", name))
		}
		role, exist := condi["roleCode"]
		if exist {
			where = where.And("role_code", role)
		}
		ids, exist := condi["ids"]
		if exist {
			userIds := fn.ToInterfaceSlice(ids.([]int))
			where = where.And("user_id__in", userIds...)
		}
		querySeter = querySeter.SetCond(where)
	}
	_, err := querySeter.Limit(-1).All(&users)
	if err != nil {
		logs.Error(err)
	}
	return users
}

// 查询项目成员用户
func FindUserByProjId(projId int) []*UserInfo {
	var users []*UserInfo
	var querySql = `select u.user_id,u.user_name,u.real_name,u.user_status from user_info u join 
	rel_proj_user r on r.user_id=u.user_id where proj_id=?`
	_, err := orm.NewOrm().Raw(querySql, projId).QueryRows(&users)
	if checkQueryErr(err) {
		return nil
	}
	return users
}

// 查询非项目成员用户
func FindNotProjMember(projId int) []*UserInfo {
	var users []*UserInfo
	var querySql = `select u.user_id,u.user_name,u.real_name from user_info u where u.user_status=1 and 
	not EXISTS (select 1 from rel_proj_user r where r.proj_id=? and r.user_id=u.user_id)`
	_, err := orm.NewOrm().Raw(querySql, projId).QueryRows(&users)
	if checkQueryErr(err) {
		return nil
	}
	return users
}

func SaveUser(user *UserInfo) (int64, error) {
	db := orm.NewOrm()
	isExist := db.QueryTable("user_info").Filter("user_name", user.UserName).Exist()
	if isExist {
		return 0, errors.New("用户名已存在")
	}
	if user.Passwd != "" {
		user.EncryptPasswd()
	}
	user.UserStatus = Enabled
	return db.Insert(user)
}

func UpdateUser(user *UserInfo) (int64, error) {
	oldUser := GetUserById(user.UserId)
	if oldUser == nil {
		return 0, errors.New("用户不存在")
	}
	user.UpdateTime = time.Now()
	db := orm.NewOrm()
	return db.Update(user, "real_name", "user_status", "role_code", "update_by", "update_time")
}

func UpdateUserPasswd(user *UserInfo) error {
	oldUser := GetUserById(user.UserId)
	if oldUser == nil {
		return errors.New("用户不存在")
	}
	if user.Passwd == "" || user.Passwd != user.RePasswd {
		return errors.New("密码错误")
	}
	user.EncryptPasswd()
	user.UpdateTime = time.Now()
	db := orm.NewOrm()
	_, err := db.Update(user, "passwd", "update_by", "update_time")
	return err
}
