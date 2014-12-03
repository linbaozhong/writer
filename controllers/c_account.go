package controllers

import (
	"fmt"
	"writer/models"
)

type Account struct {
	Authorized
}

func (a *Account) Login() {
	account := &models.Account{}

	fmt.Println(account)

	a.toString(fmt.Sprintf("O(∩_∩)O哈哈~ %d", models.Deleted))
}
