package domain

import "time"

var ServiceName string

type User struct {
	Id       uint
	Uuid     string
	NickName string
	Email    string
	Password string
	Mobile   string
}

type Friends struct {
	Id          uint
	Uuid        string    `json:"uuid"` // 业务主键
	UserAMobile string    `json:"user_a_mobile"`
	UserBMobile string    `json:"user_b_mobile"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
