package models

import (
	// "errors"
	"fmt"
	"github.com/astaxie/beego/validation"
	"strings"
	"zouzhe/utils"
)

// article视图
type Article struct {
	Id         int64  `json:"articleId"`
	DocumentId int64  `json:"documentId"`
	Title      string `json:"title" valid:"MaxSize(250)"`
	Content    string `json:"content"`
	Tags       string `json:"tags" valid:"MaxSize(250)"`
	MoreId     int64  `json:"moreId"`
	ParentId   int64  `json:"parentId"`
	Position   int    `json:"position"`
	Status     int    `json:"status" valid:"Range(0,1)"`
	Deleted    int    `json:"deleted" valid:"Range(0,1)"`
	Creator    int64  `json:"creator"`
	Created    int64  `json:"created"`
	Updator    int64  `json:"updator"`
	Updated    int64  `json:"updated"`
	Ip         string `json:"ip" valid:"MaxSize(23)"`
}

// Articles表
type Articles struct {
	Id         int64  `json:"articleId"`
	DocumentId int64  `json:"documentId"`
	Status     int    `json:"status" valid:"Range(0,1)"`
	Deleted    int    `json:"deleted" valid:"Range(0,1)"`
	Creator    int64  `json:"creator"`
	Created    int64  `json:"created"`
	Updator    int64  `json:"updator"`
	Updated    int64  `json:"updated"`
	Ip         string `json:"ip" valid:"MaxSize(23)"`
}

// ArticleMore表 -- article多对多映射
type ArticleMore struct {
	Id        int64  `json:"moreId"`
	ArticleId int64  `json:"articleId"`
	ParentId  int64  `json:"parentId"`
	Position  int    `json:"position"`
	Depth     string `json:"depth"`
	Updator   int64  `json:"updator"`
	Updated   int64  `json:"updated"`
	Ip        string `json:"ip" valid:"MaxSize(23)"`
}

// Documents表
type Documents struct {
	Id      int64  `json:"documentId"`
	Title   string `json:"title" valid:"MaxSize(250)"`
	Content string `json:"content"`
	Tags    string `json:"tags" valid:"MaxSize(250)"`
	Status  int    `json:"status" valid:"Range(0,1)"`
	Deleted int    `json:"deleted" valid:"Range(0,1)"`
	Creator int64  `json:"creator"`
	Created int64  `json:"created"`
	Updator int64  `json:"updator"`
	Updated int64  `json:"updated"`
	Ip      string `json:"ip" valid:"MaxSize(23)"`
}

// 文章是否存在
func (this *Article) Exists() (bool, error) {
	return db.Get(this)
}

// 自定义数据验证
func (this *Article) Valid(v *validation.Validation) {

}

// 新文章
func (this *Article) Update() (error, []Error) {
	//数据有效性检验
	if d, err := dataCheck(this); err != nil {
		return err, d
	}

	session := db.NewSession()
	defer session.Close()
	// 事务开始
	err := session.Begin()

	if err != nil {
		session.Rollback()
		return err, nil
	}

	// articles对象
	_article := new(Articles)
	_article.Id = this.Id
	_article.DocumentId = this.DocumentId
	_article.Creator = this.Creator
	_article.Created = this.Created
	_article.Updator = this.Updator
	_article.Updated = this.Updated
	_article.Ip = this.Ip

	// articlemore对象
	_more := new(ArticleMore)
	_more.ParentId = this.ParentId
	_more.Updator = this.Updator
	_more.Updated = this.Updated
	_more.Ip = this.Ip

	// documents对象
	_document := new(Documents)
	_document.Id = this.DocumentId
	_document.Title = this.Title
	_document.Content = this.Content
	_document.Tags = this.Tags
	_document.Creator = this.Creator
	_document.Created = this.Created
	_document.Updator = this.Updator
	_document.Updated = this.Updated
	_document.Ip = this.Ip

	// 如果是新增，按index腾出位置
	if this.Id == 0 {
		// 找到id=this.Position参考文档的position
		var positionSql string
		if this.Position > 0 {
			if _results, err := session.Query("select position,depth,articleid from articlemore where id=?", this.Position); len(_results) > 0 && err == nil {
				_more.Position = 0 //utils.Bytes2int(_results[0]["position"])
				_more.ParentId = utils.Bytes2int64(_results[0]["articleid"])
				_more.Depth = fmt.Sprintf("%s%d,", string(_results[0]["depth"]), _more.ParentId)
			}
			positionSql = "update articlemore set position = position+2 , updated = ? , ip = ? where parentId = ? and position >= ?"
		} else {
			if _results, err := session.Query("select position,depth,parentid from articlemore where id=?", this.MoreId); len(_results) > 0 && err == nil {
				_more.Position = utils.Bytes2int(_results[0]["position"])
				_more.ParentId = utils.Bytes2int64(_results[0]["parentid"])
				_more.Depth = string(_results[0]["depth"])
			}
			positionSql = "update articlemore set position = position+2 , updated = ? , ip = ? where parentId = ? and position > ?"
		}
		// 更新其后文档的position
		if _, err = session.Exec(positionSql, this.Updated, this.Ip, _more.ParentId, _more.Position); err != nil {
			session.Rollback()
			return err, nil
		}
		// 如果关联文档不存在，创建关联文档
		if this.DocumentId <= 0 {
			// 先 insert documents 附表
			if _, err = session.Insert(_document); err != nil {
				session.Rollback()
				return err, nil
			}
		}

		// insert articles 主表
		_article.DocumentId = _document.Id //主附表映射
		// 为返回对象赋值
		this.DocumentId = _article.DocumentId
		this.Position = _more.Position
		// 创建新文档
		if _, err = session.Insert(_article); err == nil {
			// 为返回对象赋值
			this.Id = _article.Id
			// 层次深度 depth = 父条目depth + 本条目Id
			if _article.Id > 0 {
				_more.Position += 1
				_more.ArticleId = _article.Id
				//_more.Depth = fmt.Sprintf("%s%d,", _more.Depth, _more.ParentId)

				// Update articlemore 多对多表
				if _, err = session.Insert(_more); err == nil {
					this.MoreId = _more.Id
				} else {
					session.Rollback()
					return err, nil
				}
			}
		} else {
			session.Rollback()
			return err, nil
		}
	} else {
		// Update articles 主表
		if _, err = session.Id(_article.Id).Cols("updator", "updated", "ip").Update(_article); err != nil {
			session.Rollback()
			return err, nil
		}
		// 读取关联的文档id
		if ok, err := this.GetEx(); !ok {
			session.Rollback()
			return err, nil
		}
		// Update documents 附表
		_document.Id = this.DocumentId
		if _, err = session.Id(_document.Id).Cols("title", "content", "tags", "updator", "updated", "ip").Update(_document); err != nil {
			session.Rollback()
			return err, nil
		}
	}

	// 检查tags是否为空
	if len(this.Tags) != 0 {
		_tags := strings.Split(this.Tags, ",")
		// 从数据库中读取所有标签的id,name
		tags := new(Tags)
		rows, err := session.In("name", _tags).Rows(tags)

		if err != nil {
			session.Rollback()
			return err, nil
		}
		defer rows.Close()

		ids := make([]int64, 0)
		for rows.Next() {
			err = rows.Scan(tags)
			//...
			ids = append(ids, tags.Id)
			// 移除已经存在标签
			if utils.StringsContains(_tags, tags.Name) {
				_tags = utils.RemoveStringSlice(tags.Name, _tags)
			}
		}
		// 遍历标签，找到数据库中不存在的标签，insert into标签表
		for _, v := range _tags {
			_tag := new(Tags)
			_tag.Name = v

			_, err = session.Insert(_tag)
			if err == nil {
				ids = append(ids, _tag.Id)
			} else {
				session.Rollback()
				return err, nil
			}
		}
		// 清除旧的标签-文章的索引
		_del := new(TagDocument)
		_, err = session.Where("documentId = ?", this.Id).Delete(_del)
		if err != nil {
			session.Rollback()
			return err, nil
		}
		// 建立新的标签-文章的索引
		tagDocuments := make([]TagDocument, 0)
		for _, id := range ids {
			tagDocuments = append(tagDocuments, TagDocument{TagId: id, DocumentId: this.Id})
		}
		fmt.Println(tagDocuments)
		_, err = session.Insert(tagDocuments)
		if err != nil {
			session.Rollback()
			return err, nil
		}
	}
	// 提交事务
	err = session.Commit()

	return err, nil
}

// 只读取可见的
func (this *Article) Get() (bool, error) {
	return this._get(true)
}

// 读取全部
func (this *Article) GetEx() (bool, error) {
	return this._get(false)
}

// 读取
func (this *Article) _get(all bool) (bool, error) {
	// Dal对象
	_dal := &Dal{}
	_dal.Field = "articles.*,documents.title,documents.content"
	_dal.From = "articles,documents"
	_dal.Where = "documents.id = articles.documentId and articles.id=?"

	// 可见的
	if all {
		_dal.Where += fmt.Sprintf(" and articles.status=%d and articles.deleted=%d and documents.status=%d and documents.deleted=%d", Unlock, Undelete, Unlock, Undelete)
	}

	return db.Sql(_dal.Select(), this.Id).Get(this)
}

// 分页列表可见的
func (this *Article) List(page *Pagination, condition string, params ...interface{}) ([]Article, error) {
	return this._list(true, page, condition, params...)
}

// 分页列表全部的
func (this *Article) ListEx(page *Pagination, condition string, params ...interface{}) ([]Article, error) {
	return this._list(false, page, condition, params...)
}

// 分页列表
func (this *Article) _list(view bool, page *Pagination, condition string, params ...interface{}) ([]Article, error) {
	// Dal对象
	_dal := &Dal{}
	_dal.From = "articlemore,articles,documents"
	_dal.Where = "articlemore.articleId = articles.id and documents.id = articles.documentId"
	_dal.OrderBy = "articlemore.parentId,articlemore.position"

	// 可见的
	if view {
		_dal.Where += fmt.Sprintf(" and articles.status=%d and articles.deleted=%d and documents.status=%d and documents.deleted=%d", Unlock, Undelete, Unlock, Undelete)
	}
	// 条件
	if strings.TrimSpace(condition) != "" {
		_dal.Where += " and " + condition
	}
	// slice承载返回的结果
	as := make([]Article, 0)

	// 读取符合条件的记录总数
	if page.Count == 0 {
		// 读取符合条件的记录总数
		if rows := _dal.Count(params...); rows > 0 {
			// 读取总页数
			getPageCount(rows, page)
		}
	}

	if page.Count > 0 {

		_dal.Size = page.Size
		_dal.Offset = page.Index * page.Size

		_dal.Field = "articlemore.id as moreid,articlemore.parentid,articlemore.position,articlemore.updator,articles.id,articles.documentid,articles.creator,documents.title,documents.content"
		err := db.Sql(_dal.Select(), params...).Find(&as)
		return as, err
	}
	return as, nil
}

// 设置locked、deleted等状态字段的公共方法
func (this *Article) SetStatus(action string) error {
	a := new(Articles)
	a.Id = this.Id

	// 验证当前用户的权限
	if has, err := db.Id(a.Id).Cols("creator").Get(a); has {
		// 如果不是作者，只能删除映射关系
		if a.Creator != this.Updator {
			_, err = db.Exec("delete from articlemore where articleid = ? and parentid = ?", this.Id, this.ParentId)
			return err
		}
	} else {
		return err
	}

	a.Updated = this.Updated
	a.Updator = this.Updator
	a.Ip = this.Ip

	d := new(Documents)
	d.Updated = this.Updated
	d.Updator = this.Updator
	d.Ip = this.Ip

	switch strings.ToLower(action) {
	case "lock":
		action = "status"
		a.Status = Locked
		d.Status = Locked
	case "unlock":
		action = "status"
		a.Status = Unlock
		d.Status = Unlock
	case "delete":
		action = "deleted"
		a.Deleted = Deleted
		d.Deleted = Deleted
	case "undelete":
		action = "deleted"
		a.Deleted = Undelete
		d.Deleted = Undelete
	}

	// 事务
	session := db.NewSession()
	defer session.Close()
	// 事务开始
	err := session.Begin()

	if err != nil {
		session.Rollback()
		return err
	}
	// 修改articles状态
	_, err = session.Id(a.Id).Cols(action, "updator", "updated", "ip").Update(a)

	if err != nil {
		session.Rollback()
		return err
	}

	// 读取关联的文档id
	if ok, err := this.GetEx(); !ok {
		session.Rollback()
		return err
	}
	d.Id = this.DocumentId

	// 修改documents的状态
	_, err = session.Id(d.Id).Where("creator=?", this.Creator).Cols(action, "updator", "updated", "ip").Update(d)
	if err != nil {
		session.Rollback()
		return err
	}
	// 完成事务提交
	session.Commit()

	return err
}

// 设置文档节点位置
func (this *Article) SetPosition() (bool, error, *ArticleMore) {
	// articlemore对象
	_more := new(ArticleMore)
	_more.Id = this.Id
	_more.Updator = this.Updator
	_more.Updated = this.Updated
	_more.Ip = this.Ip

	session := db.NewSession()
	defer session.Close()
	// 事务开始
	err := session.Begin()

	if err != nil {
		session.Rollback()
		return false, err, _more
	}

	// 参考点的position
	var positionSql string
	if this.Position > 0 {
		if _results, err := session.Query("select position,parentid,depth from articlemore where id=?", this.Position); len(_results) > 0 && err == nil {
			_more.Position = utils.Bytes2int(_results[0]["position"])
			_more.ParentId = utils.Bytes2int64(_results[0]["parentid"])
			_more.Depth = string(_results[0]["depth"])
		} else {
			session.Rollback()
			return false, err, _more
		}
		positionSql = "update articlemore set position = position+2 , updated = ? , ip = ? where parentId = ? and position > ?"
	} else if this.MoreId > 0 {
		if _results, err := session.Query("select position,articleid,depth from articlemore where id=?", this.MoreId); len(_results) > 0 && err == nil {
			_more.ParentId = utils.Bytes2int64(_results[0]["articleid"])
			_more.Depth = fmt.Sprintf("%s%d,", string(_results[0]["depth"]), _more.ParentId)
		} else {
			session.Rollback()
			return false, err, _more
		}
		positionSql = "update articlemore set position = position+2 , updated = ? , ip = ? where parentId = ? and position >= ?"
	} else {
		_more.ParentId = 0
		positionSql = "update articlemore set position = position+2 , updated = ? , ip = ? where parentId = ? and position >= ?"
	}
	// 更新其后文档的position
	if _, err = session.Exec(positionSql, this.Updated, this.Ip, _more.ParentId, _more.Position); err != nil {
		session.Rollback()
		return false, err, _more
	}

	_more.Position += 1

	// 源文档的作者
	var _updator int64

	// 更新所有子节点的depth
	if _results, err := session.Query("select parentid,articleid,depth,updator from articlemore where id=?", _more.Id); len(_results) > 0 && err == nil {
		_updator = utils.Bytes2int64(_results[0]["updator"])
		// 父节点已改变
		if _more.ParentId != utils.Bytes2int64(_results[0]["parentid"]) {
			// 更新本条目所有子条目的 depth
			_more.ArticleId = utils.Bytes2int64(_results[0]["articleid"])
			_old_Depth := fmt.Sprintf("%s%d,", string(_results[0]["depth"]), _more.ArticleId)
			_new_Depth := fmt.Sprintf("%s%d,", _more.Depth, _more.ArticleId)

			if _, err = session.Exec(fmt.Sprintf("update articlemore set depth = concat('%s',substring(depth,%d)) where depth like '%s%s'", _new_Depth, len(_old_Depth)+1, _old_Depth, "%")); err != nil {
				session.Rollback()
				return false, err, _more
			}
		}
	} else {
		session.Rollback()
		return false, err, _more
	}

	// 如果修改自己的文档
	if _updator == _more.Updator {
		_, err = session.Id(_more.Id).Cols("parentId", "depth", "position", "updator", "updated", "ip").Update(_more)
	} else {
		_more.Id = 0
		_, err = session.Insert(_more)
	}
	if err != nil {
		session.Rollback()
		return false, err, _more
	}
	// 提交事务
	err = session.Commit()
	return err == nil, err, _more
}

// 读取目录
func (this *Article) Catalog(id int64) ([]Article, error) {
	// Dal对象
	_dal := &Dal{}
	_dal.Field = "articlemore.id as moreid,articlemore.parentid,articlemore.position,articlemore.updator,articles.id,articles.documentid,articles.creator,documents.title"
	_dal.From = "articlemore,articles,documents"
	_dal.Where = fmt.Sprintf("documents.id = articles.documentId and articles.id = articlemore.articleid and articlemore.depth like '%d,%s'", id, "%")
	_dal.OrderBy = "articlemore.depth,articlemore.position"

	// 可见的
	_dal.Where += fmt.Sprintf(" and articles.status=%d and articles.deleted=%d and documents.status=%d and documents.deleted=%d", Unlock, Undelete, Unlock, Undelete)

	// slice承载返回的结果
	as := make([]Article, 0)

	err := db.Sql(_dal.Select()).Find(&as)
	return as, err
}

// 内容分页列表（可见子孙节点）
func (this *Article) GetContent(page *Pagination, id int64, condition string, params ...interface{}) ([]Article, error) {
	return this._content(true, page, id, condition, params...)
}

// 内容分页列表（全部子孙节点）
func (this *Article) GetContentEx(page *Pagination, id int64, condition string, params ...interface{}) ([]Article, error) {
	return this._content(false, page, id, condition, params...)
}

// 内容分页列表（子孙节点）
func (this *Article) _content(view bool, page *Pagination, id int64, condition string, params ...interface{}) ([]Article, error) {
	// Dal对象
	_dal := &Dal{}
	_dal.From = "articlemore,articles,documents"
	_dal.Where = fmt.Sprintf("documents.id = articles.documentId and articles.id = articlemore.articleid and articlemore.depth like '%d,%s'", id, "%")

	// 可见的
	if view {
		_dal.Where += fmt.Sprintf(" and articles.status=%d and articles.deleted=%d and documents.status=%d and documents.deleted=%d", Unlock, Undelete, Unlock, Undelete)
	}
	// 条件
	if strings.TrimSpace(condition) != "" {
		_dal.Where += " and " + condition
	}
	// slice承载返回的结果
	as := make([]Article, 0)
	// 读取符合条件的记录总数
	if rows := _dal.Count(params...); rows > 0 {

		getPageCount(rows, page)

		_dal.Size = page.Size
		_dal.Offset = page.Index * page.Size

		_dal.Field = "articlemore.id as moreid,articlemore.parentid,articlemore.position,articlemore.updator,articles.id,articles.documentid,articles.creator,documents.title,documents.content"
		_dal.OrderBy = "articlemore.depth,articlemore.position"

		err := db.Sql(_dal.Select(), params...).Find(&as)
		return as, err
	}
	return as, nil
}

// 读取一个文档
func (this *Article) GetSingle(condition string, params ...interface{}) error {
	return this._single(true, condition, params...)
}

// 读取一个文档
func (this *Article) _single(view bool, condition string, params ...interface{}) error {
	// Dal对象
	_dal := &Dal{}
	_dal.Field = "articles.id,articles.documentid,articles.creator,documents.title,documents.content"
	_dal.From = "articles,documents"
	_dal.Where = "documents.id = articles.documentId"

	// 可见的
	if view {
		_dal.Where += fmt.Sprintf(" and articles.status=%d and articles.deleted=%d and documents.status=%d and documents.deleted=%d", Unlock, Undelete, Unlock, Undelete)
	}
	// 条件
	if strings.TrimSpace(condition) != "" {
		_dal.Where += " and " + condition
	}

	_, err := db.Sql(_dal.Select(), params...).Get(this)
	return err
}
