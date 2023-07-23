//go:generate buf generate

package proto

import (
	"context"
	"sync"

	"github.com/bufbuild/connect-go"
	userv1 "github.com/owenthereal/gostart/proto/gen/user/v1"
	"github.com/owenthereal/gostart/proto/gen/user/v1/v1connect"
)

var _ v1connect.UserServiceHandler = (*UserService)(nil)

func NewUserService() *UserService {
	return &UserService{
		Users:  make(map[int64]*userv1.User),
		NextId: 1000,
	}
}

type UserService struct {
	userv1.UnimplementedUserServiceServer

	Users  map[int64]*userv1.User
	NextId int64
	Lock   sync.Mutex
}

func (s *UserService) CreateUser(ctx context.Context, req *connect.Request[userv1.CreateUserRequest]) (*connect.Response[userv1.CreateUserResponse], error) {
	if err := req.Msg.ValidateAll(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	s.Lock.Lock()
	defer s.Lock.Unlock()

	user := &userv1.User{
		Id:    s.NextId,
		Email: req.Msg.User.Email, // TODO: check if user exists
	}
	s.Users[user.Id] = user
	s.NextId = s.NextId + 1

	return connect.NewResponse(&userv1.CreateUserResponse{
		User: &userv1.User{
			Email: req.Msg.User.Email,
		},
	}), nil
}

func (s *UserService) ListUsers(ctx context.Context, req *connect.Request[userv1.ListUsersRequest]) (*connect.Response[userv1.ListUsersResponse], error) {
	if err := req.Msg.ValidateAll(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	s.Lock.Lock()
	defer s.Lock.Unlock()

	result := []*userv1.User{}
	for _, user := range s.Users {
		if emails := req.Msg.Emails; emails != nil {
			for _, t := range emails {
				if user.Email == t {
					result = append(result, user)
				}
			}
		} else {
			result = append(result, user)
		}

		if limit := req.Msg.Limit; limit != nil {
			l := int(*limit)
			if len(result) >= l {
				break
			}
		}
	}

	return connect.NewResponse(&userv1.ListUsersResponse{
		Users: result,
	}), nil
}
