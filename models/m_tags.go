package models

type Tags struct {
	Id    int64  `json:"tagsId"`
	Name  string `json:"name" valid:"MaxSize(100)"`
	Times int64  `json:"times"`
}

//type TagDocument struct {
//	Id         int64
//	TagId      int64
//	DocumentId int64
//}
type TagArticle struct {
	TagId     int64
	ArticleId int64
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

func (this *Tags) GetArticleIdsByTags(tags []string) ([]int64, error) {
	ids = make([]int64, 0)

	rows, err := db.In("name", tags).Cols("id").Rows(this)
	if err != nil {
		return ids, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(this)
		if err != nil {
			continue
		}
		ids = append(ids, this.Id)
	}
	return ids, nil
}
