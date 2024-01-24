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
	if err == nil && user != nil {
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

func AddFriend(ctx context.Context, userAMobile, userBMobile string) error {
	return data.CreatFriendShip(ctx, domain.Friends{
		UserAMobile: userAMobile,
		UserBMobile: userBMobile,
	})

}

func GetUserFriends(ctx context.Context, userMobile string) ([]domain.User, error) {
	//先获取userFriends的mobile
	friends, err := data.GetUserFriendsByUserPhone(ctx, userMobile)
	if err != nil {
		return nil, err
	}
	mobiles := make([]string, 0)
	for _, v := range friends {
		var friendPhone string
		if v.UserAMobile != userMobile {
			friendPhone = v.UserAMobile
		} else {
			friendPhone = v.UserBMobile
		}
		mobiles = append(mobiles, friendPhone)
	}
	return data.ListUserInfoByMobile(ctx, mobiles)
}
