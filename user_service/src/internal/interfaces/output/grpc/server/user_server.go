package server

import (
	"context"
	"user_service/internal/usecase"
	pb "user_service/proto"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
	userUsecase usecase.UserUsecase
}

func NewUserServer(uc usecase.UserUsecase) *UserServer {
	return &UserServer{userUsecase: uc}
}

func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	userID := req.GetId()

	user, err := s.userUsecase.GetProfile(int(userID))
	if err != nil {
		return nil, err
	}

	return &pb.GetUserResponse{
		Id:       int32(user.ID),
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
