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
		log.Logger.Error(errcode.WithCode(errcode.Database_err, "数据库错误"))
		return errcode.WithCode(errcode.Database_err, "数据库错误")
	}
	return nil
}

func ListUserFiles(ctx context.Context, userFile domain.UserFile) ([]*UserFile, int64, error) {
	db := GetData()
	list := make([]*UserFile, 0)
	var sum int64
	if err := db.DB(ctx).Where("user_uuid=?", userFile.UserUuid).Count(&sum).Find(&list).Error; err != nil {
		log.Logger.Error(errcode.WithCode(errcode.Database_err, "数据库错误"))
		return nil, 0, errcode.WithCode(errcode.Database_err, "数据库错误")
	}
	return list, sum, nil
}

func DeleteUserFile(ctx context.Context, userFile domain.UserFile) error {
	db := GetData()
	if err := db.DB(ctx).Where("file_hash = ?and user_uuid = ?", userFile.FileHash, userFile.UserUuid).Delete(&File{}).Error; err != nil {
		log.Logger.Error(errcode.WithCode(errcode.Database_err, "数据库错误"))
		return errcode.WithCode(errcode.Database_err, "数据库错误")
	}
	return nil
}
