package psql

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/IamStubborN/petstore/db/models"
	"github.com/jmoiron/sqlx"
)

func TestDatabase_CreateUser(t *testing.T) {
	pool, mock := generateMockedDB()
	defer checkError(pool.Close)

	type fields struct {
		pool *sqlx.DB
	}
	type args struct {
		ctx    context.Context
		user   *models.User
		mockFn func()
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Success 1",
			fields: fields{pool: pool},
			args: args{
				ctx: context.Background(),
				user: &models.User{
					ID:         1,
					Email:      "itu@gmail.com",
					Username:   "username",
					Password:   "password",
					FirstName:  "User",
					LastName:   "Name",
					Phone:      "+39000000",
					UserStatus: 1,
				},
				mockFn: func() {
					mock.ExpectExec(`insert into "user" (.+) values (.+)`).
						WillReturnResult(sqlmock.NewResult(1, 1))
				},
			},
			wantErr: false,
		},
		{
			name:   "Failure 1",
			fields: fields{pool: pool},
			args: args{
				ctx: context.Background(),
				user: &models.User{
					ID:         1,
					Email:      "itu@gmail.com",
					Username:   "username",
					Password:   "password",
					FirstName:  "User",
					LastName:   "Name",
					Phone:      "+39000000",
					UserStatus: 1,
				},
				mockFn: func() {
					mock.ExpectExec(`insert into "user" (.+) values (.+)`).
						WillReturnError(errors.New("can't insert into user: pg duplicate user_name"))
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Database{
				pool: tt.fields.pool,
			}

			tt.args.mockFn()

			if err := d.CreateUser(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDatabase_CreateUsersFromList(t *testing.T) {
	pool, mock := generateMockedDB()
	defer checkError(pool.Close)

	type fields struct {
		pool *sqlx.DB
	}
	type args struct {
		ctx    context.Context
		list   *models.UserList
		mockFn func()
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Success 1",
			fields: fields{pool: pool},
			args: args{
				ctx: context.Background(),
				list: &models.UserList{
					{
						ID:         1,
						Email:      "itu@gmail.com",
						Username:   "username",
						Password:   "password",
						FirstName:  "User",
						LastName:   "Name",
						Phone:      "+390000001",
						UserStatus: 1,
					},
					{
						ID:         2,
						Email:      "itu2@gmail.com",
						Username:   "username2",
						Password:   "password2",
						FirstName:  "User2",
						LastName:   "Name2",
						Phone:      "+390000002",
						UserStatus: 2,
					},
				},
				mockFn: func() {
					mock.ExpectBegin()
					mock.ExpectPrepare(`insert into "user" (.+) values (.+)`)
					mock.ExpectExec(`insert into "user" (.+) values (.+)`).
						WillReturnResult(sqlmock.NewResult(1, 1))
					mock.ExpectExec(`insert into "user" (.+) values (.+)`).
						WillReturnResult(sqlmock.NewResult(2, 2))
					mock.ExpectCommit()
				},
			},
			wantErr: false,
		},
		{
			name:   "Failure 1",
			fields: fields{pool: pool},
			args: args{
				ctx:  context.Background(),
				list: nil,
				mockFn: func() {
					mock.ExpectBegin()
					mock.ExpectPrepare(`insert into "user" (.+) values (.+)`).
						WillReturnError(errors.New("can't prepare statement"))
					mock.ExpectRollback()
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Database{
				pool: tt.fields.pool,
			}

			tt.args.mockFn()

			if err := d.CreateUsersFromList(tt.args.ctx, tt.args.list); (err != nil) != tt.wantErr {
				t.Errorf("CreateUsersFromList() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDatabase_DeleteUser(t *testing.T) {
	pool, mock := generateMockedDB()
	defer checkError(pool.Close)

	type fields struct {
		pool *sqlx.DB
	}
	type args struct {
		ctx      context.Context
		username string
		mockFn   func()
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Success 1",
			fields: fields{pool: pool},
			args: args{
				ctx:      context.Background(),
				username: "DelUser",
				mockFn: func() {
					mock.ExpectExec(`delete from "user" where user_name=.`).
						WithArgs("DelUser").
						WillReturnResult(sqlmock.NewResult(0, 1))
				},
			},
			wantErr: false,
		},
		{
			name:   "Failure 1",
			fields: fields{pool: pool},
			args: args{
				ctx:      context.Background(),
				username: "DelUser",
				mockFn: func() {
					mock.ExpectExec(`delete from "user" where user_name=.`).
						WithArgs("DelUser").
						WillReturnResult(sqlmock.NewResult(0, 0))
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Database{
				pool: tt.fields.pool,
			}

			tt.args.mockFn()

			if err := d.DeleteUser(tt.args.ctx, tt.args.username); (err != nil) != tt.wantErr {
				t.Errorf("DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDatabase_GetUserByName(t *testing.T) {
	pool, mock := generateMockedDB()
	defer checkError(pool.Close)

	type fields struct {
		pool *sqlx.DB
	}
	type args struct {
		ctx      context.Context
		username string
		mockFn   func()
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.User
		wantErr bool
	}{
		{
			name:   "Success 1",
			fields: fields{pool: pool},
			args: args{
				ctx:      context.Background(),
				username: "peterpandam",
				mockFn: func() {
					mock.ExpectQuery(`select (.+) from "user"`).
						WithArgs("peterpandam").
						WillReturnRows(sqlmock.NewRows([]string{
							"id", "user_name", "password", "first_name",
							"last_name", "email", "phone", "user_status_id"}).
							AddRow(1, "peterpandam", "pass", "Peter", "Loot",
								"ppd@gmail.com", "3800000000", 1))
				},
			},
			want: &models.User{
				ID:         1,
				Email:      "ppd@gmail.com",
				Username:   "peterpandam",
				Password:   "pass",
				FirstName:  "Peter",
				LastName:   "Loot",
				Phone:      "3800000000",
				UserStatus: 1,
			},
			wantErr: false,
		},
		{
			name:   "Failure 1",
			fields: fields{pool: pool},
			args: args{
				ctx:      context.Background(),
				username: "NOT_EXIST_USER",
				mockFn: func() {
					mock.ExpectQuery(`select (.+) from "user"`).
						WithArgs("NOT_EXIST_USER").
						WillReturnError(errors.New("can't get user from db"))
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Database{
				pool: tt.fields.pool,
			}

			tt.args.mockFn()

			got, err := d.GetUserByName(tt.args.ctx, tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserByName() got = %v, want %v", got, tt.want)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDatabase_Login(t *testing.T) {
	pool, mock := generateMockedDB()
	defer checkError(pool.Close)

	type fields struct {
		pool *sqlx.DB
	}
	type args struct {
		ctx      context.Context
		username string
		password string
		mockFn   func()
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:   "Success 1",
			fields: fields{pool: pool},
			args: args{
				ctx:      context.Background(),
				username: "admin",
				password: "password",
				mockFn: func() {
					mock.ExpectQuery(`select (.+) from user_info where user_name=.`).
						WithArgs("admin").
						WillReturnRows(sqlmock.NewRows([]string{"password", "allowed_methods"}).
							AddRow("$2a$10$fuyAMkr.wnE1O.maxUv1l.p6FTSpQMv2d1TbBbtHnIXXdg4NvZYJ6", "'{get}'"))
				},
			},
			want:    "'{get}'",
			wantErr: false,
		},
		{
			name:   "Failure 1",
			fields: fields{pool: pool},
			args: args{
				ctx:      context.Background(),
				username: "wrong",
				password: "password",
				mockFn: func() {
					mock.ExpectQuery(`select (.+) from user_info where user_name=.`).
						WithArgs("wrong").
						WillReturnError(errors.New("can't find username"))
				},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Database{
				pool: tt.fields.pool,
			}

			tt.args.mockFn()

			got, err := d.Login(tt.args.ctx, tt.args.username, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("Login() got = %v, want %v", got, tt.want)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDatabase_UpdateUser(t *testing.T) {
	pool, mock := generateMockedDB()
	defer checkError(pool.Close)

	type fields struct {
		pool *sqlx.DB
	}
	type args struct {
		ctx    context.Context
		user   *models.User
		mockFn func()
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Success 1",
			fields: fields{pool: pool},
			args: args{
				ctx: context.Background(),
				user: &models.User{
					ID:         1,
					Email:      "ppd@gmail.com",
					Username:   "peterpandam",
					Password:   "pass",
					FirstName:  "Peter",
					LastName:   "Loot",
					Phone:      "3800000000",
					UserStatus: 1,
				},
				mockFn: func() {
					mock.ExpectExec(`update "user" set (.+)`).
						WillReturnResult(sqlmock.NewResult(1, 1))
				},
			},
			wantErr: false,
		},
		{
			name:   "Failure 1",
			fields: fields{pool: pool},
			args: args{
				ctx: context.Background(),
				user: &models.User{
					ID:         1,
					Email:      "ppd@gmail.com",
					Username:   "peterpandam",
					Password:   "pass",
					FirstName:  "Peter",
					LastName:   "Loot",
					Phone:      "3800000000",
					UserStatus: 1,
				},
				mockFn: func() {
					mock.ExpectExec(`update "user" set (.+)`).
						WillReturnResult(sqlmock.NewResult(0, 0))
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Database{
				pool: tt.fields.pool,
			}

			tt.args.mockFn()

			if err := d.UpdateUser(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func Test_comparePassword(t *testing.T) {
	type args struct {
		encryptedPass string
		password      string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success 1",
			args: args{
				encryptedPass: "$2a$10$UwzV8tkdBpq7/WTzLocU8edTLQTT2k.4YCm15oBGiQ2gUqCU6R5Te",
				password:      "pass",
			},
			wantErr: false,
		},
		{
			name: "Failure 1",
			args: args{
				encryptedPass: "$2a$10$UwzV8tkdBpq7/WTzLocU8edTLQTT2k.4YCm15oBGiQ2gUqCU6R5Te",
				password:      "wrong",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := comparePassword(tt.args.encryptedPass, tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("comparePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_encryptPassword(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name:    "Success 1",
			args:    args{password: "pass"},
			want:    60,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encryptPassword(tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("encryptPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != tt.want {
				t.Errorf("encryptPassword() got = %v, want %v", got, tt.want)
			}
		})
	}
}
