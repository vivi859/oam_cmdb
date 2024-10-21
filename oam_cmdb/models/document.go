package models

import (
	"OAM/conf"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

/* 文档 */
type Document struct {
	DocId      int `orm:"PK;auto"`
	Title      string
	DocType    string
	ProjId     int
	AuthorId   int
	CreateBy   string
	UpdateBy   string
	CreateTime time.Time `orm:"auto_now_add;type(datetime)" json:",omitempty"`
	UpdateTime time.Time `orm:"auto_now;type(datetime)" json:",omitempty"`
	Content    string
}

const table_doc_name = "document"

func (t *Document) TableName() string {
	return table_doc_name
}

//查询项目下所有文档
func FindDocsByProjId(projId int) []*Document {
	var fields []*Document
	orm.NewOrm().QueryTable(table_doc_name).Filter("proj_id", projId).OrderBy("-doc_id").All(&fields)
	return fields
}

func GetDocumentById(id int) *Document {
	doc := Document{DocId: id}
	query := orm.NewOrm()
	err := query.Read(&doc)
	if checkQueryErr(err) {
		return nil
	}
	return &doc
}

// 根据id查询文档,但不返回内容,内容较大,只需要其他信息是使用
func GetDocumentByIdNoContent(id int) *Document {
	doc := Document{}
	query := orm.NewOrm()
	err := query.QueryTable(table_doc_name).Filter("doc_id", id).One(&doc, "doc_id", "title", "doc_type", "proj_id", "update_time", "create_by", "author_id")
	if checkQueryErr(err) {
		return nil
	}
	return &doc
}

func SaveDocument(doc *Document) (int64, error) {
	return orm.NewOrm().Insert(doc)
}

func UpdateDocument(doc *Document) (int64, error) {
	return orm.NewOrm().Update(doc, "title", "proj_id", "content", "update_by", "update_time")
}

// 删除文档
func DeleteDocumentById(id int) error {
	doc := Document{DocId: id}
	db := orm.NewOrm()
	err := db.Read(&doc)
	if err != nil {
		return err
	}

	_, err = orm.NewOrm().Delete(&doc)
	if err == nil {
		if doc.DocType != "md" {
			os.Remove(filepath.Join(conf.GlobalCfg.DATA_STORE_PATH, doc.Content))
		} else {
			//删除文章中的图片文件
			if len(doc.Content) > 20 {
				reg, _ := regexp.Compile(`\[.+\.\w{3,4}\]\(/` + conf.STATIC_EXT_BASE_URL + `(/[\w/]+\.\w{3,4})\)`)
				mdR := reg.FindAllStringSubmatch(doc.Content, -1)
				if len(mdR) > 0 {
					for _, v := range mdR {
						os.Remove(filepath.Join(conf.GlobalCfg.DATA_STORE_PATH, v[1]))
						logs.Debug("删除文档图片:%s,docId=%d", v[1], id)
					}
				}
			}
		}

	}
	return err
}

func FindDocForPage(row int, curpage int, condi map[string]interface{}) Page[Document] {
	pageData := Page[Document]{RowPerPage: row, PageNum: curpage}
	query := orm.NewOrm()
	selector := query.QueryTable(table_doc_name)
	if condi != nil {
		where := orm.NewCondition()
		kw, exist := condi["keyword"]
		if exist {
			where = where.And("title__icontains", kw.(string))
		}
		projId, exist := condi["projId"]
		if exist {
			where = where.And("proj_id", projId.(int))
		}

		selector = selector.SetCond(where)
	}
	pageData.QueryPage(selector.OrderBy("-doc_id"))

	return pageData
}
