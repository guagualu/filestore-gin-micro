package data

import (
	"context"
	"fileStore/internel/domain"
	"fileStore/internel/pkg/uuid"
	"fileStore/log"
	"time"
)

type imSession struct {
	Id          uint      `gorm:"primarykey autoIncrement"`                                                         // 物理主键
	SessionUuid string    `gorm:"column:session_uuid;unique;type:char(36);not null;default:''" json:"session_uuid"` // 业务主键
	UserAUuid   string    `gorm:"column:user_a_uuid;type:char(36);not null;default:''" json:"user_a_uuid"`          // 发送者uuid
	UserBUuid   string    `gorm:"column:user_b_uuid;type:char(36);not null;default:''" json:"user_b_uuid"`          // 接受者uuid
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime comment '创建时间';not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime comment '更新时间';not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
}

func CreateSession(ctx context.Context, userAUuid, userBUuid string) error {
	line := imSession{
		SessionUuid: uuid.NewUuid(),
		UserAUuid:   userAUuid,
		UserBUuid:   userBUuid,
	}
	if err := GetData().DB(ctx).Create(&line).Error; err != nil {
		log.Logger.Error("CreateSession err:", err)
		return err
	}
	return nil
}

func ListUserSessionsByUserUuid(ctx context.Context, userUuid string) ([]domain.ImSession, error) {
	list := make([]imSession, 0)
	if err := GetData().DB(ctx).Where("user_a_uuid = ? or user_b_uuid= ?", userUuid, userUuid).Order("updated_at desc").Find(&list).Error; err != nil {
		log.Logger.Error("ListUserSessionsByUserUuid err:", err)
		return nil, err
	}
	res := make([]domain.ImSession, 0)
	for _, v := range list {
		res = append(res, domain.ImSession{
			Id:          v.Id,
			SessionUuid: v.SessionUuid,
			UserAUuid:   v.UserAUuid,
			UserBUuid:   v.UserBUuid,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		})
	}
	return res, nil
}

func UpdateUserSessionUpdatedAtBySessionUuid(ctx context.Context, sessionUuid string, updatedTime time.Time) error {
	if err := GetData().DB(ctx).Table("im_session").Where("session_uuid = ? ", sessionUuid).Update("updated_at", updatedTime).Error; err != nil {
		log.Logger.Error("UpdateUserSessionUpdatedAtBySessionUuid err:", err)
		return err
	}
	return nil
}
