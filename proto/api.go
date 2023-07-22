//go:generate buf generate

package proto

import (
	"context"

	userv1 "github.com/owenthereal/gostart/proto/gen/user/v1"
)

var _ userv1.UserServiceServer = (*UserService)(nil)

type UserService struct {
	userv1.UnimplementedUserServiceServer
}

func (s *UserService) CreateUser(ctx context.Context, req *userv1.CreateUserRequest) (*userv1.CreateUserResponse, error) {
	return nil, nil
}
