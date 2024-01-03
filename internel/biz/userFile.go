package biz

import (
	"context"
	"fileStore/depart/file/client"
	filePb "fileStore/depart/file/proto"
	"fileStore/internel/data"
	"fileStore/internel/domain"
)

func GetUserFileList(ctx context.Context, userUuid string, page, pageSize int) ([]domain.UserFile, int64, error) {
	return data.ListUserFiles(ctx, userUuid, page, pageSize)
}

func DeletedUserFileList(ctx context.Context, userUuid string, fileIds []int) error {
	return data.DeleteUserFiles(ctx, fileIds, userUuid)
}

func RenameUserFile(ctx context.Context, userUuid string, fileHash string, fileName, fileOldName string) error {
	return data.RenameUserFile(ctx, userUuid, fileHash, fileName, fileOldName)
}

func GetSoftDeletedUserFileList(ctx context.Context, userUuid string, page, pageSize int) ([]domain.UserFile, int64, error) {
	return data.GetSoftDeletedUserFiles(ctx, userUuid, page, pageSize)
}

func RealDeletedUserFileList(ctx context.Context, userUuid string, fileIds []int) error {
	return data.RealDeleteUserFiles(ctx, fileIds, userUuid)
}

// 获取hash与 size 映射的map
func GetFileHashAndFileSizeMap(ctx context.Context, fileHashs []string, destMap map[string]int) error {
	files, err := client.GetFileClient().ListFile(ctx, &filePb.ListFileReq{FileHash: fileHashs})
	if err != nil {
		return err
	}
	for _, v := range files.FileList {
		destMap[v.FileHash] = int(v.FileSize)
	}
	return nil
}
