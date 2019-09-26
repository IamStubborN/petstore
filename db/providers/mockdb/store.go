package mockdb

import (
	"context"
	"errors"

	"github.com/IamStubborN/petstore/db/models"
)

func (d *Database) GetInventories(ctx context.Context) (map[string]int64, error) {
	return map[string]int64{
		"available": 353,
		"pending":   853,
		"sold":      351,
	}, nil
}

func (d *Database) CreateOrder(ctx context.Context, order *models.Order) (*models.Order, error) {
	if order.ID != 1 {
		return nil, errors.New("bad order input")
	}

	return order, nil
}

func (d *Database) FindOrderByID(ctx context.Context, orderID int64) (*models.Order, error) {
	if orderID != 1 {
		return nil, errors.New("bad orderID")
	}

	return &models.Order{
		ID:       1,
		PetID:    1,
		UserID:   1,
		Quantity: 12,
		ShipDate: "2019-09-05T15:35:12",
		Status:   "placed",
		Complete: false,
	}, nil
}

func (d *Database) DeleteOrderByID(ctx context.Context, orderID int64) error {
	var err error
	if orderID != 1 {
		err = errors.New("bad orderID")
	}

	return err
}

func (d *Database) CreateInvoiceByDates(ctx context.Context, from, to string) ([]*models.InvoiceItem, error) {
	if from != "2019-09-07" || to != "2019-09-08" {
		return nil, errors.New("bad dates inputs")
	}

	return []*models.InvoiceItem{
		{
			ID:       1,
			User:     "Jack The Ripper",
			Pet:      "John Snow",
			Category: "Cat",
			ShipDate: "2019-09-07T15:35:04",
			Quantity: 15,
			Price:    35.00,
		},
		{
			ID:       2,
			User:     "Ginger",
			Pet:      "Phil Heat",
			Category: "Dog",
			ShipDate: "2019-09-08T15:35:04",
			Quantity: 34,
			Price:    49.99,
		},
	}, nil
}
