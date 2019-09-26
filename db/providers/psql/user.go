package psql

import (
	"context"
	"database/sql"

	"github.com/IamStubborN/petstore/db/models"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"golang.org/x/crypto/bcrypt"
)

func (d *Database) CreateUser(ctx context.Context, user *models.User) error {
	encryptedPass, err := encryptPassword(user.Password)
	if err != nil {
		return errors.Wrap(err, "can't generate encrypted password")
	}

	user.Password = encryptedPass

	_, err = d.pool.NamedExecContext(ctx, qm[userCreateQ], &user)
	if err != nil {
		return errors.Wrap(err, "can't insert into user")
	}

	return nil
}

func (d *Database) CreateUsersFromList(ctx context.Context, list *models.UserList) error {
	tx, err := d.pool.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return errors.Wrap(err, "can't start transaction")
	}

	defer func() {
		if err != nil {
			checkError(tx.Rollback)
		}
		checkError(tx.Commit)
	}()

	var stmt *sqlx.NamedStmt
	stmt, err = d.pool.PrepareNamedContext(ctx, qm[userCreateQ])
	if err != nil {
		return errors.Wrap(err, "can't prepare statement")
	}
	defer checkError(stmt.Close)

	for _, user := range *list {
		encryptedPass, err := encryptPassword(user.Password)
		if err != nil {
			return errors.Wrap(err, "can't generate encrypted password")
		}

		user.Password = encryptedPass

		_, err = stmt.ExecContext(ctx, user)
		if err != nil {
			return errors.Wrap(err, "can't insert into user")
		}
	}

	return nil
}

func (d *Database) Login(ctx context.Context, username, password string) (string, error) {
	var encryptedPass, allowMethods string
	err := d.pool.QueryRowxContext(ctx, qm[userGetAllowedMethodsAndPassQ],
		username).Scan(&encryptedPass, &allowMethods)
	if err != nil {
		return allowMethods, errors.Wrap(err, "can't scan data for user")
	}

	err = comparePassword(encryptedPass, password)
	if err != nil {
		return allowMethods, errors.Wrap(err, "can't login, wrong password")
	}

	return allowMethods, nil
}

func (d *Database) GetUserByName(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := d.pool.QueryRowxContext(ctx, qm[userGetByNameQ], username).
		StructScan(&user)
	if err != nil {
		return nil, errors.Wrap(err, "can't get user from db")
	}

	return &user, nil
}

func (d *Database) UpdateUser(ctx context.Context, user *models.User) error {
	encryptedPass, err := encryptPassword(user.Password)
	if err != nil {
		return errors.Wrap(err, "can't generate encrypted password")
	}

	user.Password = encryptedPass

	res, err := d.pool.NamedExecContext(ctx, qm[userUpdateQ], user)
	if err != nil {
		return errors.Wrap(err, "can't exec update user")
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.Errorf("user %s doesn't exist", user.Username)
	}

	return nil
}

func (d *Database) DeleteUser(ctx context.Context, username string) error {
	res, err := d.pool.ExecContext(ctx, qm[userDeleteQ], username)
	if err != nil {
		return errors.Wrap(err, "can't delete user")
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.Errorf("user %s doesn't exist", username)
	}

	return nil
}

func encryptPassword(password string) (string, error) {
	encryptedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(encryptedPass), nil
}

func comparePassword(encryptedPass, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(encryptedPass), []byte(password))
}
