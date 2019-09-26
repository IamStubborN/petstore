package mockdb

import (
	"context"
	"errors"

	"github.com/IamStubborN/petstore/db/models"
)

const testUserName = "admin"

func (d *Database) CreateUser(ctx context.Context, user *models.User) error {
	return nil
}

func (d *Database) CreateUsersFromList(ctx context.Context, list *models.UserList) error {
	return nil
}

func (d *Database) Login(ctx context.Context, username string, password string) (string, error) {
	return "{'GET','POST','PUT','DELETE'}", nil
}

func (d *Database) GetUserByName(ctx context.Context, username string) (*models.User, error) {
	if username == testUserName {
		return &models.User{
			ID:         1,
			Email:      "itu@gmail.com",
			Username:   "admin",
			Password:   "password",
			FirstName:  "User",
			LastName:   "Name",
			Phone:      "+39000000",
			UserStatus: 1,
		}, nil
	}

	return nil, errors.New("invalid username")
}

func (d *Database) UpdateUser(ctx context.Context, user *models.User) error {
	var err error
	if user.Username != testUserName {
		err = errors.New("invalid user")
	}

	return err
}

func (d *Database) DeleteUser(ctx context.Context, username string) error {
	var err error
	if username != testUserName {
		err = errors.New("invalid user")
	}

	return err
}
