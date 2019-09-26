package psql

import (
	"context"
	"strconv"
	"time"

	"github.com/IamStubborN/petstore/db/models"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func (d *Database) GetInventories(ctx context.Context) (map[string]int64, error) {
	rows, err := d.pool.QueryContext(ctx, qm[storeInventoriesQ])
	if err != nil {
		return nil, errors.Wrap(err, "can't get data from order_info")
	}
	defer checkError(rows.Close)

	result := make(map[string]int64)
	for rows.Next() {
		var orderStatus string
		var quantity int64

		if err = rows.Scan(&orderStatus, &quantity); err != nil {
			return nil, errors.Wrap(err, "can't scan data from order_info")
		}

		result[orderStatus] = quantity
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (d *Database) CreateOrder(ctx context.Context, order *models.Order) (*models.Order, error) {
	orderStatusID, err := d.getOrderStatusID(ctx, order.Status)
	if err != nil {
		return nil, errors.Wrap(err, "can't find order_status from db")
	}
	orderStatusName := order.Status
	order.Status = strconv.FormatInt(orderStatusID, 10)

	shipDate := time.Now().UTC()
	order.ShipDate = shipDate.Format(time.RFC3339)

	rows, err := d.pool.NamedQueryContext(ctx, qm[storeCreateQ], order)
	if err != nil {
		return nil, errors.Wrap(err, "can't insert into order table")
	}
	defer checkError(rows.Close)

	for rows.Next() {
		if err := rows.Scan(&order.ID); err != nil {
			return nil, errors.Wrap(err, "can't scan order_id from order table")
		}
	}

	order.Status = orderStatusName

	return order, nil
}

func (d *Database) FindOrderByID(ctx context.Context, orderID int64) (*models.Order, error) {
	var order models.Order
	err := d.pool.GetContext(ctx, &order, qm[storeFindByIDQ], orderID)
	if err != nil {
		return nil, errors.Wrap(err, "can't scan data from db")
	}

	return &order, nil
}

func (d *Database) DeleteOrderByID(ctx context.Context, orderID int64) error {
	res, err := d.pool.ExecContext(ctx, qm[storeDeleteByIDQ], orderID)
	if err != nil {
		return errors.Wrap(err, "can't find order")
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.Errorf("order № %d doesn't exist", orderID)
	}

	return nil
}

func (d *Database) getOrderStatusID(ctx context.Context, orderStatusName string) (int64, error) {
	var orderStatusID int64
	err := d.pool.GetContext(ctx, &orderStatusID, qm[storeGetStatusQ], orderStatusName)
	if err != nil {
		return 0, err
	}

	return orderStatusID, nil
}

func (d *Database) CreateInvoiceByDates(ctx context.Context, from, to string) ([]*models.InvoiceItem, error) {
	argQ := map[string]interface{}{
		"from": from,
		"to":   to,
	}

	rows, err := d.pool.NamedQueryContext(ctx, qm[storeCreateInvoiceByDatesQ], argQ)
	if err != nil {
		return nil, errors.Wrap(err, "can't get data for invoice")
	}
	defer checkError(rows.Close)

	var invoices []*models.InvoiceItem
	if err := sqlx.StructScan(rows, &invoices); err != nil {
		return nil, errors.Wrap(err, "can't scan data for invoice")
	}

	if len(invoices) == 0 {
		return nil, errors.New("no data with this dates")
	}

	return invoices, nil
}
