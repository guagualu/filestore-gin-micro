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

type MultipartUploadInfo struct {
	FileHash   string `json:"file_hash"`
	FileSize   int    `json:"file_size"`
	UploadID   string `json:"upload_id"`
	ChunkSize  int    `json:"chunk_size"`
	ChunkCount int    `json:"chunk_count"`
	ChunkIndex int    `json:"chunk_index"`
}
