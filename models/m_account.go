package models

import (
	"fmt"
)

type Account struct {
	Id       int64  `json:"id"`
	OpenId   string `json:"openId"`
	OpenFrom int    `json:"openFrom"`
	OpenName string `json:"openName"`
	RegTime  int64  `json:"regTime"`
	Status   int    `json:"status"`
}

func (a *Account) Login() {
	//defer db.Close()

	has, err := db.Id(0).Get(a)
	fmt.Println(a)
	fmt.Println(&a)
	if err != nil {
		trace(err)
	}
	if has {
		trace(a)
		trace(&a)
	}
}
