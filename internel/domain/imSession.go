package domain

import (
	"sync"
	"time"
)

type ImSessionContent struct {
	Id             uint      `json:"id"`
	SessionUuid    string    `json:"session_uuid"`
	SendUserUuid   string    `json:"send_user_uuid"` // 发送者uuid
	ToUserUuid     string    `json:"to_user_uuid"`   // 接受者uuid
	MessageType    int       `json:"message_type"`   //0为普通文本消息，1为文件消息 content存的是文件hash
	MessageContent string    `json:"message_content"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type ImSession struct {
	Id          uint      `json:"id"`           // 物理主键
	SessionUuid string    `json:"session_uuid"` // 业务主键
	UserAUuid   string    `json:"user_a_uuid"`  // 发送者uuid
	UserBUuid   string    `json:"user_b_uuid"`  // 接受者uuid
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ImSendMsg struct {
	SessionUuid  string `json:"session_uuid"`   // 业务主键
	SendUserUuid string `json:"send_user_uuid"` // 发送者uuid
	ToUserUuid   string `json:"to_user_uuid"`
	Message      string `json:"message"`
	MessageType  int    `json:"message_type"`
}

// singel 单例
var imChannelMap sync.Map //并发安全的map
var once sync.Once

func GetImChannelMap() sync.Map {

	return imChannelMap
}

func GetKey(key string) chan ImSendMsg {
	//先查询是否key已经存在，不存在就建立，存在就直接获取
	val, ok := imChannelMap.Load(key)
	if !ok {
		return nil
	}
	return val.(chan ImSendMsg)
}

func StoreKey(key string, value chan ImSendMsg) {
	imChannelMap.Store(key, value)
}

func DeleteKey(key string) {
	imChannelMap.Delete(key)
	return
}
