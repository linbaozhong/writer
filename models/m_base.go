package models

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"math"
)

const (
	ReviewQuestion = 0
)

//页面公共信息
type Page struct {
	SiteName    string //网站名称
	Title       string //页面标题
	Company     string //公司名称
	Domain      string //域名
	Copyright   string //版权
	Keywords    string //Seo关键词
	Description string //Seo描述
	Author      string //作者
	Product     string //产品名称
	Version     string //版本
}
type Current struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	From   string `json:"from"`
	Role   string `json:"role"`
}

//公共字段
type Field struct {
	Sequence int
	Status   int8
	Deleted  int8
	Updator  int64
	Updated  int64
	Ip       string
	Name     string
	Role     string
}

//错误信息
type Error struct {
	Key     string
	Message string
}

//分页
type Pagination struct {
	Count int
	Prev  int
	Index int
	Next  int
	Size  int
}

//列表选项
type SelectItem struct {
	Key      string
	Value    string
	Selected bool
}

//上传文件
type UploadFile struct {
	Name string
	Ext  string
	Path string
	Size int64
}

var db *xorm.Engine

func init() {
	var err error
	db, err = xorm.NewEngine("mysql", "root:goldapple@127.0.0.1:3306/writer?charset=utf8")

	if err != nil {
		beego.Trace(err)
	}

	db.SetMapper(core.SameMapper{})

	db.ShowInfo = true
	//db.ShowDebug = true
	db.ShowSQL = true
	//db.ShowErr = true
	//db.ShowWarn = true
}

/*
* 数据有效性检查
 */
func dataCheck(d interface{}) ([]Error, error) {
	//数据有效性检验
	valid := validation.Validation{}
	b, err := valid.Valid(d)

	if err != nil {
		return nil, err
	}
	if !b {
		// 整理错误信息
		es := make([]Error, 0)

		for _, err := range valid.Errors {
			es = append(es, Error{Key: err.Key, Message: err.Message})
			beego.Error("无效数据：%s-%s", err.Key, err.Message)
		}
		return es, errors.New("无效数据")
	}
	return nil, nil
}

//xorm的补充
func parseDb(dbs []map[string][]byte) []map[string]string {
	_st := make([]map[string]string, 0)
	for _, value := range dbs {
		_mt := make(map[string]string)
		for k, v := range value {
			_mt[k] = string(v)
		}
		_st = append(_st, _mt)
	}
	return _st
}

// 根据记录总数，返回总页数
func getPageCount(rows int64, page *Pagination) {
	page.Count = int(math.Ceil(float64(rows / int64(page.Size))))
}
