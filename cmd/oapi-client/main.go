package main

import (
	"context"
	"log"

	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	"github.com/owenthereal/gostart/oapi"
)

func main() {
	basicAuthProvider, err := securityprovider.NewSecurityProviderBasicAuth("user", "pass")
	if err != nil {
		log.Fatal(err)
	}

	c, err := oapi.NewClientWithResponses(
		"http://localhost:8080",
		oapi.WithRequestEditorFn(basicAuthProvider.Intercept),
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
