package main

import (
	"context"
	"fmt"

	"github.com/yourusername/kitex-gateway-demo/kitex_gen/api"
)

type UserServiceImpl struct{}

func (s *UserServiceImpl) GetUser(ctx context.Context, req *api.GetUserRequest) (*api.GetUserResponse, error) {
	return &api.GetUserResponse{
		User: &api.User{
			Id:   req.Id,
			Name: "Kitex User",
		},
	}, nil
}

func (s *UserServiceImpl) StreamUsers(stream api.UserService_StreamUsersServer) error {
	for i := 1; i <= 5; i++ {
		err := stream.Send(&api.GetUserResponse{
			User: &api.User{
				Id:   int64(i),
				Name: fmt.Sprintf("User %d", i),
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}
