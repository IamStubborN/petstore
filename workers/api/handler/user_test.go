package handler

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/IamStubborN/petstore/config"
	"github.com/IamStubborN/petstore/db"
	"github.com/IamStubborN/petstore/workers/api/auth"
	"github.com/stretchr/testify/assert"

	"github.com/go-chi/chi"
)

const host = "http://localhost:5555"

var jwtToken string

//nolint:gochecknoinits
func init() {
	db.InitDatabase(&config.Config{DB: config.DB{Provider: "mockdb"}})
	auth.InitJWTAuth(&config.Config{JWT: config.JWT{
		KeysPath: "../../../",
		TTL:      time.Minute,
	}})
}

func TestHandler_createUser(t *testing.T) {
	raw := `
{
  "id": 1,
  "user_name": "admin",
  "first_name": "qer",
  "last_name": "qerer",
  "email": "test1@grrmail.com",
  "password": "password",
  "phone": "+3809555555",
  "user_status_id": 1
}`
	request, err := http.NewRequest("POST", "/api/v2/user", strings.NewReader(raw))
	if err != nil {
		log.Println(err)
	}

	request = addToCtxWriteTimeout(request)
	response := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post("/api/v2/user", createUser)

	r.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "OK response is expected")
}

func TestHandler_createUsersFromList(t *testing.T) {
	raw := `
[
  {
  "id": 0,
  "user_name": "TestUser1",
  "first_name": "Test1",
  "last_name": "User1",
  "email": "test@gmail.com1",
  "password": "test1test",
  "phone": "+13809555555",
  "user_status_id": 1
},
  {
    "id": 0,
    "user_name": "TestUser2",
    "first_name": "Test2",
    "last_name": "User2",
    "email": "test@gmail.com2",
    "password": "test2test",
    "phone": "+23809555555",
    "user_status_id": 2
  },
  {
    "id": 0,
    "user_name": "TestUser3",
    "first_name": "Test3",
    "last_name": "User3",
    "email": "test@gmail.com3",
    "password": "test3test",
    "phone": "+33809555555",
    "user_status_id": 3
  }
]`

	request, err := http.NewRequest("POST", "/api/v2/user/createWithList", strings.NewReader(raw))
	if err != nil {
		log.Println(err)
	}

	request = addToCtxWriteTimeout(request)
	response := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post("/api/v2/user/createWithList", createUsersFromList)

	r.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "OK response is expected")
}

func TestHandler_getUserByName(t *testing.T) {
	request, err := http.NewRequest("GET", "/api/v2/user/admin", nil)
	if err != nil {
		log.Println(err)
	}

	request = addToCtxWriteTimeout(request)
	response := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/api/v2/user/{username}", getUserByName)

	r.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "OK response is expected")
}

func TestHandler_login(t *testing.T) {
	reqArgs := url.Values{}
	reqArgs.Add("username", "admin")
	reqArgs.Add("password", "password")
	reqURL, err := url.Parse(host)
	if err != nil {
		log.Println(err)
	}

	reqURL.Path = "/api/v2/user/login"
	reqURL.RawQuery = reqArgs.Encode()

	request, err := http.NewRequest("GET", reqURL.String(), nil)
	if err != nil {
		log.Println(err)
	}

	request = addToCtxWriteTimeout(request)
	response := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/api/v2/user/login", login)

	r.ServeHTTP(response, request)
	jwtToken = response.Header().Get("Authorization")
	assert.Equal(t, 200, response.Code, "OK response is expected")
}

func TestHandler_logout(t *testing.T) {
	request, err := http.NewRequest("GET", "/api/v2/user/logout", nil)
	if err != nil {
		log.Println(err)
	}

	request = addToCtxWriteTimeout(request)
	request.Header.Add("Authorization", jwtToken)
	response := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/api/v2/user/logout", logout)

	r.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "OK response is expected")
}

func TestHandler_updateUserByName(t *testing.T) {
	raw := `
{
  "id": 0,
  "user_name": "admin",
  "first_name": "qwer3333",
  "last_name": "User12",
  "email": "test@gmail.com",
  "password": "password",
  "phone": "+3809555555111",
  "user_status_id": 2
}`

	request, err := http.NewRequest("POST", "/api/v2/user/admin", strings.NewReader(raw))
	if err != nil {
		log.Println(err)
	}

	request = addToCtxWriteTimeout(request)
	response := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post("/api/v2/user/{username}", updateUserByName)

	r.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "OK response is expected")
}

func TestHandler_deleteUserByName(t *testing.T) {
	request, err := http.NewRequest("DELETE", "/api/v2/user/admin", nil)
	if err != nil {
		log.Println(err)
	}

	request = addToCtxWriteTimeout(request)
	response := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Delete("/api/v2/user/{username}", deleteUserByName)

	r.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "OK response is expected")
}

func addToCtxWriteTimeout(r *http.Request) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), http.ServerContextKey,
		&http.Server{WriteTimeout: 10 * time.Second}))
}
