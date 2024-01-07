package data

import (
	"context"
	"errors"
	"fileStore/internel/domain"
	"fileStore/internel/pkg/code/errcode"
	"fileStore/log"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"gorm.io/gorm"
	"strconv"
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

func GetFileByFileHash(ctx context.Context, fileHash string) (*domain.File, error) {
	db := GetData()
	var file File
	if err := db.DB(ctx).Where("file_hash = ?", fileHash).First(&file).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.WithCode(errcode.NotFoundFile, "未找到file", nil)
		}
		log.Logger.Error("数据库错误:", err)
		return nil, errcode.WithCode(errcode.Database_err, "数据库错误")
	}
	return &domain.File{
		ID:       file.ID,
		FileHash: file.FileHash,
		FileName: file.FileName,
		FileSize: file.FileSize,
		FileAddr: file.FileAddr,
		CreateAt: file.CreateAt,
		UpdateAt: file.UpdateAt,
		Status:   file.Status,
	}, nil
}
func ListFileInfo(ctx context.Context, fileHashs []string) ([]*domain.File, error) {
	db := GetData()
	var file []File
	if err := db.DB(ctx).Where("file_hash in ?", fileHashs).Find(&file).Error; err != nil {
		log.Logger.Error("数据库错误:", err)
		return nil, err
	}
	res := make([]*domain.File, 0)
	for _, v := range file {
		res = append(res, &domain.File{
			ID:       v.ID,
			FileHash: v.FileHash,
			FileName: v.FileName,
			FileSize: v.FileSize,
			FileAddr: v.FileAddr,
			CreateAt: v.CreateAt,
			UpdateAt: v.UpdateAt,
		})
	}
	return res, nil
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
		log.Logger.Error("数据库错误:", err)
		return errcode.WithCode(errcode.Database_err, "数据库错误")
	}
	return nil
}

func DeleteFile(ctx context.Context, fileHash string) error {
	db := GetData()
	if err := db.DB(ctx).Where("file_hash = ?", fileHash).Delete(&File{}).Error; err != nil {
		log.Logger.Error("数据库错误:", err)
		return errcode.WithCode(errcode.Database_err, "数据库错误")
	}
	return nil
}

func UpdataFileLocated(ctx context.Context, fileHash, located string) error {
	db := GetData()
	if err := db.DB(ctx).Debug().Where("file_hash = ?", fileHash).Table("file").Update("file_addr", located).Error; err != nil {
		log.Logger.Error("UpdataFileLocated err:", err)
		return errcode.WithCode(errcode.Database_err, "数据库错误")
	}
	return nil
}

func SaveFileUploadInfo(upInfo domain.MultipartUploadInfo) error {
	conn := GetData().red.Get()
	defer conn.Close()
	_, err := conn.Do("HSET", "MP_"+upInfo.UploadID, "chunkcount", upInfo.ChunkCount)
	if err != nil {
		return err
	}
	_, err = conn.Do("HSET", "MP_"+upInfo.UploadID, "filehash", upInfo.FileHash)
	if err != nil {
		return err
	}
	_, err = conn.Do("HSET", "MP_"+upInfo.UploadID, "filesize", upInfo.FileSize)
	if err != nil {
		return err
	}
	return nil
}

func SaveFileMpUpload(upInfo domain.MultipartUploadInfo) error {
	conn := GetData().RDB().Get()
	defer conn.Close()
	_, err := conn.Do("HSET", "MP_"+upInfo.UploadID, "chunkindex"+strconv.Itoa(upInfo.ChunkIndex), 1)
	//reply, err := conn.Do("SET", "test1", 1)
	if err != nil {
		return err
	}
	return nil
}

func GetFileMpUploadSum(upInfo domain.MultipartUploadInfo) (int, error) {
	conn := GetData().red.Get()
	defer conn.Close()
	sum := 0
	for i := 1; i <= upInfo.ChunkCount; i++ {
		reply, err := redis.Int(conn.Do("HGET", "MP_"+upInfo.UploadID, "chunkindex"+strconv.Itoa(i)))
		if err != nil || reply != 1 {
			continue
		}
		sum++
	}
	return sum, nil
}

func GetFailedChunk(uploadId string, chunkCount int) ([]int, error) {
	conn := GetData().red.Get()
	defer conn.Close()
	res := make([]int, 0)
	for i := 1; i <= chunkCount; i++ {
		reply, err := redis.Int(conn.Do("HGET", "MP_"+uploadId, "chunkindex"+strconv.Itoa(i)))
		fmt.Println(reply)
		fmt.Println(err)
		if err != nil || reply != 1 {
			res = append(res, i)
			continue
		}
	}
	return res, nil
}
