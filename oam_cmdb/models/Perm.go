package models

import (
	fn "OAM/util"
	"errors"
	"regexp"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	"golang.org/x/exp/slices"
)

const (
	api_pattern = "^(/[a-zA-Z0-9_\\-\\.]+\\/?)+\\*?$"
	//allow_http_methods=["get","post","put","delete"]
)

type UrlPerm struct {
	//methods  []string
	wildcard bool
	path     string
}

type LoginUser struct {
	UserId     int
	UserName   string
	RealName   string
	UserStatus int8
	RoleCode   string
	FunCodes   []string
	UrlPerms   []UrlPerm
}

func ValidUrl(url string) bool {
	matched, _ := regexp.MatchString(api_pattern, url)
	return matched
}

func NewUrlPerm(url string) (UrlPerm, error) {
	var perm UrlPerm
	if len(url) == 0 {
		return perm, errors.New("url must be not empty")
	}

	if !ValidUrl(url) {
		return perm, errors.New("url perm format error")
	}

	url = strings.TrimSuffix(url, "/")
	if strings.HasSuffix(url, "/*") {
		perm.wildcard = true
		perm.path = strings.TrimSuffix(url, "/*")
	} else {
		perm.wildcard = false
		perm.path = url
	}

	return perm, nil
}

func (p UrlPerm) Implies(reqUrl string) bool {
	reqUrl = strings.TrimSuffix(reqUrl, "/")
	if p.path == reqUrl {
		return true
	}
	if p.wildcard {
		if strings.HasPrefix(reqUrl, p.path) {
			//*只支持一级目录匹配,如果后面还有目录不匹配
			if strings.LastIndex(reqUrl, "/") < len(p.path) {
				return true
			}
		}
	}

	return false
}

//创建登录用户,加载用户角色权限
func CreateLoginUser(user UserInfo) LoginUser {
	loginUser := user.ToLoginUser()
	if user.RoleCode == ROLE_ROOT {
		return loginUser
	}
	roleFuns := FindFunByRoleCode(user.RoleCode, 0)
	if len(roleFuns) > 0 {
		var funCodes []string
		var perms []UrlPerm
		//组装url权限
		for _, fun := range roleFuns {
			funCodes = append(funCodes, fun.FunCode)
			if len(fun.FunUrl) > 0 {
				urls := fn.Split(fun.FunUrl)
				for _, url := range urls {
					p, err := NewUrlPerm(url)
					if err != nil {
						logs.Warn("urlperm format error:" + url + ",err:" + err.Error())
						continue
					}
					perms = append(perms, p)
				}
			}
		}
		loginUser.FunCodes = funCodes
		loginUser.UrlPerms = perms
	}
	return loginUser
}

//判断登录用户角色
func (user LoginUser) HasRole(roleCode string) bool {
	return user.RoleCode == roleCode
}

//判断是否为超级管理员角色
func (user LoginUser) IsSupervisor() bool {
	return user.RoleCode == ROLE_ROOT
}

//基于功能标识判断登录用户是否有权限
func (user LoginUser) HasPermByFunCode(funCode string) bool {
	if user.IsSupervisor() {
		return true
	}
	return len(user.FunCodes) > 0 && slices.Contains[string](user.FunCodes, funCode)
}

func (user LoginUser) HasPermByFunCodes(funCodes ...string) map[string]bool {
	result := make(map[string]bool)
	if user.IsSupervisor() {
		for _, f := range funCodes {
			result[f] = true
		}
	} else {
		for _, f := range funCodes {
			result[f] = slices.Contains[string](user.FunCodes, f)
		}
	}
	return result
}

//基于接口url判断登录用户是否有权限
func (user LoginUser) HasPermByUrl(reqUrl string) bool {
	if user.RoleCode == ROLE_ROOT {
		return true
	}
	if len(user.UrlPerms) == 0 {
		return false
	}

	for _, p := range user.UrlPerms {
		if p.Implies(reqUrl) {
			return true
		}
	}
	return false
}

/*基于权限标识的判断,类似shiro
type Perm struct {
	parts []sets.Set
}

func NewPerm(wildcardStr string) Perm {
	wildcardStr = strings.ToLower(util.Trim(wildcardStr))
	tmpParts := strings.Split(wildcardStr, ":")
	var p Perm
	for _, part := range tmpParts {
		subParts := hashset.New()
		subParts.Add(strings.Split(part, ","))
		p.parts = append(p.parts, subParts)
	}
	if p.parts == nil {
		panic("")
	}
	return p
}

//授权判断
func (p Perm) Implies(np Perm) bool {
	i := 0
	plen := len(p.parts)
	for _, op := range np.parts {
		if plen-1 < i {
			return true
		} else {
			if !p.parts[i].Contains("*") && !p.parts[i].Contains(op) {
				return false
			}
			i++
		}
	}

	for _, op := range p.parts {
		if !op.Contains("*") {
			return false
		}
	}
	return true
}
*/
