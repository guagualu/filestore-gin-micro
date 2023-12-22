package data

import (
	"context"
	"fileStore/internel/domain"
	"fileStore/internel/pkg/code/errcode"
	"fileStore/log"
	"gorm.io/gorm"
	"time"
)

type File struct {
	ID       uint           `gorm:"column:id;type:uint;primary_key;autoIncrement;comment:'物理主键'" json:"id"`
	FileHash string         `gorm:"column:file_hash;type:char(40);Index:idx_file_hash;not null;comment:'文件hash值'" json:"file_hash"`
	FileName string         `gorm:"column:file_name;type:varchar(256);not null;comment:'文件名'"  json:"file_name"`
	FileSize int64          `gorm:"column:file_size;type:int;default: 0;not null;comment:'文件大小'"  json:"file_size"`
	FileAddr string         `gorm:"column:file_addr;type:varchar(256);not null;comment:'文件存储地址'"  json:"file_addr"`
	CreateAt time.Time      `gorm:"column:created_at;type:datetime comment '创建时间';not null;default:CURRENT_TIMESTAMP"  json:"create_at"`
	UpdateAt time.Time      `gorm:"column:updated_at;type:datetime comment '更新时间';not null;default:CURRENT_TIMESTAMP;autoUpdateTime" json:"update_at"`
	Status   gorm.DeletedAt `json:"status"`
}

func SaveFile(ctx context.Context, file domain.File) error {
	db := GetData()
	u := File{
		FileHash: file.FileHash,
		FileName: file.FileName,
		FileSize: file.FileSize,
		FileAddr: file.FileAddr,
	}
	if err := db.DB(ctx).Omit("created_at", "updated_at").Create(&u).Error; err != nil {
		log.Logger.Error(errcode.WithCode(errcode.Database_err, "数据库错误"))
		return errcode.WithCode(errcode.Database_err, "数据库错误")
	}
	return nil
}

func DeleteFile(ctx context.Context, fileHash string) error {
	db := GetData()
	if err := db.DB(ctx).Where("file_hash = ?", fileHash).Delete(&File{}).Error; err != nil {
		log.Logger.Error(errcode.WithCode(errcode.Database_err, "数据库错误"))
		return errcode.WithCode(errcode.Database_err, "数据库错误")
	}
	return nil
}
