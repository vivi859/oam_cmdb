package controllers

import (
	"OAM/conf"
	"OAM/models"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/beego/beego/v2/core/logs"
)

type DocController struct {
	AuthController
}

func (c *DocController) ToAddDoc() {
	projId, _ := c.GetInt("projId", 0)
	doc := models.Document{ProjId: projId}
	c.Data["doc"] = doc

	projs := models.FindProjectForMap()
	c.Data["projs"] = projs

	c.TplName = "doc-edit.html"
}

func (c *DocController) ToEditDoc() {
	docId, err := c.GetInt("docId")
	if err != nil {
		c.toErrorPage(err.Error())
		return
	}
	doc := models.GetDocumentById(docId)
	if doc == nil {
		c.toErrorPage("文档不存在")
		return
	}
	if doc.DocType != "md" {
		c.toErrorPage("不支持编辑该类型文档")
		return
	}
	projs := models.FindProjectForMap()
	c.Data["projs"] = projs

	c.Data["doc"] = doc
	c.TplName = "doc-edit.html"
}

func (c *DocController) View() {
	idstr, ok := c.Ctx.Input.Params()["0"]
	if !ok {
		c.toErrorPage("参数错误")
		return
	}
	docId, err := strconv.Atoi(idstr)
	if err != nil {
		c.toErrorPage("参数错误")
		return
	}
	doc := models.GetDocumentById(docId)
	if doc == nil {
		c.toErrorPage("文档不存在")
		return
	}

	c.Data["doc"] = doc
	if doc.DocType == "md" {
		c.TplName = "doc-view.html"
	} else {
		//c.outdoc(*doc)
		//c.TplName = "doc-preview.html"
		err = c.Preview(doc)
		if err != nil {
			c.toErrorPage(err.Error())
		}
	}
}

func (c *DocController) Preview(doc *models.Document) error {
	showName := url.QueryEscape(doc.Title)
	output := c.Ctx.Output
	switch doc.DocType {
	case "pdf":
		output.Header("Content-Type", "application/pdf")
		output.Header("Content-Disposition", "inline;fileName="+showName+";fileName*=UTF-8''"+showName)
	case "doc":
		showName = showName + ".doc"
		output.Header("Content-Type", "application/msword")
		output.Header("Content-Disposition", "attachment; fileName="+showName+";fileName*=UTF-8''"+showName)
	case "docx":
		showName = showName + ".docx"
		output.Header("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
		output.Header("Content-Disposition", "attachment; fileName="+showName+";fileName*=UTF-8''"+showName)
	default:
		return errors.New("不支持的文档格式")
	}

	var content []byte
	var err error
	content, err = ioutil.ReadFile(filepath.Join(conf.GlobalCfg.DATA_STORE_PATH, doc.Content))

	if err != nil {
		http.Error(output.Context.ResponseWriter, "读取文档失败，文档不存在或已损坏", http.StatusInternalServerError)
		return err
	}
	return output.Body(content)
}

func (c *DocController) SaveDoc() {
	c.justPost()
	doc := models.Document{}
	c.BindJSON(&doc)
	if len(doc.Title) > 100 {
		c.JsonParamError("标题长度小于100")
	}
	if doc.Title == "" {
		doc.Title = "无标题"
	}
	curUser := c.getLoginUser()
	doc.UpdateBy = curUser.UserName
	var err error
	if doc.DocId > 0 {
		oldDoc := models.GetDocumentByIdNoContent(doc.DocId)
		if oldDoc == nil {
			c.JsonParamError("文档不存在")
			return
		}
		doc.DocType = oldDoc.DocType
		doc.CreateBy = oldDoc.CreateBy
		_, err = models.UpdateDocument(&doc)

	} else {
		doc.DocType = "md"
		doc.AuthorId = curUser.UserId
		if curUser.RealName != "" {
			doc.CreateBy = curUser.RealName
		} else {
			doc.CreateBy = curUser.UserName
		}
		_, err = models.SaveDocument(&doc)
	}
	if err != nil {
		logs.Error("保存文档失败", err)
		c.JsonFailed("保存失败")
	} else {
		doc.Content = ""
		c.JsonOk(doc)
	}
}

func (c *DocController) ImportDoc() {
	projId, _ := c.GetInt("projId", 0)
	t := c.GetString("docTitle")
	if t == "" {
		t = "无标题"
	}
	docFile, err := c.uploadFile("docFile", 1)
	if err != nil {
		c.JsonFailed(err.Error())
	}

	docType := path.Ext(docFile)[1:]
	curUser := c.getLoginUser()

	var doc = models.Document{ProjId: projId, Title: t, Content: docFile, DocType: docType, UpdateBy: curUser.UserName, AuthorId: curUser.UserId}
	if curUser.RealName != "" {
		doc.CreateBy = curUser.RealName
	} else {
		doc.CreateBy = curUser.UserName
	}
	//md文件改为内容存数据库字段
	if docType == "md" {
		var docFileAbsPath = filepath.Join(conf.GlobalCfg.DATA_STORE_PATH, docFile)
		var contentBytes []byte
		contentBytes, err = ioutil.ReadFile(docFileAbsPath)
		if err != nil {
			c.JsonError()
		}
		doc.Content = string(contentBytes)
		os.Remove(docFileAbsPath)
	}
	_, err = models.SaveDocument(&doc)
	if err != nil {
		logs.Error("导入文档失败", err)
		c.JsonFailed("导入失败")
	} else {
		doc.Content = ""
		c.JsonOk(doc)
	}
}

func (c *DocController) DeleteDoc() {
	docId, err := c.GetInt("docId")
	if err != nil {
		c.JsonParamError("参数错误")
	}
	curUser := c.getLoginUser()
	err = models.DeleteDocumentById(docId)
	if err == nil {
		logs.Info("[%s]删除文档id:%v", curUser.UserName, docId)
		c.JsonOk(docId)
	} else {
		c.JsonFailed(err.Error())
	}
}

func (c *DocController) ToDocList() {
	projs := models.FindProjectForMap()
	c.Data["projItems"] = projs
	c.TplName = "doc.html"
}

func (c *DocController) DocPage() {
	where := make(map[string]interface{})
	ip := c.GetString("keyword")
	if ip != "" {
		where["keyword"] = ip
	}
	pid, err := c.GetInt("projId")
	if err == nil {
		where["projId"] = pid
	}
	row, _ := c.GetInt("rows", models.DEFAULT_PAGE_SIZE)
	page, _ := c.GetInt("page", 1)
	pageData := models.FindDocForPage(row, page, where)
	c.JsonOk(pageData)
}
