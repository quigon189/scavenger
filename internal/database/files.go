package database

import "scavenger/internal/models"

func (d *Database) AddStoredFile(file *models.StoredFile) error {
	result, err := d.db.Exec(
		CreateStoredFileQuery,
		file.Path,
		file.URL,
		file.Filename,
		file.Size,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil
	}

	file.ID = int(id)
	return nil
}

func (d *Database) GetStoredFileByURL(url string) (*models.StoredFile, error) {
	file := &models.StoredFile{}

	err := d.db.QueryRow(GetStoredFileByURLQuery, url).Scan(
		&file.ID,
		&file.Path,
		&file.URL,
		&file.Filename,
		&file.Size,
	)
	if err != nil {
		return &models.StoredFile{}, err
	}

	return file, nil

}

func (d *Database) GetStoredFile(id int) (*models.StoredFile, error) {
	file := &models.StoredFile{}

	err := d.db.QueryRow(GetStoredFileQuery, id).Scan(
		&file.ID,
		&file.Path,
		&file.URL,
		&file.Filename,
		&file.Size,
	)
	if err != nil {
		return &models.StoredFile{}, err
	}

	return file, nil
}
