package biz

import (
	"context"
	"errors"
	"fileStore/internel/data"
	"fileStore/internel/domain"
	"fileStore/internel/pkg/uuid"
	"fileStore/log"
	"github.com/gomodule/redigo/redis"
)

func SignupUserinfo(ctx context.Context, userName string, userPsd, email, phone string) error {
	return data.CreatUser(ctx, domain.User{
		NickName: userName,
		Email:    email,
		Password: userPsd,
		Mobile:   phone,
	})
}

func SignInUserinfo(ctx context.Context, userPsd, phone string) (*domain.User, error) {
	return data.GetUserByPhoneAndPsd(ctx, domain.User{
		Password: userPsd,
		Mobile:   phone,
	})
}

func GetUserInfo(ctx context.Context, userUuid string) (*domain.User, error) {
	//获取分布式锁结束
	//3、查询用户信息  先查缓存 如果没有 在查mysql 并对redis作缓存
	muxUuid := uuid.NewUuid()
	cancleCtx, cancel := context.WithCancel(ctx)
	err := data.SetMutex(muxUuid, cancleCtx)
	//先删除锁 再 进行watch程序的取消
	defer cancel()
	defer data.DeleteMutex(muxUuid)
	if err != nil {
		log.Logger.Error("分布式锁获取失败")
		return nil, err
	}
	user, err := data.GetUserInfoByCache(ctx, userUuid)
	if err == nil {
		return user, nil
	}
	if err != nil && errors.Is(err, redis.ErrNil) {
		//没有或者是已经过期 重新从mysql中获取
		res, err := data.GetUserInfo(ctx, userUuid)
		if err != nil {
			return nil, err
		}
		//写入缓存
		data.SetUserInfoByCache(ctx, *res)
		return res, nil
	} else {
		return nil, err
	}

}
