//go:generate oapi-codegen --config=cfg.yaml ./oapi.yaml

package oapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

var _ ServerInterface = (*UserService)(nil)

func NewUserService() *UserService {
	return &UserService{
		Users: make(map[uuid.UUID]User),
	}
}

type UserService struct {
	Users map[uuid.UUID]User
	Lock  sync.Mutex
}

func (s *UserService) CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser NewUser
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		renderError(w, r, http.StatusBadRequest, "Invalid format for NewUser")
		return
	}

	s.Lock.Lock()
	defer s.Lock.Unlock()

	var user User
	user.Email = newUser.Email // TODO: check for uniqueness
	user.Id = uuid.New()

	s.Users[user.Id] = user

	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, user)
}

func (s *UserService) FindUsers(w http.ResponseWriter, r *http.Request, params FindUsersParams) {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	result := []User{}
	for _, user := range s.Users {
		if params.Emails != nil {
			for _, t := range *params.Emails {
				if user.Email == openapi_types.Email(t) {
					result = append(result, user)
				}
			}
		} else {
			result = append(result, user)
		}

		if params.Limit != nil {
			l := int(*params.Limit)
			if len(result) >= l {
				break
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, result)
}

func (s *UserService) DeleteUser(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	_, found := s.Users[id]
	if !found {
		renderError(w, r, http.StatusNotFound, fmt.Sprintf("Could not find user with ID %d", id))
		return
	}
	delete(s.Users, id)

	w.WriteHeader(http.StatusNoContent)
}

func (s *UserService) GetUserById(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	user, found := s.Users[id]
	if !found {
		renderError(w, r, http.StatusNotFound, fmt.Sprintf("Could not find user with ID %d", id))
		return
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, user)
}

func renderError(w http.ResponseWriter, r *http.Request, code int, message string) {
	petErr := Error{
		Code:    int32(code),
		Message: message,
	}
	w.WriteHeader(code)
	render.JSON(w, r, petErr)
}
