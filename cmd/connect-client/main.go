package main

import (
	"context"
	"log"
	"net/http"

	"github.com/owenthereal/gostart/proto"
	userv1 "github.com/owenthereal/gostart/proto/gen/user/v1"
	"github.com/owenthereal/gostart/proto/gen/user/v1/v1connect"

	"github.com/bufbuild/connect-go"
)

func main() {
	client := v1connect.NewUserServiceClient(
		http.DefaultClient,
		"http://localhost:8080",
		connect.WithInterceptors(proto.NewAuthInterceptor()),
	)
	res, err := client.CreateUser(
		context.Background(),
		connect.NewRequest(&userv1.CreateUserRequest{
			User: &userv1.NewUser{
				Email: "o@owenou.com",
			},
		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(res.Msg.User)

	resp, err := client.ListUsers(context.Background(), connect.NewRequest(&userv1.ListUsersRequest{}))
	if err != nil {
		log.Fatal(err)
	}

	log.Println(resp.Msg.Users)
}
