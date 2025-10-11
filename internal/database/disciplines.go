package database

import (
	"log"
	"scavenger/internal/models"
	"strconv"
	"time"
)

func (d *Database) CreateDiscipline(disc *models.Discipline) error {
	result, err := d.db.Exec(CreateDisciplineQuery, disc.Name, disc.GroupID)
	if err != nil {
		return err
	}

	dID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	disc.ID = int(dID)
	return nil
}

func (d *Database) GetDisciplines() ([]models.Discipline, error) {
	discs := []models.Discipline{}
	row, err := d.db.Query(GetAllDisciplinesQuery)
	if err != nil {
		return discs, err
	}

	for row.Next() {
		var disc models.Discipline
		err := row.Scan(&disc.ID, &disc.Name, &disc.GroupID)
		if err == nil {
			discs = append(discs, disc)
		}
	}

	return discs, nil
}

func (d *Database) GetDisciplinesByGroupID(groupID int) ([]models.Discipline, error) {
	discs := []models.Discipline{}
	row, err := d.db.Query(GetAllDisciplinesByGroupIDQuery, &groupID)
	if err != nil {
		return discs, err
	}

	for row.Next() {
		var disc models.Discipline
		err := row.Scan(&disc.ID, &disc.Name, &disc.GroupID)
		if err == nil {
			discs = append(discs, disc)
		}
	}

	return discs, nil
}

func (d *Database) GetDisciplinesWithGroup() ([]models.Discipline, error) {
	discs, err := d.GetDisciplines()
	if err != nil {
		return []models.Discipline{}, err
	}

	for i := range discs {
		if discs[i].GroupID == nil {
			discs[i].Group.Name = "Без группы"
		} else {
			group, _ := d.GetGroupByID(*discs[i].GroupID)
			discs[i].Group = *group
		}
	}

	return discs, nil
}

func (d *Database) GetDisciplinesWithoutGroup() ([]models.Discipline, error) {
	discs := []models.Discipline{}

	row, err := d.db.Query(GetDisciplinesWithoutGroupQuery)
	if err != nil {
		return discs, err
	}

	for row.Next() {
		var disc models.Discipline
		err := row.Scan(&disc.ID, &disc.Name, &disc.GroupID)
		if err == nil {
			discs = append(discs, disc)
		}
	}

	return discs, nil
}

func (d *Database) GetDisciplineByID(id int) (*models.Discipline, error) {
	disc := &models.Discipline{}

	err := d.db.QueryRow(GetDisciplineByIDQuery, id).Scan(&disc.ID, &disc.Name, &disc.GroupID)
	if err != nil {
		return nil, err
	}

	return disc, nil
}

func (d *Database) GetDisciplineWithLabs(id int) (*models.Discipline, error) {
	disc, err := d.GetDisciplineByID(id)
	if err != nil {
		return &models.Discipline{}, err
	}

	labs, err := d.GetDisciplineLabs(id)
	if err != nil {
		return &models.Discipline{}, err
	}

	disc.Labs = append(disc.Labs, labs...)

	return disc, nil
}

func (d *Database) UpdateDiscipline(disc *models.Discipline) error {
	_, err := d.db.Exec(UpdateDisciplineQuery, disc.Name, disc.GroupID, disc.ID)
	return err
}

func (d *Database) DeleteDiscipline(disc *models.Discipline) error {
	_, err := d.db.Exec(DeleteDisciplineQuery, disc.ID)
	return err
}

func (d *Database) AddDisciplineLab(lab *models.Lab) error {
	result, err := d.db.Exec(
		CreateDisciplineLabQuery,
		lab.Name,
		lab.Description,
		lab.MDFileID,
		lab.Deadline.Unix(),
		lab.DisciplineID,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	lab.ID = strconv.Itoa(int(id))

	err = d.AddLabFiles(int(id), lab.StoredFiles)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) AddLabFiles(labID int, files []models.StoredFile) error {
	var err error = nil

	for _, file := range files {
		_, f_err := d.db.Exec(CreateLabFileQuery, file.ID, labID)
		if f_err != nil {
			err = f_err
		}
	}
	return err
}

func (d *Database) RemoveLabFiles(labID int, files []models.StoredFile) error {
	var err error = nil

	for _, file := range files {
		_, f_err := d.db.Exec(DeleteLabFileQuery, file.ID, labID)
		if f_err != nil {
			err = f_err
		}
	}

	return err
}

func (d *Database) GetDisciplineLabs(id int) ([]models.Lab, error) {
	labs := []models.Lab{}

	row, err := d.db.Query(GetDisciplineLabsQuery, id)
	if err != nil {
		return labs, err
	}

	for row.Next() {
		lab := models.Lab{}
		var deadline int64
		err = row.Scan(
			&lab.ID,
			&lab.Name,
			&lab.Description,
			&lab.MDFileID,
			&deadline,
		)
		if err != nil {
			continue
		}

		lab.Deadline = time.Unix(deadline, 0)

		mdFile, _ := d.GetStoredFile(lab.MDFileID)
		lab.MDFile = *mdFile

		files, err := d.GetLabFiles(lab.ID)
		if err != nil {
			log.Printf("err to get lab files: %v", err)
			labs = append(labs, lab)
			continue
		}

		lab.StoredFiles = append(lab.StoredFiles, files...)
		labs = append(labs, lab)
	}

	return labs, nil
}

func (d *Database) GetLabFiles(id string) ([]models.StoredFile, error) {
	var files []models.StoredFile

	row, err := d.db.Query(GetLabFilesQuery, id)
	if err != nil {
		return files, err
	}
	for row.Next() {
		var file models.StoredFile
		err := row.Scan(
			&file.ID,
			&file.Path,
			&file.URL,
			&file.Filename,
			&file.Size,
		)
		if err != nil {
			return []models.StoredFile{}, err
		}
		files = append(files, file)
	}

	return files, nil
}

func (d *Database) GetLabByID(id int) (*models.Lab, error) {
	lab := &models.Lab{}
	var deadline int64

	err := d.db.QueryRow(GetLabByIDQuery, id).Scan(
		&lab.ID,
		&lab.Name,
		&lab.Description,
		&lab.MDFileID,
		&deadline,
		&lab.DisciplineID,
	)
	if err != nil {
		return &models.Lab{}, err
	}

	lab.Deadline = time.Unix(deadline, 0)

	mdfile, _ := d.GetStoredFile(lab.MDFileID)
	lab.MDFile = *mdfile

	files, err := d.GetLabFiles(lab.ID)
	if err != nil {
		return &models.Lab{}, err
	}

	lab.StoredFiles = append(lab.StoredFiles, files...)

	return lab, nil
}

func (d *Database) UpdateLab(lab *models.Lab) error {
	_, err := d.db.Exec(
		UpdateLabQuery,
		lab.Name,
		lab.Description,
		lab.MDFileID,
		lab.Deadline.Unix(),
		lab.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) DeleteLab(labID int) error {
	_, err := d.db.Exec(DeleteLabQuery, labID)
	return err
}
