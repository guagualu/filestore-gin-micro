package domain

var ServiceName string

type User struct {
	Id       uint
	Uuid     string
	NickName string
	Email    string
	Password string
	Mobile   string
}
