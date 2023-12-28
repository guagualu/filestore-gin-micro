package handle

import (
	"context"
	"errors"
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

func (s *UserFileRpcServiceStruct) GetUserFile(ctx context.Context, in *pb.UserFileReq) (*pb.UserFileRsp, error) {
	userFile := domain.UserFile{
		FileHash: in.FileHash,
		FileName: in.FileName,
		UserUuid: in.UserUuid,
	}
	res, err := data.GetUserFiles(ctx, userFile)
	if err != nil {
		return nil, errors.New("test")
	}
	return &pb.UserFileRsp{
		FileHash:  res.FileHash,
		FileName:  res.FileName,
		UserUuid:  res.UserUuid,
		Id:        uint32(res.ID),
		CreatedAt: res.CreateAt.String(),
		UpdatedAt: res.UpdateAt.String(),
	}, nil
}
