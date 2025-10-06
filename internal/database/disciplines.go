package database

import (
	"scavenger/internal/models"
	"strconv"
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
	return nil
}
