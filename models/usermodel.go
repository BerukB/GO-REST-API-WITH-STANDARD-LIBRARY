package usermodel

type User struct {
	ID       string      `json:"id"`
	UserName string      `json:"username,omitempty"`
	Email    string      `json:"email"`
	Phone    PhoneNumber `json:"phone,omitempty"`
	PassWord string      `json:"PassWord"`
	Address  string      `json:"address"`
}

type PhoneNumber string

func (p PhoneNumber) Format() string {
	last9Digits := string(p)[len(p)-9:]
	return "251" + last9Digits
}
