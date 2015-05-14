package models

import (
	// "errors"
	"fmt"
	"github.com/astaxie/beego/validation"
)

// 账户
type Accounts struct {
	Id           int64  `json:"accoundId"`
	LoginName    string `json:"loginName" valid:"MaxSize(50)"`
	Password     string `json:"password"`
	RealName     string `json:"realName" valid:"MaxSize(50)"`
	OpenId       string `json:"openId" valid:"MaxSize(32)"`
	OpenFrom     string `json:"openFrom" valid:"MaxSize(10)"`
	NickName     string `json:"nickName" valid:"MaxSize(50)"`
	Gender       int    `json:"gender" valid:"Range(0,1)"`
	Avatar_1     string `json:"avatar_1" valid:"MaxSize(250)"`
	Avatar_2     string `json:"avatar_2" valid:"MaxSize(250)"`
	AccessToken  string `json:"accessToken" valid:"MaxSize(150)"`
	RefreshToken string `json:"refreshToken" valid:"MaxSize(150)"`
	Status       int    `json:"status" valid:"Range(0,1)"`
	Deleted      int    `json:"deleted" valid:"Range(0,1)"`
	Updator      int64  `json:"updator"`
	Updated      int64  `json:"updated"`
	Ip           string `json:"ip" valid:"MaxSize(23)"`
}

// 账号是否存在
func (this *Accounts) Exists() (bool, error) {
	fmt.Println(this)
	return db.Get(this)
}

// 自定义数据验证
func (this *Accounts) Valid(v *validation.Validation) {
	//登录名必须是email
	if this.LoginName != "" {
		v.Email(this.LoginName, "loginName")
	}
}

// 增加新账户
func (this *Accounts) Post() (error, []Error) {

	if this.RealName == "" {
		this.RealName = this.NickName
	}
	//数据有效性检验
	if d, err := dataCheck(this); err != nil {
		return err, d
	}
	//
	_, err := db.Insert(this)
	return err, nil
}

// 更新第三方登录的refreshToken
func (this *Accounts) RefreshAccessToken() (int64, error) {
	return db.Id(this.Id).Cols("accessToken", "refreshToken", "updated").Update(this)
}

// --------------------------------------------------
