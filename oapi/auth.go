package oapi

import (
	"context"
	"fmt"

	"github.com/getkin/kin-openapi/openapi3filter"
)

func NewAuthenticator() openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		if input.SecuritySchemeName != "basicAuth" {
			return fmt.Errorf("security scheme %s != 'basicAuth'", input.SecuritySchemeName)
		}

		username, password, ok := input.RequestValidationInput.Request.BasicAuth()
		if !ok {
			return fmt.Errorf("no basic auth provided")
		}

		if username == "user" && password == "pass" {
			return nil
		}

		return fmt.Errorf("invalid username or password")
	}
}
