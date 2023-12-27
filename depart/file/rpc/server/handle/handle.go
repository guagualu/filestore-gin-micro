package handle

import (
	"context"
	pb "fileStore/depart/file/proto"
	"fileStore/internel/data"
	"fileStore/internel/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type FileRpcServiceStruct struct {
	pb.UnimplementedFileServiceServer
}

func (s *FileRpcServiceStruct) SaveFile(ctx context.Context, in *pb.FileReq) (*emptypb.Empty, error) {
	//先检查是否file已经存在
	existFile, _ := data.GetFileByFileHash(ctx, in.FileHash)
	if existFile != nil {
		return &emptypb.Empty{}, status.New(codes.OK, "").Err()
	}
	//存入
	file := domain.File{
		FileHash: in.FileHash,
		FileName: in.FileName,
		FileSize: in.FileSize,
		FileAddr: in.LocatedAt,
	}
	err := data.SaveFile(ctx, file)

	if err != nil {
		return &emptypb.Empty{}, status.New(codes.Aborted, "").Err()
	}
	return &emptypb.Empty{}, status.New(codes.OK, "").Err()
}

func (s *FileRpcServiceStruct) DeleteFile(ctx context.Context, in *pb.FileReq) (*emptypb.Empty, error) {
	//存入
	err := data.DeleteFile(ctx, in.FileHash)
	if err != nil {
		return &emptypb.Empty{}, status.New(codes.Aborted, "").Err()
	}
	return &emptypb.Empty{}, status.New(codes.OK, "").Err()
}
