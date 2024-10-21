package controllers

import (
	"OAM/conf"
	"OAM/models"
	"OAM/util"
	"errors"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	beecontext "github.com/beego/beego/v2/server/web/context"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/exp/slices"
)

type BaseController struct {
	beego.Controller
}

// 用jsoniter替换原生json
func (c *BaseController) ServeJSON(encoding ...bool) error {
	var (
		hasIndent = beego.BConfig.RunMode != beego.PROD
		//hasEncoding = len(encoding) > 0 && encoding[0]
	)
	data := c.Data["json"]
	output := c.Ctx.Output
	output.Header("Content-Type", "application/json; charset=utf-8")
	var content []byte
	var err error
	if hasIndent {
		content, err = jsoniter.MarshalIndent(data, "", "  ")
	} else {
		content, err = jsoniter.Marshal(data)
	}
	if err != nil {
		http.Error(output.Context.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return err
	}
	//if encoding {
	//	content = []byte(stringsToJSON(string(content)))
	//}
	return output.Body(content)
}

func (c *BaseController) BindJSON(obj interface{}) error {
	return jsoniter.Unmarshal(c.Ctx.Input.RequestBody, obj)
}

// 输出CallResult JSON
func (c *BaseController) JsonResult(r models.CallResult) {
	c.Data["json"] = r
	c.ServeJSON()
	c.StopRun()
}

func (c *BaseController) JsonStatus(status int, message string) {
	r := models.CallResult{Status: status, Message: message}
	c.JsonResult(r)
}

// 输出成功返回JSON,数据以CallResult封装
func (c *BaseController) JsonOk(data interface{}) {
	var r models.CallResult
	r.Ok(data)
	c.JsonResult(r)
}

// 输出参数错误JSON,数据以CallResult封装
func (c *BaseController) JsonParamError(message string) {
	var r models.CallResult
	r.ParamError(message)
	c.JsonResult(r)
}

func (c *BaseController) JsonFailed(message string) {
	var r models.CallResult
	r.Failed(message)
	c.JsonResult(r)
}

func (c *BaseController) JsonError() {
	var r models.CallResult
	r.Error()
	c.JsonResult(r)
}

func (c *BaseController) JsonForbidden() {
	var r models.CallResult
	r.Forbidden()
	c.JsonResult(r)
}

func (c *BaseController) toErrorPage(errmsg string) {
	c.Data["content"] = errmsg
	c.TplName = "error.html"
}

// 是否POST提交
func (c *BaseController) justPost() {
	if !c.Ctx.Input.IsPost() {
		c.Abort("405")
	}
}

var AllowedImgExts = []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".ico"}
var AllowedDocExts = []string{".pdf", ".doc", ".docx", ".md", ".xls", ".xlsx", ".ppt", ".pptx"}

// 图片文件最大限制
const img_max_size int64 = 1024 * 1024 * 5

// 文档文件最大限制
const doc_max_size int64 = 1024 * 1024 * 50

// 文件上传处理
func (c *BaseController) uploadFile(fileFormName string, fileType int) (string, error) {
	f, h, err := c.GetFile(fileFormName)
	if err != nil {
		logs.Error("上传图片异常", err)
		return "", err
	}

	defer f.Close()

	ext := path.Ext(h.Filename)
	if fileType == 0 {
		if !slices.Contains(AllowedImgExts, strings.ToLower(ext)) {
			return "", errors.New("不支持的图片类型")
		}
		if h.Size > img_max_size {
			return "", errors.New("超过图片限制大小：5M")
		}
	}
	if fileType == 1 {
		if !slices.Contains(AllowedDocExts, strings.ToLower(ext)) {
			return "", errors.New("不支持的文档类型")
		}
		if h.Size > doc_max_size {
			return "", errors.New("超过文件限制大小：50M")
		}
	}
	savePath := getSavePath(fileType)
	newFilePath := filepath.Join(savePath, h.Filename)
	if util.FileIsExists(newFilePath) {
		//newFileName := generateAttachName(ext, fileType)
		newFileName := util.RandomStr(5) + h.Filename
		newFilePath = filepath.Join(savePath, newFileName)
	}
	err = c.SaveToFile(fileFormName, filepath.Join(newFilePath))
	if err != nil {
		return "", err
	} else {
		relpath, err := filepath.Rel(conf.GlobalCfg.DATA_STORE_PATH, newFilePath)
		if err == nil && runtime.GOOS == "windows" {
			relpath = filepath.ToSlash(relpath)
		}
		return relpath, err
	}
}

func getSavePath(fileType int) string {
	dir := strconv.Itoa(int(time.Now().Month()))
	if fileType == 0 {
		dir = "/images/" + dir
	} else {
		dir = "/doc/" + dir
	}
	absDir := filepath.Join(conf.GlobalCfg.DATA_STORE_PATH, dir)
	if !util.FileIsExists(absDir) {
		err := os.MkdirAll(absDir, os.ModeDir|os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	return absDir
}

func generateAttachName(ext string, fileType int) string {
	name := strconv.FormatInt(time.Now().UnixMilli(), 10) + util.RandomStr(5) + ext
	return name
}

type AuthController struct {
	BaseController
}

// 登录过滤器排除路径
var loginAuthFilterExcludeUrls = []string{"/login"}

const LOGIN_USER_KEY = "currentLoginUser"

// 登录拦截器
var LoginAuthFilter = func(ctx *beecontext.Context) {
	reqPath := ctx.Request.URL.Path
	if !slices.Contains[string](loginAuthFilterExcludeUrls, reqPath) {
		loginUser, ok := ctx.Input.Session(LOGIN_USER_KEY).(models.LoginUser)
		if ok {
			//登录用户信息放入ctx
			ctx.Input.SetData(LOGIN_USER_KEY, loginUser)
		} else {
			logs.Warn("get loginuser fail,requestURI: %s,remoteAddr:%s", ctx.Request.RequestURI, ctx.Input.IP())
			if ctx.Input.IsAjax() {
				ctx.Abort(401, "Unauthorized")
			} else {
				ctx.Redirect(302, "/login")
			}
		}
	}
}

// 从web context中获取当前登录用户信息
func (c *AuthController) getLoginUser() models.LoginUser {
	return c.Ctx.Input.GetData(LOGIN_USER_KEY).(models.LoginUser)
}

var accessFilterExcludeUrls = []string{"/login", "/logout", "/index", "/pubkey", "/prikey", "/uploadimg", "/uploadfile", "/user/profile",
	"/user/updateloginuserpasswd", "/fun/navmenu", "/fun/getloginuserperms", "/error/403"}

// 权限控制拦截器,放在登录拦截后
var AccessFilter = func(ctx *beecontext.Context) {
	reqPath := ctx.Request.URL.Path

	if !slices.Contains[string](accessFilterExcludeUrls, reqPath) {
		loginUser := ctx.Input.GetData(LOGIN_USER_KEY).(models.LoginUser)
		if !loginUser.IsSupervisor() && !loginUser.HasPermByUrl(reqPath) {
			logs.Warn("user [%s] access forbidden,requestURI: %s,remoteAddr:%s", loginUser.UserName, ctx.Request.RequestURI, ctx.Input.IP())
			if ctx.Input.IsAjax() {
				ctx.Abort(403, "无权限访问")
			} else {
				ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
				ctx.WriteString("<p style='text-align:center'>没有权限访问该页面</p>")
			}
		}
	}
}
