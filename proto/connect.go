//go:generate buf generate

package proto

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"sync"

	"github.com/bufbuild/connect-go"
	userv1 "github.com/owenthereal/gostart/proto/gen/user/v1"
	"github.com/owenthereal/gostart/proto/gen/user/v1/v1connect"
)

var _ v1connect.UserServiceHandler = (*ConnectUserService)(nil)

func NewConnectUserService() *ConnectUserService {
	return &ConnectUserService{
		Users:  make(map[int64]*userv1.User),
		NextId: 1000,
	}
}

type ConnectUserService struct {
	Users  map[int64]*userv1.User
	NextId int64
	Lock   sync.Mutex
}

func (s *ConnectUserService) CreateUser(ctx context.Context, req *connect.Request[userv1.CreateUserRequest]) (*connect.Response[userv1.CreateUserResponse], error) {
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

func (s *ConnectUserService) ListUsers(ctx context.Context, req *connect.Request[userv1.ListUsersRequest]) (*connect.Response[userv1.ListUsersResponse], error) {
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

const authHeader = "Authentication"

func NewAuthInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			if req.Spec().IsClient {
				req.Header().Set(authHeader, fmt.Sprintf("Basic: "+base64.StdEncoding.EncodeToString([]byte("user:pass"))))
				return next(ctx, req)
			}

			auth := req.Header().Get(authHeader)
			if auth == "" {
				return nil, connect.NewError(
					connect.CodeUnauthenticated,
					errors.New("no basic auth provided"),
				)
			}

			return next(ctx, req)
		})
	}

	return connect.UnaryInterceptorFunc(interceptor)
}
