package data

import (
	"context"
	"encoding/json"
	"errors"
	"fileStore/internel/domain"
	"fileStore/internel/pkg/uuid"
	"fileStore/log"
	"github.com/gomodule/redigo/redis"
	"gorm.io/gorm"
	"time"
)

type User struct {
	Id        uint      `gorm:"primarykey autoIncrement"`                                         // 物理主键
	Uuid      string    `gorm:"column:uuid;unique;type:char(36);not null;default:''" json:"uuid"` // 业务主键
	NickName  string    `gorm:"column:nick_name;type:varchar(255) comment '用户名';not null;default:''" json:"nick_name"`
	Email     string    `gorm:"column:email;type:varchar(50) comment '邮箱';not null;default:''" json:"email"`
	Password  string    `gorm:"column:password;type:varchar(255) comment '密码';not null;default:''" json:"password"`
	Mobile    string    `gorm:"column:mobile;unique;type:varchar(20) comment '手机号';not null default:''" json:"mobile"`
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
	u.Uuid = uuid.NewUuid()
	if err := db.DB(ctx).Omit("created_at", "updated_at", "next_expire_time").Create(&u).Error; err != nil {
		log.Logger.Error("数据库错误:", err)
		return err
	}
	return nil
}

func GetUserByPhoneAndPsd(ctx context.Context, user domain.User) (*domain.User, error) {
	db := GetData()
	var u User
	if err := db.DB(ctx).Where("mobile = ? and password = ?", user.Mobile, user.Password).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		log.Logger.Error("数据库错误:", err)
		return nil, err
	}
	res := new(domain.User)
	res.Email = u.Email
	res.Mobile = u.Email
	res.Uuid = u.Uuid
	res.Id = u.Id
	res.NickName = u.NickName
	res.Password = u.NickName
	return res, nil
}

func GetUserInfoByCache(ctx context.Context, userUuid string) (*domain.User, error) {
	user := new(domain.User)
	conn := GetData().red.Get()
	defer conn.Close()
	res, err := redis.String(conn.Do("GET", "user_"+userUuid))
	if err != nil {
		return nil, err
	} else {
		err := json.Unmarshal([]byte(res), user)
		if err != nil {
			log.Logger.Error("json反序列化 user 失败")
			return nil, err
		}
		return user, nil
	}
}
func SetUserInfoByCache(ctx context.Context, userInfo domain.User) error {
	conn := GetData().red.Get()
	defer conn.Close()
	info, _ := json.Marshal(userInfo)
	_, err := conn.Do("SET", "user_"+userInfo.Uuid, info)
	if err != nil {
		return err
	}
	return nil
}

func GetUserInfo(ctx context.Context, userUuid string) (*domain.User, error) {
	var u User
	if err := GetData().DB(ctx).Where("uuid = ?", userUuid).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &domain.User{
		Id:       u.Id,
		Uuid:     u.Uuid,
		NickName: u.NickName,
		Email:    u.Email,
		Password: u.Password,
		Mobile:   u.Mobile,
	}, nil
}

func ListUserInfoByMobile(ctx context.Context, mobiles []string) ([]domain.User, error) {
	var u []User
	if err := GetData().DB(ctx).Where("mobile in ?", mobiles).Find(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	res := make([]domain.User, 0)
	for _, v := range u {
		res = append(res, domain.User{
			Id:       v.Id,
			Uuid:     v.Uuid,
			NickName: v.NickName,
			Email:    v.Email,
			Password: v.Password,
			Mobile:   v.Mobile,
		})
	}
	return res, nil
}
