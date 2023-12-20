package handle

import (
	"context"
	"fileStore/depart/user/proto"
)

type UserRpcServiceStruct struct {
	pb.UnimplementedUserServiceServer
}

func (s *UserRpcServiceStruct) UserInfo(ctx context.Context, in *pb.ReqUserInfo) (*pb.RespUserInfo, error) {
	return &pb.RespUserInfo{
		Id:       0,
		Uuid:     "123",
		NickName: "423",
		Email:    "2354",
		Password: "234",
		Mobile:   "423",
	}, nil
}
