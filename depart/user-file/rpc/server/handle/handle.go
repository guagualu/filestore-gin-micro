package handle

import (
	"context"
	pb "fileStore/depart/user-file/proto"
	"fileStore/internel/data"
	"fileStore/internel/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserFileRpcServiceStruct struct {
	pb.UnimplementedUserFileServiceServer
}

func (s *UserFileRpcServiceStruct) SaveUserFile(ctx context.Context, in *pb.UserFileReq) (*emptypb.Empty, error) {
	//存入
	userFile := domain.UserFile{
		FileHash: in.FileHash,
		FileName: in.FileName,
		UserUuid: in.UserUuid,
	}
	err := data.SaveUserFile(ctx, userFile)
	if err != nil {
		return &emptypb.Empty{}, status.New(codes.Aborted, "").Err()
	}
	return &emptypb.Empty{}, status.New(codes.OK, "").Err()
}

func (s *UserFileRpcServiceStruct) DeleteUserFile(ctx context.Context, in *pb.UserFileReq) (*emptypb.Empty, error) {
	//存入
	userFile := domain.UserFile{
		FileHash: in.FileHash,
		FileName: in.FileName,
		UserUuid: in.UserUuid,
	}
	err := data.DeleteUserFile(ctx, userFile)
	if err != nil {
		return &emptypb.Empty{}, status.New(codes.Aborted, "").Err()
	}
	return &emptypb.Empty{}, status.New(codes.OK, "").Err()
}
