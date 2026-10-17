package database

import "scavenger/internal/models"

func (d *Database) AddStoredFile(file *models.StoredFile) error {
	var id int32
	err := d.db.QueryRow(
		CreateStoredFileQuery,
		file.Path,
		file.URL,
		file.Filename,
		file.Size,
	).Scan(&id)
	if err != nil {
		return err
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
