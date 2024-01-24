package biz

import (
	"context"
	"errors"
	"fileStore/internel/data"
	"fileStore/internel/domain"
	"gorm.io/gorm"
	"time"
)

func SaveASessionContent(ctx context.Context, content domain.ImSessionContent) (int, error) {
	return data.CreateASessionMessage(ctx, content)
}

func UpdateSessionUpdateTime(ctx context.Context, sessionUuid string) error {
	return data.UpdateUserSessionUpdatedAtBySessionUuid(ctx, sessionUuid, time.Now())
}

func CreateSession(ctx context.Context, userAUuid, userBUuid string) (*domain.ImSession, error) {
	return data.CreateSession(ctx, userAUuid, userBUuid)
}

func GetUserAllSession(ctx context.Context, userUuid string) ([]domain.ImSession, error) {
	list, err := data.ListUserSessionsByUserUuid(ctx, userUuid)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func GetUserASession(ctx context.Context, sessionUuid string) ([]domain.ImSessionContent, error) {
	list, err := data.GetSessionAllMessage(ctx, sessionUuid)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func IsSeesionExist(ctx context.Context, userAUuid, userBUuid string) (*domain.ImSession, bool, error) {
	session, err := data.GetSessionByUsers(ctx, userAUuid, userBUuid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return session, true, nil
}
