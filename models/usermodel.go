package usermodel

type User struct {
	ID       string `json:"id"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	PassWord string `json:"PassWord"`
}
