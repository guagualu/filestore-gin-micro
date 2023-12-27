package data

import (
	"context"
	"fileStore/internel/domain"
	"fileStore/internel/pkg/code/errcode"
	"fileStore/log"
	"gorm.io/gorm"
	"time"
)

type UserFile struct {
	ID       uint           `gorm:"column:id;type:uint;primary_key;autoIncrement;comment:'物理主键'" json:"id"`
	FileHash string         `gorm:"column:file_hash;type:char(40);Index:idx_file_hash;not null;comment:'文件hash值'" json:"file_hash"`
	UserUuid string         `gorm:"column:user_uuid;type:varchar(40);not null;comment:'用户uuid'"  json:"user_uuid"`
	FileName string         `gorm:"column:file_name;type:varchar(256);not null;comment:'文件名'"  json:"file_name"`
	CreateAt time.Time      `gorm:"column:created_at;type:datetime comment '创建时间';not null;default:CURRENT_TIMESTAMP"  json:"create_at"`
	UpdateAt time.Time      `gorm:"column:updated_at;type:datetime comment '更新时间';not null;default:CURRENT_TIMESTAMP;autoUpdateTime" json:"update_at"`
	Status   gorm.DeletedAt `json:"status"`
}

func SaveUserFile(ctx context.Context, userFile domain.UserFile) error {
	db := GetData()
	uf := UserFile{
		FileHash: userFile.FileHash,
		FileName: userFile.FileName,
		UserUuid: userFile.UserUuid,
	}
	if err := db.DB(ctx).Omit("created_at", "updated_at").Create(&uf).Error; err != nil {
		log.Logger.Error("SaveUserFile err:", err)
		return errcode.WithCode(errcode.Database_err, "数据库错误")
	}
	return nil
}

func GetUserFiles(ctx context.Context, userFile domain.UserFile) (*domain.UserFile, error) {
	db := GetData()
	var res UserFile
	if err := db.DB(ctx).Where("user_uuid=? and file_hash =? and file_name = ?", userFile.UserUuid, userFile.FileHash, userFile.FileName).First(&res).Error; err != nil {
		log.Logger.Error("GetUserFiles err:", err)
		return nil, errcode.WithCode(errcode.Database_err, "数据库错误")
	}
	return &domain.UserFile{
		ID:       res.ID,
		FileHash: res.FileHash,
		UserUuid: res.UserUuid,
		FileName: res.FileName,
		CreateAt: res.CreateAt,
		UpdateAt: res.UpdateAt,
	}, nil
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}

		switch {
		case pageSize > 10000:
			pageSize = 10000
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func ListUserFiles(ctx context.Context, userUuid string, page, pageSize int) ([]domain.UserFile, int64, error) {
	db := GetData()
	list := make([]*UserFile, 0)
	var sum int64
	if err := db.DB(ctx).Where("user_uuid=?", userUuid).Count(&sum).Find(&list).Scopes(Paginate(page, pageSize)).Error; err != nil {
		log.Logger.Error(errcode.WithCode(errcode.Database_err, "数据库错误"))
		return nil, 0, errcode.WithCode(errcode.Database_err, "数据库错误")
	}
	res := make([]domain.UserFile, 0)
	for _, v := range list {
		res = append(res, domain.UserFile{
			ID:       v.ID,
			FileHash: v.FileHash,
			UserUuid: v.UserUuid,
			FileName: v.FileName,
			CreateAt: v.CreateAt,
			UpdateAt: v.UpdateAt,
		})
	}
	return res, sum, nil
}

func DeleteUserFile(ctx context.Context, userFile domain.UserFile) error {
	db := GetData()
	if err := db.DB(ctx).Where("file_hash = ?and user_uuid = ?", userFile.FileHash, userFile.UserUuid).Delete(&File{}).Error; err != nil {
		log.Logger.Error("DeleteUserFile err:", err)
		return errcode.WithCode(errcode.Database_err, "数据库错误")
	}
	return nil
}

func DeleteUserFiles(ctx context.Context, fileHashs []string, userUuid string) error {
	db := GetData()
	if err := db.DB(ctx).Where("file_hash in ? and user_uuid = ?", fileHashs, userUuid).Delete(&File{}).Error; err != nil {
		log.Logger.Error("DeleteUserFile err:", err)
		return errcode.WithCode(errcode.Database_err, "数据库错误")
	}
	return nil
}

func RenameUserFile(ctx context.Context, userUuid string, fileHash string, fileName string) error {
	db := GetData()
	if err := db.DB(ctx).Where("file_hash = ? and user_uuid = ?", fileHash, userUuid).Table("user_file").Update("file_name", fileName).Error; err != nil {
		log.Logger.Error("DeleteUserFile err:", err)
		return errcode.WithCode(errcode.Database_err, "数据库错误")
	}
	return nil
}
