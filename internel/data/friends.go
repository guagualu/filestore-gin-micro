package data

import (
	"context"
	"fileStore/internel/domain"
	"fileStore/internel/pkg/uuid"
	"fileStore/log"
	"time"
)

type Friends struct {
	Id          uint      `gorm:"primarykey autoIncrement"`                                         // 物理主键
	Uuid        string    `gorm:"column:uuid;unique;type:char(36);not null;default:''" json:"uuid"` // 业务主键
	UserAMobile string    `gorm:"column:user_a_mobile;type:varchar(20) comment '用户a手机号';not null default:''" json:"user_a_mobile"`
	UserBMobile string    `gorm:"column:user_b_mobile;type:varchar(20) comment '用户b手机号';not null default:''" json:"user_b_mobile"`
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime comment '创建时间';not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime comment '更新时间';not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
}

func CreatFriendShip(ctx context.Context, friend domain.Friends) error {
	db := GetData()
	f := Friends{
		Uuid:        uuid.NewUuid(),
		UserAMobile: friend.UserAMobile,
		UserBMobile: friend.UserBMobile,
	}
	if err := db.DB(ctx).Omit("created_at", "updated_at").Create(&f).Error; err != nil {
		log.Logger.Error("数据库错误:", err)
		return err
	}
	return nil
}

func GetUserFriendsByUserPhone(ctx context.Context, userPhone string) ([]domain.Friends, error) {
	db := GetData()
	var list []Friends
	if err := db.DB(ctx).Where("user_a_mobile= ? or user_b_mobile = ?", userPhone, userPhone).Find(&list).Error; err != nil {
		log.Logger.Error("数据库错误:", err)
		return nil, err
	}
	res := make([]domain.Friends, 0)
	for _, v := range list {
		res = append(res, domain.Friends{
			Id:          v.Id,
			Uuid:        v.Uuid,
			UserAMobile: v.UserBMobile,
			UserBMobile: v.UserBMobile,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		})
	}
	return res, nil
}
