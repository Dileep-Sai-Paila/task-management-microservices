package grpcclient

import (
	"context"
	"fmt"
	pb "task_service/proto"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// client for the User grpc service.
type UserClient struct {
	client pb.UserServiceClient
}

func NewUserClient(address string) (*UserClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("could not connect to user service: %w", err)
	}

	client := pb.NewUserServiceClient(conn)

	return &UserClient{client: client}, nil
}

// calls the GetUser rpc on the User Service.
func (c *UserClient) GetUser(ctx context.Context, userID int32) (*pb.GetUserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	req := &pb.GetUserRequest{Id: userID}

	res, err := c.client.GetUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("grpc call to GetUser failed: %w", err)
	}

	return res, nil
}
