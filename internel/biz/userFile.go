package biz

import (
	"context"
	"fileStore/internel/data"
	"fileStore/internel/domain"
)

func GetUserFileList(ctx context.Context, userUuid string, page, pageSize int) ([]domain.UserFile, int64, error) {
	return data.ListUserFiles(ctx, userUuid, page, pageSize)
}

func DeletedUserFileList(ctx context.Context, userUuid string, fileHashs []string) error {
	return data.DeleteUserFiles(ctx, fileHashs, userUuid)
}

func RenameUserFile(ctx context.Context, userUuid string, fileHash string, fileName string) error {
	return data.RenameUserFile(ctx, userUuid, fileHash, fileName)
}
