package database

import "scavenger/internal/models"

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
			discs =append(discs, disc)
		}
	}

	return discs, nil
}

func (d *Database) GetDisciplinesByGroupID(groupID int) ([]models.Discipline, error) {
	discs := []models.Discipline{}
	row, err := d.db.Query(GetAllDisciplinesByGroupIDQuery, groupID)
	if err != nil {
		return discs, err
	}

	for row.Next() {
		var disc models.Discipline
		err := row.Scan(&disc.ID, &disc.Name, &disc.GroupID)
		if err == nil {
			discs =append(discs, disc)
		}
	}

	return discs, nil

}
