package models

type Account struct {
	Id     int64
	OpenId string
}

func (a *Account) Login() {
	defer db.Close()

	has, err := db.Id(0).Get(a)
	if err != nil {
		trace(err)
	}
	if has {

	}
}
