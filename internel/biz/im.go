package biz

import (
	"context"
	"fileStore/internel/data"
	"fileStore/internel/domain"
	"time"
)

func SaveASessionContent(ctx context.Context, content domain.ImSessionContent) error {
	return data.CreateASessionMessage(ctx, content)
}

func UpdateSessionUpdateTime(ctx context.Context, sessionUuid string) error {
	return data.UpdateUserSessionUpdatedAtBySessionUuid(ctx, sessionUuid, time.Now())
}
