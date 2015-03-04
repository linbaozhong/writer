package models

type Tags struct {
	Id    int64  `json:"tagsId"`
	Name  string `json:"name" valid:"MaxSize(100)"`
	Times int64  `json:"times"`
}
type TagArticle struct {
	TagId      int64
	DocumentId int64
}

// 标签是否存在
func (this *Tags) Exists() (bool, error) {
	return db.Get(this)
}

//
func (this *Tags) List() ([]Tags, error) {
	// slice承载返回的结果
	ts := make([]Tags, 0)

	err := db.Desc("times").Find(&ts)
	return ts, err
}
