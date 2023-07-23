package main

import (
	"context"
	"log"

	"github.com/owenthereal/gostart/oapi"
)

func main() {
	c, err := oapi.NewClientWithResponses("http://localhost:8080")
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
