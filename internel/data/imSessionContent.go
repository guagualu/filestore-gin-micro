package data

import (
	"context"
	"fileStore/internel/domain"
	"fileStore/log"
	"time"
)

type imSessionContent struct {
	Id             uint      `gorm:"primarykey autoIncrement"`                                                      // 物理主键
	SessionUuid    string    `gorm:"column:session_uuid;type:char(36);not null;default:''" json:"session_uuid"`     // 业务主键
	SendUserUuid   string    `gorm:"column:send_user_uuid;type:char(36);not null;default:''" json:"send_user_uuid"` // 发送者uuid
	ToUserUuid     string    `gorm:"column:to_user_uuid;type:char(36);not null;default:''" json:"to_user_uuid"`     // 接受者uuid
	MessageType    int       `gorm:"column:message_type;type:int(3);not null;default:0" json:"message_type"`        //0为普通文本消息，1为文件消息 content存的是文件hash
	MessageContent string    `gorm:"column:message_content;type:text;not null;" json:"message_content"`
	CreatedAt      time.Time `gorm:"column:created_at;type:datetime comment '创建时间';not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at;type:datetime comment '更新时间';not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
}

func GetSessionAllMessage(ctx context.Context, SessionUuid string) ([]domain.ImSessionContent, error) {
	sessionRes := make([]imSessionContent, 0)
	if err := GetData().DB(ctx).Debug().Where("session_uuid = ?", SessionUuid).Find(&sessionRes).Error; err != nil {
		log.Logger.Error("GetSessionAllMessage err:", err)
		return nil, err
	}
	res := make([]domain.ImSessionContent, 0)
	for _, v := range sessionRes {
		res = append(res, domain.ImSessionContent{
			Id:             v.Id,
			SessionUuid:    v.SessionUuid,
			SendUserUuid:   v.SendUserUuid,
			ToUserUuid:     v.ToUserUuid,
			MessageType:    v.MessageType,
			MessageContent: v.MessageContent,
			CreatedAt:      v.CreatedAt,
			UpdatedAt:      v.UpdatedAt,
		})
	}
	return res, nil
}

func CreateASessionMessage(ctx context.Context, message domain.ImSessionContent) (int, error) {
	msg := imSessionContent{
		SessionUuid:    message.SessionUuid,
		SendUserUuid:   message.SendUserUuid,
		ToUserUuid:     message.ToUserUuid,
		MessageType:    message.MessageType,
		MessageContent: message.MessageContent,
	}
	if err := GetData().DB(ctx).Create(&msg).Error; err != nil {
		log.Logger.Error("SaveASessionMessage err:", err)
		return 0, err
	}
	return int(msg.Id), nil
}
