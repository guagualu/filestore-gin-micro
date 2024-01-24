package handle

import (
	"context"
	"fileStore/depart/user/proto"
	"fileStore/internel/biz"
)

type UserRpcServiceStruct struct {
	pb.UnimplementedUserServiceServer
}

func (s *UserRpcServiceStruct) UserInfo(ctx context.Context, in *pb.ReqUserInfo) (*pb.RespUserInfo, error) {
	user, err := biz.GetUserInfo(ctx, in.Uuid)
	if err != nil {
		return nil, err
	}
	return &pb.RespUserInfo{
		Id:       int32(user.Id),
		Uuid:     user.Uuid,
		NickName: user.NickName,
		Email:    user.Email,
		Password: user.Password,
		Mobile:   user.Mobile,
	}, nil
}
