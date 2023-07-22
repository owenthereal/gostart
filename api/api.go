//go:generate oapi-codegen --config=cfg.yaml ./api.yaml

package api

import "net/http"

var _ ServerInterface = (*Server)(nil)

type Server struct {
}

func (s *Server) GetUsers(w http.ResponseWriter, r *http.Request) {
}
