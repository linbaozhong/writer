package controllers

import "writer/models"

type Account struct {
	Authorized
}

func (a *Account) Login() {
	account := &models.Account{}
	account.Login()
	a.Ctx.Output.Json(account, true, true)

	//a.toString(fmt.Sprintf("O(∩_∩)O哈哈~ %d", models.Deleted))
}
