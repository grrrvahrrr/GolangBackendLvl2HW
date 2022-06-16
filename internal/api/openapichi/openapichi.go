package openapichi

import (
	"GoBeLvl2/internal/api/handlers"
	"GoBeLvl2/internal/entities"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type OpenApiChi struct {
	*chi.Mux
	hs *handlers.Handlers
}

type RouterUser struct {
	User entities.User
}

type RouterEnv struct {
	Project      entities.Project
	Organization entities.Organization
	CorpGroup    entities.CorpGroup
	Community    entities.Community
}

func NewOpenApiRouter(hs *handlers.Handlers) *OpenApiChi {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	ret := &OpenApiChi{
		hs: hs,
	}

	r.Mount("/", Handler(ret))
	swg, err := GetSwagger()
	if err != nil {
		log.Fatal("swagger fail")
	}

	r.Get("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		enc := json.NewEncoder(w)
		_ = enc.Encode(swg)
	})

	ret.Mux = r

	return ret
}

// Add a new enviroment
// (POST /addEnv)
func (rt *OpenApiChi) AddEnv(w http.ResponseWriter, r *http.Request) {
	var re RouterEnv
	err := rt.hs.AddEnv(r.Context(), re)
	if err != nil {
		fmt.Println(err)
	}
}

// Add a new user
// (POST /addUser)
func (rt *OpenApiChi) AddUser(w http.ResponseWriter, r *http.Request) {
	var ru RouterUser
	err := rt.hs.AddUser(r.Context(), ru.User)
	if err != nil {
		fmt.Println(err)
	}
}

// search for env by user names or enviroment params
// (GET /searchEnv)
func (rt *OpenApiChi) SearchEnv(w http.ResponseWriter, r *http.Request, params SearchEnvParams) {
	var paramsstr []string
	i, err := rt.hs.SearchEnv(r.Context(), paramsstr)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(i)

}

// search user by name or enviroment
// (GET /searchUser)
func (rt *OpenApiChi) SearchUser(w http.ResponseWriter, r *http.Request, params SearchUserParams) {
	var paramsstr []string
	user, err := rt.hs.SearchUser(r.Context(), paramsstr)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(user)
}

// set/delete user in enviroment
// (POST /updateUserInEnv)
func (rt *OpenApiChi) UpdateUserInEnv(w http.ResponseWriter, r *http.Request) {
	var operation string
	var ru RouterUser
	var re RouterEnv
	rt.hs.UpdateUserInEnv(r.Context(), ru.User, operation, re)
}
