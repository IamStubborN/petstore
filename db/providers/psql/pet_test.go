package psql

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/IamStubborN/petstore/db/models"
	"github.com/lib/pq"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/jmoiron/sqlx"
)

func TestDatabase_AddPetToStore(t *testing.T) {
	pool, mock := generateMockedDB()
	defer checkError(pool.Close)

	type fields struct {
		pool *sqlx.DB
	}
	type args struct {
		ctx    context.Context
		pet    *models.Pet
		mockFn func()
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Pet
		wantErr bool
	}{
		{
			name:   "Success 1",
			fields: fields{pool: pool},
			args: args{
				ctx: context.Background(),
				pet: &models.Pet{
					ID: 1,
					Category: models.Category{
						ID:    1,
						Name:  "Cat",
						Price: 35.00,
					},
					Name:      "Soo",
					PhotoURLs: []string{"1", "2", "3"},
					Tags: []models.Tag{
						{ID: 1, Name: "small"},
						{ID: 4, Name: "best"},
						{ID: 5, Name: "cool"},
					},
					Status: "pending",
				},
				mockFn: func() {
					mock.ExpectBegin()
					mock.ExpectQuery(`select id from pet_status`).
						WithArgs("pending").
						WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
					mock.ExpectPrepare(`insert into pet (.+) values (.+) returning id`).ExpectQuery().
						WithArgs(1, "Soo", pq.Array([]string{"1", "2", "3"}), 1).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(45))
					mock.ExpectExec(`insert into pet_tag`).
						WithArgs(45, 1).WillReturnResult(sqlmock.NewResult(1, 1))
					mock.ExpectExec(`insert into pet_tag`).
						WithArgs(45, 4).WillReturnResult(sqlmock.NewResult(1, 1))
					mock.ExpectExec(`insert into pet_tag`).
						WithArgs(45, 5).WillReturnResult(sqlmock.NewResult(1, 1))
					mock.ExpectCommit()
				},
			},
			want: &models.Pet{
				ID: 45,
				Category: models.Category{
					ID:    1,
					Name:  "Cat",
					Price: 35.00,
				},
				Name:      "Soo",
				PhotoURLs: []string{"1", "2", "3"},
				Tags: []models.Tag{
					{ID: 1, Name: "small"},
					{ID: 4, Name: "best"},
					{ID: 5, Name: "cool"},
				},
				Status: "pending",
			},
			wantErr: false,
		},
		{
			name:   "Failure 1",
			fields: fields{pool: pool},
			args: args{
				ctx: context.Background(),
				pet: &models.Pet{
					ID: 1,
					Category: models.Category{
						ID:    1,
						Name:  "Cat",
						Price: 35.00,
					},
					Name:      "Soo",
					PhotoURLs: []string{"1", "2", "3"},
					Tags: []models.Tag{
						{ID: 1, Name: "small"},
						{ID: 4, Name: "best"},
						{ID: 5, Name: "cool"},
					},
					Status: "BAD_STATUS",
				},
				mockFn: func() {
					mock.ExpectBegin()
					mock.ExpectQuery(`select id from pet_status`).
						WithArgs("BAD_STATUS").
						WillReturnError(errors.New("can't find pet status"))
					mock.ExpectRollback()
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

			got, err := d.AddPetToStore(tt.args.ctx, tt.args.pet)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddPetToStore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddPetToStore() got = %v, want %v", got, tt.want)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDatabase_DeletePetByID(t *testing.T) {
	pool, mock := generateMockedDB()
	defer checkError(pool.Close)

	type fields struct {
		pool *sqlx.DB
	}
	type args struct {
		ctx    context.Context
		petID  int64
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
				ctx:   context.Background(),
				petID: 1,
				mockFn: func() {
					mock.ExpectExec(`delete from pet`).
						WithArgs(1).
						WillReturnResult(sqlmock.NewResult(0, 1))
				},
			},
			wantErr: false,
		},
		{
			name:   "Failure 1",
			fields: fields{pool: pool},
			args: args{
				ctx:   context.Background(),
				petID: 99,
				mockFn: func() {
					mock.ExpectExec(`delete from pet`).
						WithArgs(99).
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

			if err := d.DeletePetByID(tt.args.ctx, tt.args.petID); (err != nil) != tt.wantErr {
				t.Errorf("DeletePetByID() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDatabase_FindPetsByStatus(t *testing.T) {
	pool, mock := generateMockedDB()
	defer checkError(pool.Close)

	type fields struct {
		pool *sqlx.DB
	}
	type args struct {
		ctx    context.Context
		status []string
		mockFn func()
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.PetList
		wantErr bool
	}{
		{
			name:   "Success 1",
			fields: fields{pool: pool},
			args: args{
				ctx:    context.Background(),
				status: []string{"pending", "available"},
				mockFn: func() {
					mock.ExpectQuery(`select (.+) from pet_info`).
						WithArgs(pq.Array([]string{"pending", "available"})).
						WillReturnRows(sqlmock.NewRows(
							[]string{"id", "category_id", "category_name",
								"name", "photo_urls", "pet_status_name"}).
							AddRow(1, 1, "Cat", "Soo",
								pq.Array([]string{"1", "2", "3"}), "pending").
							AddRow(2, 2, "Dog", "Sylar",
								pq.Array([]string{"1", "2", "3"}), "available"))
					mock.ExpectQuery(`select (.+) from pet_tag`).
						WithArgs(1).
						WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
							AddRow(1, "small").AddRow(4, "best").AddRow(5, "cool"))
					mock.ExpectQuery(`select (.+) from pet_tag`).
						WithArgs(2).
						WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
							AddRow(1, "small").AddRow(2, "average").AddRow(3, "large"))
				},
			},
			want: &models.PetList{
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
					Status: "pending",
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
			},
			wantErr: false,
		},
		{
			name:   "Failure 1",
			fields: fields{pool: pool},
			args: args{
				ctx:    context.Background(),
				status: []string{"bad status"},
				mockFn: func() {
					mock.ExpectQuery(`select (.+) from pet_info`).
						WithArgs(pq.Array([]string{"bad status"})).
						WillReturnError(errors.New("can't get data from pet_info"))
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

			got, err := d.FindPetsByStatus(tt.args.ctx, tt.args.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindPetsByStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != nil {
				for idx := range *(got) {
					if !reflect.DeepEqual(*(*(got))[idx], *(*(tt.want))[idx]) {
						t.Errorf("FindPetsByStatus() got = %v, want %v", *(*(got))[idx], *(*(tt.want))[idx])
					}
				}
			}

			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDatabase_GetPetByID(t *testing.T) {
	pool, mock := generateMockedDB()
	defer checkError(pool.Close)

	type fields struct {
		pool *sqlx.DB
	}
	type args struct {
		ctx    context.Context
		id     int64
		mockFn func()
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Pet
		wantErr bool
	}{
		{
			name:   "Success 1",
			fields: fields{pool: pool},
			args: args{
				ctx: context.Background(),
				id:  1,
				mockFn: func() {
					mock.ExpectQuery(`select (.+) from pet_info`).
						WithArgs(1).
						WillReturnRows(sqlmock.NewRows([]string{
							"id", "category_id",
							"category_name", "name", "photo_urls",
							"pet_status_name",
						}).AddRow(1, 1, "Cat", "Soo",
							pq.Array([]string{"1", "2", "3"}), "pending"))
					mock.ExpectQuery(`select (.+) from pet_tag`).
						WithArgs(1).
						WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
							AddRow(1, "small").AddRow(4, "best").AddRow(5, "cool"))
				},
			},
			want: &models.Pet{
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
				Status: "pending",
			},
			wantErr: false,
		},
		{
			name:   "Failure 1",
			fields: fields{pool: pool},
			args: args{
				ctx: context.Background(),
				id:  99,
				mockFn: func() {
					mock.ExpectQuery(`select (.+) from pet_info`).
						WithArgs(99).
						WillReturnError(errors.New("can't find by id"))
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

			got, err := d.GetPetByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPetByID() got = %v, want %v", got, tt.want)
			}

			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDatabase_UpdatePetInStoreByBody(t *testing.T) {
	pool, mock := generateMockedDB()
	defer checkError(pool.Close)

	type fields struct {
		pool *sqlx.DB
	}
	type args struct {
		ctx    context.Context
		pet    *models.Pet
		mockFn func()
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Pet
		wantErr bool
	}{
		{
			name:   "Success 1",
			fields: fields{pool: pool},
			args: args{
				ctx: context.Background(),
				pet: &models.Pet{
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
					Status: "pending",
				},
				mockFn: func() {
					mock.ExpectBegin()
					mock.ExpectQuery(`select id from pet_status`).
						WithArgs("pending").
						WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
					mock.ExpectPrepare(`update pet set`).ExpectQuery().
						WithArgs(1, "Soo", pq.Array([]string{"1", "2", "3"}), 1, 1).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(45))
					mock.ExpectExec(`delete from pet_tag`).
						WithArgs(45).WillReturnResult(sqlmock.NewResult(45, 1))
					mock.ExpectExec(`insert into pet_tag`).
						WithArgs(45, 1).WillReturnResult(sqlmock.NewResult(45, 1))
					mock.ExpectExec(`insert into pet_tag`).
						WithArgs(45, 4).WillReturnResult(sqlmock.NewResult(45, 1))
					mock.ExpectExec(`insert into pet_tag`).
						WithArgs(45, 5).WillReturnResult(sqlmock.NewResult(45, 1))
					mock.ExpectCommit()
				},
			},
			want: &models.Pet{
				ID: 45,
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
				Status: "pending",
			},
			wantErr: false,
		},
		{
			name:   "Failure 1",
			fields: fields{pool: pool},
			args: args{
				ctx: context.Background(),
				pet: &models.Pet{
					ID: 1,
					Category: models.Category{
						ID:   99,
						Name: "WOOOOHOOOO",
					},
					Name:      "Soo",
					PhotoURLs: []string{"1", "2", "3"},
					Tags: []models.Tag{
						{ID: 1, Name: "small"},
						{ID: 4, Name: "best"},
						{ID: 5, Name: "cool"},
					},
					Status: "pending",
				},
				mockFn: func() {
					mock.ExpectBegin()
					mock.ExpectQuery(`select id from pet_status`).
						WithArgs("pending").
						WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
					mock.ExpectPrepare(`update pet set`).ExpectQuery().
						WithArgs(99, "Soo", pq.Array([]string{"1", "2", "3"}), 1, 1).
						WillReturnError(errors.New("can't exec add query"))
					mock.ExpectRollback()
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

			got, err := d.UpdatePetInStoreByBody(tt.args.ctx, tt.args.pet)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdatePetInStoreByBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != nil && !reflect.DeepEqual(*(got), *(tt.want)) {
				t.Errorf("UpdatePetInStoreByBody() got = %v, want %v", got, tt.want)
			}

			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDatabase_UpdatePetInStoreByForm(t *testing.T) {
	pool, mock := generateMockedDB()
	defer checkError(pool.Close)

	type fields struct {
		pool *sqlx.DB
	}
	type args struct {
		ctx    context.Context
		petID  int64
		name   string
		status string
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
				ctx:    context.Background(),
				petID:  1,
				name:   "Soo",
				status: "pending",
				mockFn: func() {
					mock.ExpectQuery(`select id from pet_status`).
						WithArgs("pending").
						WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
					mock.ExpectExec(`update pet set`).
						WithArgs("Soo", 1, 1).
						WillReturnResult(sqlmock.NewResult(0, 1))
				},
			},
			wantErr: false,
		},
		{
			name:   "Failure 1",
			fields: fields{pool: pool},
			args: args{
				ctx:    context.Background(),
				petID:  99,
				name:   "Soo",
				status: "pending",
				mockFn: func() {
					mock.ExpectQuery(`select id from pet_status`).
						WithArgs("pending").
						WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
					mock.ExpectExec(`update pet set`).
						WithArgs("Soo", 1, 99).
						WillReturnError(errors.New("can't update pet"))
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

			err := d.UpdatePetInStoreByForm(tt.args.ctx, tt.args.petID, tt.args.name, tt.args.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdatePetInStoreByForm() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDatabase_UpdatePetPhotosByID(t *testing.T) {
	pool, mock := generateMockedDB()
	defer checkError(pool.Close)

	type fields struct {
		pool *sqlx.DB
	}
	type args struct {
		ctx       context.Context
		petID     int64
		imagesURL []string
		mockFn    func()
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
				ctx:       context.Background(),
				petID:     1,
				imagesURL: []string{"1", "2", "3"},
				mockFn: func() {
					mock.ExpectExec(`update pet set`).
						WithArgs(pq.Array([]string{"1", "2", "3"}), 1).
						WillReturnResult(sqlmock.NewResult(1, 1))
				},
			},
			wantErr: false,
		},
		{
			name:   "Failure 1",
			fields: fields{pool: pool},
			args: args{
				ctx:       context.Background(),
				petID:     99,
				imagesURL: []string{"1", "2", "3"},
				mockFn: func() {
					mock.ExpectExec(`update pet set`).
						WithArgs(pq.Array([]string{"1", "2", "3"}), 99).
						WillReturnError(errors.New("pet № 99 doesn't exist"))
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

			if err := d.UpdatePetPhotosByID(tt.args.ctx, tt.args.petID, tt.args.imagesURL); (err != nil) != tt.wantErr {
				t.Errorf("UpdatePetPhotosByID() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDatabase_getTagsByPetID(t *testing.T) {
	pool, mock := generateMockedDB()
	defer checkError(pool.Close)

	type fields struct {
		pool *sqlx.DB
	}
	type args struct {
		ctx    context.Context
		petID  int64
		mockFn func()
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.Tag
		wantErr bool
	}{
		{
			name:   "Success 1",
			fields: fields{pool: pool},
			args: args{
				ctx:   context.Background(),
				petID: 1,
				mockFn: func() {
					mock.ExpectQuery(`select (.+) from pet_tag`).
						WithArgs(1).
						WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
							AddRow(1, "small").AddRow(4, "best").AddRow(5, "cool"))
				},
			},
			want: []models.Tag{
				{ID: 1, Name: "small"},
				{ID: 4, Name: "best"},
				{ID: 5, Name: "cool"},
			},
			wantErr: false,
		},
		{
			name:   "Failure 1",
			fields: fields{pool: pool},
			args: args{
				ctx:   context.Background(),
				petID: 99,
				mockFn: func() {
					mock.ExpectQuery(`select (.+) from pet_tag`).
						WithArgs(99).
						WillReturnError(errors.New("can't get pet tags"))
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

			got, err := d.getTagsByPetID(tt.args.ctx, tt.args.petID)
			if (err != nil) != tt.wantErr {
				t.Errorf("getTagsByPetID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getTagsByPetID() got = %v, want %v", got, tt.want)
			}

			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
