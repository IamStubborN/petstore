package mockdb

import (
	"context"
	"errors"
	"strings"

	"github.com/IamStubborN/petstore/db/models"
)

const testPetName = "TestPet"

func (d *Database) AddPetToStore(ctx context.Context, pet *models.Pet) (*models.Pet, error) {
	var err error
	if pet.Name != testPetName {
		err = errors.New("invalid pet")
	}

	return pet, err
}

func (d *Database) UpdatePetInStoreByBody(ctx context.Context, pet *models.Pet) (*models.Pet, error) {
	if pet.ID != 1 {
		return nil, errors.New("bad input")
	}

	return pet, nil
}

func (d *Database) UpdatePetInStoreByForm(ctx context.Context, id int64, name string, status string) error {
	var err error
	if id != 1 || name != testPetName || status != "pending" {
		err = errors.New("bad input")
	}

	return err
}

func (d *Database) FindPetsByStatus(ctx context.Context, status []string) (*models.PetList, error) {
	if strings.Contains(strings.Join(status, " "), "available") {
		return &models.PetList{
			{
				ID: 1,
				Category: models.Category{
					ID:   1,
					Name: "Cat",
				},
				Name:      "Soo",
				PhotoURLs: []string{"1", "2", "3"},
				Tags: []models.Tag{
					{ID: 1, Name: "small"},
					{ID: 4, Name: "best"},
					{ID: 5, Name: "cool"},
				},
				Status: "available",
			},
			{
				ID: 2,
				Category: models.Category{
					ID:   2,
					Name: "Dog",
				},
				Name:      "Sylar",
				PhotoURLs: []string{"1", "2", "3"},
				Tags: []models.Tag{
					{ID: 1, Name: "small"},
					{ID: 2, Name: "average"},
					{ID: 3, Name: "large"},
				},
				Status: "available",
			},
		}, nil
	}
	return nil, errors.New("can't find status")
}

func (d *Database) GetPetByID(ctx context.Context, id int64) (*models.Pet, error) {
	if id != 1 {
		return nil, errors.New("invalid pet id ")
	}

	return &models.Pet{ID: 1,
		Category: models.Category{
			ID:   1,
			Name: "Cat",
		},
		Name:      "Soo",
		PhotoURLs: []string{"1", "2", "3"},
		Tags: []models.Tag{
			{ID: 1, Name: "small"},
			{ID: 4, Name: "best"},
			{ID: 5, Name: "cool"},
		},
		Status: "available",
	}, nil
}

func (d *Database) DeletePetByID(ctx context.Context, id int64) error {
	var err error
	if id != 1 {
		err = errors.New("invalid pet id")
	}

	return err
}

func (d *Database) UpdatePetPhotosByID(ctx context.Context, id int64, imagesURL []string) error {
	var err error
	if id != 1 && imagesURL != nil {
		err = errors.New("invalid pet id or imagesURLs")
	}

	return err
}
