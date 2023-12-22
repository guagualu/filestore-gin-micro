package domain

import (
	"gorm.io/gorm"
	"time"
)

type File struct {
	ID       uint
	FileHash string
	FileName string
	FileSize int64
	FileAddr string
	CreateAt time.Time
	UpdateAt time.Time
	Status   gorm.DeletedAt
}
