package data

import (
	"gorm.io/gorm"
	"time"
)

type File struct {
	ID       uint           `gorm:"column:id;type:uint;primary_key;comment:'物理主键'" json:"id"`
	FileHash string         `gorm:"column:file_hash;type:char(40);unique;Index:idx_file_hash;not null;comment:'文件hash值'" json:"file_sha1"`
	FileName string         `gorm:"column:file_name;type:varchar(256);not null;comment:'文件名'"  json:"file_name"`
	FileSize int64          `gorm:"column:file_size;type:tinyint;default: 0;not null;comment:'文件大小'"  json:"file_size"`
	FileAddr string         `gorm:"column:file_addr;type:varchar(256);not null;comment:'文件存储地址'"  json:"file_addr"`
	CreateAt time.Time      `gorm:"column:created_at;type:datetime comment '创建时间';not null;default:CURRENT_TIMESTAMP"  json:"create_at"`
	UpdateAt time.Time      `gorm:"column:updated_at;type:datetime comment '更新时间';not null;default:CURRENT_TIMESTAMP;autoUpdateTime" json:"update_at"`
	Status   gorm.DeletedAt `json:"status"`
}
