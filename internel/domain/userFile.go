package domain

import (
	"gorm.io/gorm"
	"time"
)

type UserFile struct {
	ID       uint
	FileHash string
	UserUuid string
	FileName string
	CreateAt time.Time
	UpdateAt time.Time
	Status   gorm.DeletedAt
}
