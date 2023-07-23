package main

import (
	"context"
	"log"
	"net/http"

	"github.com/owenthereal/gostart/oapi"
)

func main() {
	c, err := oapi.NewClientWithResponses(
		"http://localhost:8080",
		oapi.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			req.SetBasicAuth("user", "pass")
			return nil
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.TODO()

	user, err := c.CreateUserWithResponse(ctx, oapi.CreateUserJSONRequestBody{
		Email: "o@owenou.com",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(user.JSON201)

	users, err := c.FindUsersWithResponse(ctx, &oapi.FindUsersParams{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(users.JSON200)
}
