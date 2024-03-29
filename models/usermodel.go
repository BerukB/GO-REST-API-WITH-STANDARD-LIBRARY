package usermodel

type User struct {
	ID       string `json:"id"`
	UserName string `json:"username"`
	PassWord string `json:"PassWord"`
}
