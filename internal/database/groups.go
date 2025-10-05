package database

import (
	"scavenger/internal/models"
)

func (d *Database) CreateGroup(group *models.Group) error {
	result, err := d.db.Exec(CreateGroupQuery, 0, group.Name)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	group.ID = int(id)
	return nil
}

func (d *Database) GetAllGroups() ([]models.Group, error) {
	groups := []models.Group{}
	row, err := d.db.Query(GetAllGroupsQuery)
	if err != nil {
		return groups, err
	}

	for row.Next() {
		var group models.Group
		err := row.Scan(&group.ID, &group.Name)
		if err == nil {
			groups = append(groups, group)
		}
	}

	return groups, nil
}

func (d *Database) GetAllGroupsWithStudents() ([]models.Group, error) {
	groups, err := d.GetAllGroups()
	if err != nil {
		return nil, err
	}
	for i := range groups {
		students, err := d.GetStudentByGroupID(groups[i].ID)
		if err == nil {
			groups[i].Students = append(groups[i].Students, students...)
		}
	}

	return groups, nil
}

func (d *Database) GetAllGroupsWithDisciplines() ([]models.Group, error) {
	groups, err := d.GetAllGroups()
	if err != nil {
		return nil, err
	}
	for i := range groups {
		discs, err := d.GetDisciplinesByGroupID(groups[i].ID)
		if err == nil {
			groups[i].Disciplines = append(groups[i].Disciplines, discs...)
		}
	}

	return groups, nil
}

func (d *Database) GetGroupByID(id int) (*models.Group, error) {
	group := &models.Group{}
	err := d.db.QueryRow(GetGroupByIDQuery, id).Scan(&group.ID, &group.Name)
	return group, err
}

func (d *Database) DeleteGroupByID(id int) error {
	_, err := d.db.Exec(DeleteGroupByIDQuery, id)
	return err
}
