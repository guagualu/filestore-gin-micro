package data

import (
	"context"
	"fileStore/internel/domain"
	"fileStore/internel/pkg/code/errcode"
	"time"
)

type User struct {
	Id        uint      `gorm:"primarykey autoIncrement"`                                         // 物理主键
	Uuid      string    `gorm:"column:uuid;unique;type:char(36);not null;default:''" json:"uuid"` // 业务主键
	NickName  string    `gorm:"column:nick_name;type:varchar(255) comment '用户名';not null;default:''" json:"nick_name"`
	Email     string    `gorm:"column:email;type:varchar(50) comment '邮箱';not null;default:''" json:"email"`
	Password  string    `gorm:"column:password;type:varchar(255) comment '密码';not null;default:''" json:"password"`
	Mobile    string    `gorm:"column:mobile;type:varchar(20) comment '手机号';not null;default:''" json:"mobile"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime comment '创建时间';not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime comment '更新时间';not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
}

func CreatUser(ctx context.Context, user domain.User) error {
	db := GetData()
	u := User{
		NickName: user.NickName,
		Email:    user.Email,
		Password: user.Password,
		Mobile:   user.Mobile,
	}
	if err := db.DB(ctx).Omit("created_at", "updated_at", "next_expire_time").Create(&u).Error; err != nil {
		return errcode.WithCode(errcode.Database_err, "数据库错误")
	}
	return nil
}
