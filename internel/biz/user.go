package biz

import (
	"context"
	"fileStore/internel/data"
	"fileStore/internel/domain"
)

func SignupUserinfo(ctx context.Context, userName string, userPsd, email, phone string) error {
	return data.CreatUser(ctx, domain.User{
		NickName: userName,
		Email:    email,
		Password: userPsd,
		Mobile:   phone,
	})
}
