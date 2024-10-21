package models

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

const DEFAULT_PAGE_SIZE = 15

// 分页数据结构
type Page[T any] struct {
	//每页行数
	RowPerPage int
	//总行数
	TotalRow int
	//总页数
	TotalPage int
	//当前页码
	PageNum int
	//当前页数据
	Rows []*T

	// 查询总行数方法
	queryTotalRow func() int `json:"-"`
	//查询当前页数据方法
	queryPageData func(pageNum int, pageSize int) []*T `json:"-"`
}

func (p *Page[T]) countTotalPage() {
	tp := p.TotalRow / p.RowPerPage
	if p.TotalRow%p.RowPerPage > 0 {
		tp = tp + 1
	}
	p.TotalPage = tp
}

//使用原生sql分页
func (p *Page[T]) RawQueryPage(sql_count string, sql_queryrow string, params []interface{}) {
	query := orm.NewOrm()
	p.queryTotalRow = func() int {
		var c int
		err := query.Raw(sql_count).SetArgs(params...).QueryRow(&c)
		if err == nil {
			return c
		} else {
			logs.Error(err)
			return 0
		}
	}
	p.queryPageData = func(pageNum int, pageSize int) []*T {
		var rows []*T
		newparams := append(params, pageSize*(pageNum-1), pageSize)
		query.Raw(sql_queryrow + " limit ?,?").SetArgs(newparams...).QueryRows(&rows)
		return rows
	}
	p.doQueryPage()
}

//使用QuerySeter分页查询,适合单表
func (p *Page[T]) QueryPage(selector orm.QuerySeter) {
	p.queryTotalRow = func() int {
		c, err := selector.Count()
		if err == nil {
			return int(c)
		} else {
			logs.Error(err)
			return 0
		}
	}
	p.queryPageData = func(pageNum int, pageSize int) []*T {
		var rows []*T
		selector.Limit(pageSize, pageSize*(pageNum-1)).All(&rows)
		return rows
	}
	p.doQueryPage()
}

func (p *Page[T]) doQueryPage() {
	if p.PageNum == 0 {
		p.PageNum = 1
	}
	if p.RowPerPage == 0 {
		p.RowPerPage = 20
	}
	if p.queryTotalRow != nil {
		p.TotalRow = p.queryTotalRow()
	}
	if p.TotalRow > 0 {
		p.countTotalPage()
		if p.queryPageData != nil {
			p.Rows = p.queryPageData(p.PageNum, p.RowPerPage)
		}
	}
}
