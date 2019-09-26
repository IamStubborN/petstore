package db

import (
	"context"

	"github.com/IamStubborN/petstore/config"
	"github.com/IamStubborN/petstore/db/models"
	"github.com/IamStubborN/petstore/db/providers/mockdb"
	"github.com/IamStubborN/petstore/db/providers/psql"
	"github.com/pkg/errors"

	"go.uber.org/zap"
)

var storage Storage

type Storage interface {
	StoreDI
	PetDI
	UserDI
	Close() error
}

type StoreDI interface {
	GetInventories(ctx context.Context) (map[string]int64, error)
	CreateOrder(ctx context.Context, order *models.Order) (*models.Order, error)
	FindOrderByID(ctx context.Context, orderID int64) (*models.Order, error)
	DeleteOrderByID(ctx context.Context, orderID int64) error
	CreateInvoiceByDates(ctx context.Context, from, to string) ([]*models.InvoiceItem, error)
}

type PetDI interface {
	AddPetToStore(ctx context.Context, pet *models.Pet) (*models.Pet, error)
	UpdatePetInStoreByBody(ctx context.Context, pet *models.Pet) (*models.Pet, error)
	UpdatePetInStoreByForm(ctx context.Context, id int64, name, status string) error
	FindPetsByStatus(ctx context.Context, status []string) (*models.PetList, error)
	GetPetByID(ctx context.Context, id int64) (*models.Pet, error)
	DeletePetByID(ctx context.Context, id int64) error
	UpdatePetPhotosByID(ctx context.Context, id int64, imagesURL []string) error
}

type UserDI interface {
	CreateUser(ctx context.Context, user *models.User) error
	CreateUsersFromList(ctx context.Context, list *models.UserList) error
	Login(ctx context.Context, username, password string) (string, error)
	GetUserByName(ctx context.Context, username string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, username string) error
}

func InitDatabase(cfg *config.Config) {
	switch cfg.DB.Provider {
	case "postgres":
		storage = psql.Database{}.InitDatabase(cfg.DB)
	case "mockdb":
		storage = mockdb.Database{}.InitDatabase(cfg.DB)
	default:
		zap.L().Fatal("wrong database provider")
	}
}

func GetPetDI() PetDI {
	return storage
}

func GetStoreDI() StoreDI {
	return storage
}

func GetUserDI() UserDI {
	return storage
}

func Close() error {
	if err := storage.Close(); err != nil {
		return errors.Wrap(err, "can't close database")
	}

	return nil
}
