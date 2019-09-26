package psql

import (
	"context"
	"database/sql"

	"github.com/IamStubborN/petstore/db/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

func (d *Database) AddPetToStore(ctx context.Context, pet *models.Pet) (*models.Pet, error) {
	tx, err := d.pool.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "can't start transaction")
	}

	defer func() {
		if err != nil {
			checkError(tx.Rollback)
		}
		checkError(tx.Commit)
	}()

	var petStatusID int64
	err = tx.GetContext(ctx, &petStatusID, qm[petGetStatusIDByNameQ], pet.Status)
	if err != nil {
		return nil, errors.Wrap(err, "can't find pet_status")
	}

	stmt, err := tx.PrepareNamedContext(ctx, qm[petAddToStoreQ])
	if err != nil {
		return nil, err
	}
	defer checkError(stmt.Close)

	argQ := map[string]interface{}{
		"category_id":   pet.Category.ID,
		"name":          pet.Name,
		"photo_urls":    pq.Array(pet.PhotoURLs),
		"pet_status_id": petStatusID,
	}

	var petID int64
	err = stmt.QueryRowxContext(ctx, argQ).Scan(&petID)
	if err != nil {
		return nil, errors.Wrap(err, "can't exec add query")
	}

	pet.ID = petID

	for _, tag := range pet.Tags {
		argQ = map[string]interface{}{
			"pet_id": petID,
			"tag_id": tag.ID,
		}

		_, err = tx.NamedExecContext(ctx, qm[tagsInsertQ], argQ)
		if err != nil {
			return nil, errors.Wrap(err, "can't insert into pet_tag table")
		}
	}

	return pet, nil
}

func (d *Database) UpdatePetInStoreByBody(ctx context.Context, pet *models.Pet) (*models.Pet, error) {
	tx, err := d.pool.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "can't start transaction")
	}

	defer func() {
		if err != nil {
			checkError(tx.Rollback)
		}
		checkError(tx.Commit)
	}()

	var petStatusID int64
	err = tx.GetContext(ctx, &petStatusID, qm[petGetStatusIDByNameQ], pet.Status)
	if err != nil {
		return nil, errors.Wrap(err, "can't find pet_status")
	}

	stmt, err := tx.PrepareNamedContext(ctx, qm[petUpdateWithBodyQ])
	if err != nil {
		return nil, err
	}
	defer checkError(stmt.Close)

	argQ := map[string]interface{}{
		"id":            pet.ID,
		"category_id":   pet.Category.ID,
		"name":          pet.Name,
		"photo_urls":    pq.Array(pet.PhotoURLs),
		"pet_status_id": petStatusID,
	}

	var petID int64
	err = stmt.QueryRowxContext(ctx, argQ).Scan(&petID)
	if err != nil {
		return nil, errors.Wrap(err, "can't exec add query")
	}

	pet.ID = petID

	if _, err = tx.ExecContext(ctx, qm[tagsDeleteQ], petID); err != nil {
		return nil, errors.Wrap(err, "can't delete tags")
	}

	for _, tag := range pet.Tags {
		argQ = map[string]interface{}{
			"pet_id": petID,
			"tag_id": tag.ID,
		}

		_, err = tx.NamedExecContext(ctx, qm[tagsInsertQ], argQ)
		if err != nil {
			return nil, errors.Wrap(err, "can't insert into pet_tag table")
		}
	}

	return pet, nil
}

func (d *Database) FindPetsByStatus(ctx context.Context, status []string) (*models.PetList, error) {
	rows, err := d.pool.QueryxContext(ctx, qm[petFindByStatusQ], pq.Array(status))
	if err != nil {
		return nil, errors.Wrap(err, "can't get data from pet_info")
	}
	defer checkError(rows.Close)

	var pets models.PetList
	for rows.Next() {
		var pet models.Pet
		err = rows.Scan(
			&pet.ID,
			&pet.Category.ID,
			&pet.Category.Name,
			&pet.Name,
			pq.Array(&pet.PhotoURLs),
			&pet.Status,
		)
		if err != nil {
			return nil, err
		}

		pet.Tags, err = d.getTagsByPetID(ctx, pet.ID)
		if err != nil {
			return nil, errors.Wrap(err, "can't get tags by id")
		}

		pets = append(pets, &pet)
	}

	return &pets, nil
}

func (d *Database) GetPetByID(ctx context.Context, id int64) (*models.Pet, error) {
	var pet models.Pet
	row := d.pool.QueryRowxContext(ctx, qm[petGetByIDQ], id)

	err := row.Scan(
		&pet.ID, &pet.Category.ID,
		&pet.Category.Name, &pet.Name,
		pq.Array(&pet.PhotoURLs), &pet.Status)
	if err != nil {
		return nil, errors.Wrap(err, "can't find by id")
	}

	pet.Tags, err = d.getTagsByPetID(ctx, pet.ID)
	if err != nil {
		return nil, errors.Wrap(err, "can't get pet tags")
	}

	return &pet, nil
}

func (d *Database) UpdatePetInStoreByForm(ctx context.Context, petID int64, name, status string) error {
	var petStatusID int64
	err := d.pool.GetContext(ctx, &petStatusID, qm[petGetStatusIDByNameQ], status)
	if err != nil {
		return errors.Wrap(err, "can't find pet_status")
	}

	argQ := map[string]interface{}{
		"id":            petID,
		"name":          name,
		"pet_status_id": petStatusID,
	}

	res, err := d.pool.NamedExecContext(ctx, qm[petUpdateWithFieldsQ], argQ)
	if err != nil {
		return errors.Wrap(err, "can't update pet")
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.Errorf("pet № %d doesn't exist", petID)
	}

	return nil
}

func (d *Database) DeletePetByID(ctx context.Context, petID int64) error {
	res, err := d.pool.ExecContext(ctx, qm[petDeleteByIDQ], petID)
	if err != nil {
		return errors.Wrap(err, "can't delete from pet")
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.Errorf("pet № %d doesn't exist", petID)
	}

	return nil
}

func (d *Database) UpdatePetPhotosByID(ctx context.Context, petID int64, imagesURL []string) error {
	argQ := map[string]interface{}{
		"id":         petID,
		"photo_urls": pq.Array(imagesURL),
	}

	res, err := d.pool.NamedExecContext(ctx, qm[petUpdatePhotosByIDQ], argQ)
	if err != nil {
		return errors.Wrap(err, "can't update photo_urls")
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.Errorf("pet № %d doesn't exist", petID)
	}

	return nil
}

func (d *Database) getTagsByPetID(ctx context.Context, petID int64) ([]models.Tag, error) {
	rows, err := d.pool.QueryxContext(ctx, qm[tagsGetByPetIDQ], petID)
	if err != nil {
		return nil, err
	}
	defer checkError(rows.Close)

	var tags []models.Tag
	if err := sqlx.StructScan(rows, &tags); err != nil {
		return nil, err
	}

	return tags, nil
}
