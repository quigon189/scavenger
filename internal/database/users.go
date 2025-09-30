package database

import (
	"database/sql"
	"errors"
	"fmt"
	"scavenger/internal/models"
)

func (d *Database) CreateRole(roleName, roleDescription string) error {
	_, err := d.db.Exec(CreateRoleQuery, roleName, roleDescription)
	return err
}

func (d *Database) CreateUser(user *models.User) error {
	role := struct {
		ID          int
		Name        string
		Description string
	}{}

	err := d.db.QueryRow(GetRoleByName, user.RoleName).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("роль %s не найдена", user.RoleName)
	}

	user.RoleID = role.ID

	result, err := d.db.Exec(CreateUserQuery, user.Username, user.Name, user.PasswordHash, user.RoleID)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = int(id)
	return nil
}

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

func (d *Database) CreateStudent(student *models.User) error {
	group := &models.Group{}

	err := d.db.QueryRow(GetGroupByName, student.GroupName).Scan(
		&group.ID,
		&group.Name,
	)
	if err != nil {
		return err
	}

	student.GroupID = group.ID
	_, err = d.db.Exec(CreateStudentQuery, student.ID, student.GroupID)
	return err
}

func (d *Database) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	err := d.db.QueryRow(GetUserByUsernameQuery, username).Scan(
		&user.ID,
		&user.Username,
		&user.Name,
		&user.PasswordHash,
		&user.RoleName,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user %s: %v", username, err)
	}
	return user, nil
}
