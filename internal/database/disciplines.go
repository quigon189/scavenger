package database

import (
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
		lab.MDPath,
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

	for _, filePath := range lab.FilesPath {
		d.db.Exec(CreateLabFilesQuery, lab.ID, filePath)
	}
	return nil
}

func (d *Database) AddLabFiles(labID int, filePaths []string) error {
	var err error = nil

	for _, filePath := range filePaths {
		_, f_err := d.db.Exec(CreateLabFilesQuery, labID, filePath)
		if f_err != nil {
			err = f_err
		}
	}
	return err
}

func (d *Database) RemoveLabFiles(labID int, filePaths []string) error {
	var err error = nil

	for _, filePath := range filePaths {
		_, f_err := d.db.Exec(DeleteLabFilesQuery, labID, filePath)
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
			&lab.MDPath,
			&deadline,
		)
		if err != nil {
			continue
		}

		lab.Deadline = time.Unix(deadline, 0)

		filesRow, err := d.db.Query(GetLabFilesQuery, lab.ID)
		if err == nil {
			for filesRow.Next() {
				var filePath string
				err := filesRow.Scan(&filePath)
				if err != nil {
					continue
				}
				lab.FilesPath = append(lab.FilesPath, filePath)
			}
		}

		labs = append(labs, lab)
	}

	return labs, nil
}

func (d *Database) GetLabByID(id int) (*models.Lab, error) {
	lab := &models.Lab{}
	var deadline int64

	err := d.db.QueryRow(GetLabByIDQuery, id).Scan(
		&lab.ID,
		&lab.Name,
		&lab.Description,
		&lab.MDPath,
		&deadline,
		&lab.DisciplineID,
	)
	if err != nil {
		return &models.Lab{}, err
	}

	row, err := d.db.Query(GetLabFilesQuery, lab.ID)
	if err != nil {
		return &models.Lab{}, err
	}

	for row.Next() {
		var filePath string
		err := row.Scan(&filePath)
		if err == nil {
			lab.FilesPath = append(lab.FilesPath, filePath)
		}
	}

	return lab, nil
}

func (d *Database) UpdateLab(lab *models.Lab) error {
	_, err := d.db.Exec(
		UpdateLabQuery,
		lab.Name,
		lab.Description,
		lab.MDPath,
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
