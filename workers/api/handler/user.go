package handler

import (
	"fmt"

	"gopkg.in/validator.v2"

	"io/ioutil"
	"net/http"
	"strings"

	"github.com/IamStubborN/petstore/db"
	"github.com/IamStubborN/petstore/db/models"
	"github.com/IamStubborN/petstore/workers/api/auth"
	"github.com/IamStubborN/petstore/workers/api/mware"
	"github.com/go-chi/chi"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/pkg/errors"
)

func UserHandlers(r chi.Router) {
	registerGroup := r.Group(nil)
	registerGroup.Get("/login", login)
	registerGroup.Post("/", createUser)

	securedGroup := r.Group(nil)
	securedGroup.Use(mware.JWT)
	securedGroup.Get("/logout", logout)
	securedGroup.Post("/createWithList", createUsersFromList)
	securedGroup.Get("/{username}", getUserByName)
	securedGroup.Put("/{username}", updateUserByName)
	securedGroup.Delete("/{username}", deleteUserByName)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := genContext(r)
	defer cancel()

	if r.Body == nil {
		respond(w, errors.New("request body is nil"),
			http.StatusMethodNotAllowed, "invalid input")
		return
	}
	defer checkErrors(r.Body.Close)

	var user models.User
	bytesBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respond(w, errors.Wrap(err, "can't read request body"),
			http.StatusInternalServerError, "providers server error")
		return
	}

	if err = user.UnmarshalJSON(bytesBody); err != nil {
		respond(w, errors.Wrap(err, "can't decode request body to user"),
			http.StatusInternalServerError, "providers server error")
		return
	}

	if err = validator.Validate(user); err != nil {
		respond(w, errors.Wrap(err, "can't validate user from body"),
			http.StatusMethodNotAllowed, "invalid input")
		return
	}

	userDI := db.GetUserDI()
	if err = userDI.CreateUser(ctx, &user); err != nil {
		respond(w, err, http.StatusBadRequest, "invalid input")
		return
	}

	respond(w, nil, http.StatusOK, "success")
}

func createUsersFromList(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := genContext(r)
	defer cancel()

	if r.Body == nil {
		respond(w, errors.New("request body is nil"),
			http.StatusMethodNotAllowed, "invalid input")
		return
	}
	defer checkErrors(r.Body.Close)

	bytesBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respond(w, errors.Wrap(err, "can't read request body"),
			http.StatusInternalServerError, "providers server error")
		return
	}

	users := &models.UserList{}
	if err = users.UnmarshalJSON(bytesBody); err != nil {
		respond(w, errors.Wrap(err, "can't decode request body to users list"),
			http.StatusInternalServerError, "providers server error")
		return
	}

	for _, user := range *users {
		if err = validator.Validate(user); err != nil {
			respond(w, errors.Wrap(err, "can't validate user from body"),
				http.StatusMethodNotAllowed, "invalid input")
			return
		}
	}

	userDI := db.GetUserDI()
	if err := userDI.CreateUsersFromList(ctx, users); err != nil {
		respond(w, err, http.StatusMethodNotAllowed, "invalid input")
		return
	}

	respond(w, nil, http.StatusOK, "success")

}

func login(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := genContext(r)
	defer cancel()

	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")

	if isNotValid(username) && isNotValid(password) {
		respond(w, fmt.Errorf("bad input %s or %s", username, password),
			http.StatusBadRequest, "invalid username/password supplied")
		return
	}

	userDI := db.GetUserDI()
	user, err := userDI.GetUserByName(ctx, username)
	if err != nil {
		respond(w, err, http.StatusBadRequest, "invalid username/password supplied")
		return
	}

	allowMethods, err := userDI.Login(ctx, username, password)
	if err != nil {
		respond(w, err, http.StatusBadRequest, "invalid username/password supplied")
		return
	}

	sessionID, err := uuid.NewV4()
	if err != nil {
		respond(w, err, http.StatusInternalServerError, "providers server error")
		return
	}

	token, err := auth.GenerateToken(sessionID.String(), user.ID, allowMethods)
	if err != nil {
		respond(w, err, http.StatusInternalServerError, "providers server error")
		return
	}

	w.Header().Add("Authorization", "Bearer "+token)

	respond(w, nil, http.StatusOK, "logged in user session:"+sessionID.String())
}

func logout(w http.ResponseWriter, r *http.Request) {
	bearer := r.Header.Get("Authorization")
	if bearer == "" {
		respond(w, errors.New("authorization token not set"),
			http.StatusBadRequest, "invalid username supplied")
		return
	}
	token := bearer[len("Bearer "):]

	if _, err := auth.ParseToken(token); err != nil {
		respond(w, errors.Wrap(err, "Authorization token invalid"),
			http.StatusBadRequest, "invalid username supplied")
		return
	}

	auth.AddToBlackList(token)

	respond(w, nil, http.StatusOK, "success")
}

func getUserByName(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := genContext(r)
	defer cancel()

	username := strings.TrimPrefix(r.URL.Path, "/api/v2/user/")
	if isNotValid(username) {
		respond(w, fmt.Errorf("bad input %s", username),
			http.StatusBadRequest, "invalid username supplied")
		return
	}

	userDI := db.GetUserDI()
	user, err := userDI.GetUserByName(ctx, username)
	if err != nil {
		respond(w, err, http.StatusNotFound, "user not found")
		return
	}

	data, err := user.MarshalJSON()
	if err != nil {
		respond(w, errors.Wrap(err, "can't marshal data"),
			http.StatusInternalServerError, "providers error")
		return
	}

	if _, err := w.Write(data); err != nil {
		respond(w, errors.Wrap(err, "can't write response"),
			http.StatusInternalServerError, "providers error")
		return
	}
}

func updateUserByName(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := genContext(r)
	defer cancel()

	if r.Body == nil {
		respond(w, errors.New("request body is nil"),
			http.StatusBadRequest, "invalid user supplied")
		return
	}
	defer checkErrors(r.Body.Close)

	bytesBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respond(w, errors.Wrap(err, "can't read request body"),
			http.StatusInternalServerError, "providers server error")
		return
	}

	var user models.User
	if err = user.UnmarshalJSON(bytesBody); err != nil {
		respond(w, errors.Wrap(err, "can't decode request body to user"),
			http.StatusInternalServerError, "providers server error")
		return
	}

	user.Username = strings.TrimPrefix(r.URL.Path, "/api/v2/user/")

	if err := validator.Validate(user); err != nil {
		respond(w, errors.Wrap(err, "can't validate user from body"),
			http.StatusBadRequest, "invalid username supplied")
		return
	}

	userDI := db.GetUserDI()
	if err := userDI.UpdateUser(ctx, &user); err != nil {
		respond(w, err, http.StatusNotFound, "user not found")
		return
	}

	respond(w, nil, http.StatusOK, "success")
}

func deleteUserByName(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := genContext(r)
	defer cancel()

	username := strings.TrimPrefix(r.URL.Path, "/api/v2/user/")

	if isNotValid(username) {
		respond(w, fmt.Errorf("bad input %s", username),
			http.StatusBadRequest, "invalid username supplied")
		return
	}

	userDI := db.GetUserDI()
	if err := userDI.DeleteUser(ctx, username); err != nil {
		respond(w, err, http.StatusBadRequest, "user not found")
		return
	}

	respond(w, nil, http.StatusOK, "success")
}
