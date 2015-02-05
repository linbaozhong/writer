package models

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	//_ "github.com/go-sql-driver/mysql"
	"fmt"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
	"math"
	"strings"
	"zouzhe/utils"
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
	Count int //总页数
	Prev  int //上页索引
	Index int //当前页
	Next  int //下页索引
	Size  int //每页条数
}

//列表选项
type SelectItem struct {
	Key      string
	Value    string
	Selected bool //是否选中项
}

//上传文件
type UploadFile struct {
	Name string //文件名
	Ext  string //扩展名（文件类型）
	Path string //路径
	Size int64  //文件大小
}

var db *xorm.Engine

func init() {
	var err error
	//db, err = xorm.NewEngine("mysql", "root:goldapple@127.0.0.1:3306/writer?charset=utf8")
	db, err = xorm.NewEngine("sqlite3", "./data/writer.db")

	if err != nil {
		beego.Trace(err)
	}

	db.SetMapper(core.SameMapper{})

	//db.ShowInfo = true
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

// ---------- 数据库 DAL 层 -------------------

// Select 语句
type Dal struct {
	Field   string
	From    string
	Where   string
	OrderBy string
	Size    int
	Offset  int
}

// 生成Select语句
func (this *Dal) Select() string {
	_sql := make([]string, 0)
	_sql = append(_sql, "select "+this.Field+" from "+this.From)

	if strings.TrimSpace(this.Where) != "" {
		_sql = append(_sql, "where "+this.Where)
	}

	if strings.TrimSpace(this.OrderBy) != "" {
		_sql = append(_sql, "order by  "+this.OrderBy)
	}

	if this.Size > 0 {
		_sql = append(_sql, fmt.Sprintf("limit %d offset %d", this.Size, this.Offset))
	}

	return strings.Join(_sql, " ")
}

// 记录数目统计
func (this *Dal) Count(params ...interface{}) int64 {
	sqlstr := make([]string, 0)
	sqlstr = append(sqlstr, "select count(*) as counts from "+this.From)

	if strings.TrimSpace(this.Where) != "" {
		sqlstr = append(sqlstr, "where "+this.Where)
	}

	return utils.Str2int64(this.Single(strings.Join(sqlstr, " "), "counts", params...))
}

// 返回第一行第一个字段的字符串形式
func (this *Dal) Single(sql string, field string, params ...interface{}) string {
	if rows, err := db.Query(sql, params...); len(rows) > 0 && err == nil {
		return string(rows[0][field])
	}
	return ""
}
