package psql

import (
	"context"
	"errors"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/IamStubborN/petstore/db/models"
	"github.com/jmoiron/sqlx"
)

func generateMockedDB() (*sqlx.DB, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		log.Println(err)
	}
	pool := sqlx.NewDb(mockDB, "sqlmock")

	return pool, mock
}

func TestDatabase_CreateInvoiceByDates(t *testing.T) {
	pool, mock := generateMockedDB()
	defer checkError(pool.Close)

	type fields struct {
		pool *sqlx.DB
	}

	type args struct {
		ctx    context.Context
		from   string
		to     string
		mockFn func()
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*models.InvoiceItem
		wantErr bool
	}{
		{
			name:   "Success 1",
			fields: fields{pool: pool},
			args: args{
				ctx:  context.Background(),
				from: "2019-09-07",
				to:   "2019-09-08",
				mockFn: func() {
					mock.ExpectQuery(`select (.+) from invoice_info where ship_date between . and .`).
						WithArgs("2019-09-07", "2019-09-08").
						WillReturnRows(sqlmock.NewRows([]string{
							"id", "user_name", "pet", "category",
							"ship_date", "quantity", "price"}).
							AddRow(1, "Jack The Ripper", "John Snow",
								"Cat", "2019-08-07T15:35:04", 15, 35.00).
							AddRow(2, "Ginger", "Phil Heat",
								"Dog", "2019-09-08T15:35:04", 34, 49.99))
				},
			},
			want: []*models.InvoiceItem{
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
			},
		},
		{
			name:   "Failure 1",
			fields: fields{pool: pool},
			args: args{
				ctx:  context.Background(),
				from: "2019-09-12",
				to:   "2019-09-09",
				mockFn: func() {
					mock.ExpectQuery(`select (.+) from invoice_info where ship_date between . and .`).
						WithArgs("2019-09-12", "2019-09-09").
						WillReturnError(errors.New("can't get data for invoice"))
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

			got, err := d.CreateInvoiceByDates(tt.args.ctx, tt.args.from, tt.args.to)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateInvoiceByDates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for idx := range got {
				if !reflect.DeepEqual(got[idx].ID, tt.want[idx].ID) {
					t.Errorf("CreateInvoiceByDates() got = %v, want %v", got[idx], tt.want[idx])
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDatabase_CreateOrder(t *testing.T) {
	pool, mock := generateMockedDB()
	defer checkError(pool.Close)

	type fields struct {
		pool *sqlx.DB
	}
	type args struct {
		ctx    context.Context
		order  *models.Order
		mockFn func()
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Order
		wantErr bool
	}{
		{
			name:   "Success 1",
			fields: fields{pool: pool},
			args: args{
				ctx: context.Background(),
				order: &models.Order{
					ID:       1,
					PetID:    1,
					UserID:   1,
					Quantity: 12,
					ShipDate: "2019-09-05T15:35:12",
					Status:   "placed",
					Complete: false,
				},
				mockFn: func() {
					mock.ExpectQuery(`select (.+) from order_status where name=(.+)`).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
					mock.ExpectQuery(`insert into "order" (.+) values (.+)`).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(23))
				},
			},
			want: &models.Order{
				ID:       23,
				PetID:    1,
				UserID:   1,
				Quantity: 12,
				ShipDate: time.Now().UTC().Format(time.RFC3339),
				Status:   "placed",
				Complete: false,
			},
			wantErr: false,
		},
		{
			name:   "Failure 1",
			fields: fields{pool: pool},
			args: args{
				ctx: context.Background(),
				order: &models.Order{
					ID:       1,
					PetID:    99,
					UserID:   1,
					Quantity: 12,
					ShipDate: "2019-09-05T15:35:12",
					Status:   "BAD_STATUS",
					Complete: false,
				},
				mockFn: func() {
					mock.ExpectQuery(`select (.+) from order_status where name=(.+)`).
						WithArgs("BAD_STATUS").
						WillReturnError(errors.New("can't get order status id by name"))
				},
			},
			wantErr: true,
		},
		{
			name:   "Failure 2",
			fields: fields{pool: pool},
			args: args{
				ctx: context.Background(),
				order: &models.Order{
					ID:       1,
					PetID:    99,
					UserID:   1,
					Quantity: 12,
					ShipDate: "2019-09-05T15:35:12",
					Status:   "placed",
					Complete: false,
				},
				mockFn: func() {
					mock.ExpectQuery(`select (.+) from order_status where name=(.+)`).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
					mock.ExpectQuery(`insert into "order" (.+) values (.+)`).
						WillReturnError(errors.New("can't get order, invalid pet_id"))
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

			got, err := d.CreateOrder(tt.args.ctx, tt.args.order)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateOrder() error = %v, wantErr %v", err, tt.wantErr)
			}

			if got != nil && !reflect.DeepEqual(*(got), *(tt.want)) {
				t.Errorf("CreateOrder() got = %v, want %v", got, tt.want)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDatabase_DeleteOrderByID(t *testing.T) {
	pool, mock := generateMockedDB()
	defer checkError(pool.Close)

	type fields struct {
		pool *sqlx.DB
	}
	type args struct {
		ctx     context.Context
		orderID int64
		mockFn  func()
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
				ctx:     context.Background(),
				orderID: 1,
				mockFn: func() {
					mock.ExpectExec(`delete from "order" where id=.`).
						WillReturnResult(sqlmock.NewResult(1, 1))

				},
			},
			wantErr: false,
		},
		{
			name:   "Failure 1",
			fields: fields{pool: pool},
			args: args{
				ctx:     context.Background(),
				orderID: 99,
				mockFn: func() {
					mock.ExpectExec(`delete from "order" where id=.`).
						WithArgs(99).
						WillReturnError(errors.New("can't find order by id"))
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

			if err := d.DeleteOrderByID(tt.args.ctx, tt.args.orderID); (err != nil) != tt.wantErr {
				t.Errorf("DeleteOrderByID() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDatabase_FindOrderByID(t *testing.T) {
	pool, mock := generateMockedDB()
	defer checkError(pool.Close)

	type fields struct {
		pool *sqlx.DB
	}
	type args struct {
		ctx     context.Context
		orderID int64
		mockFn  func()
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Order
		wantErr bool
	}{
		{
			name:   "Success 1",
			fields: fields{pool: pool},
			args: args{
				ctx:     context.Background(),
				orderID: 1,
				mockFn: func() {
					mock.ExpectQuery(`select (.+) from order_info where id=.`).
						WillReturnRows(sqlmock.NewRows([]string{
							"id", "user_id", "pet_id", "quantity",
							"ship_date", "order_status", "complete"}).
							AddRow(1, 1, 1, 12, "2019-09-05T15:35:12", "placed", false))
				},
			},
			want: &models.Order{
				ID:       1,
				PetID:    1,
				UserID:   1,
				Quantity: 12,
				ShipDate: "2019-09-05T15:35:12",
				Status:   "placed",
				Complete: false,
			},
			wantErr: false,
		},
		{
			name:   "Failure 1",
			fields: fields{pool: pool},
			args: args{
				ctx:     context.Background(),
				orderID: 99,
				mockFn: func() {
					mock.ExpectQuery(`select (.+) from order_info where id=.`).
						WithArgs(99).
						WillReturnError(errors.New("can't find order"))
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

			got, err := d.FindOrderByID(tt.args.ctx, tt.args.orderID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindOrderByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindOrderByID() got = %v, want %v", got, tt.want)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDatabase_GetInventories(t *testing.T) {
	pool, mock := generateMockedDB()
	defer checkError(pool.Close)

	type fields struct {
		pool *sqlx.DB
	}
	type args struct {
		ctx    context.Context
		mockFn func()
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]int64
		wantErr bool
	}{
		{
			name:   "Success 1",
			fields: fields{pool: pool},
			args: args{
				ctx: context.Background(),
				mockFn: func() {
					mock.ExpectQuery("select (.+) from order_info").
						WillReturnRows(sqlmock.NewRows([]string{"pet_status", "sum"}).
							AddRow("sold", 516).
							AddRow("pending", 663).
							AddRow("available", 669))
				},
			},
			want: map[string]int64{
				"available": 669,
				"pending":   663,
				"sold":      516,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Database{
				pool: tt.fields.pool,
			}

			tt.args.mockFn()

			got, err := d.GetInventories(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInventories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetInventories() got = %v, want %v", got, tt.want)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDatabase_getOrderStatusID(t *testing.T) {
	pool, mock := generateMockedDB()
	defer checkError(pool.Close)

	type fields struct {
		pool *sqlx.DB
	}
	type args struct {
		ctx             context.Context
		orderStatusName string
		mockFn          func()
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name:   "Success 1",
			fields: fields{pool: pool},
			args: args{
				ctx:             context.Background(),
				orderStatusName: "placed",
				mockFn: func() {
					mock.ExpectQuery(
						"select id from order_status where name=.").
						WithArgs("placed").
						WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				},
			},
			want:    1,
			wantErr: false,
		},
		{
			name:   "Failure 1",
			fields: fields{pool: pool},
			args: args{
				ctx:             context.Background(),
				orderStatusName: "BAD_STATUS_NAME",
				mockFn: func() {
					mock.ExpectQuery(
						"select id from order_status where name=.").
						WithArgs("BAD_STATUS_NAME").
						WillReturnError(errors.New("can't find order_status from db"))
				},
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Database{
				pool: tt.fields.pool,
			}

			tt.args.mockFn()

			got, err := d.getOrderStatusID(tt.args.ctx, tt.args.orderStatusName)
			if (err != nil) != tt.wantErr {
				t.Errorf("getOrderStatusID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("getOrderStatusID() got = %v, want %v", got, tt.want)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
