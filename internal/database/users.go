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

func (d *Database) CreateUserWithRole(user *models.User) error {
	role := struct{
		ID int
		Name string
		Description string
	}{}

	err := d.db.QueryRow(GetRoleByName, user.Role).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
	)	
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("роль %s не найдена", user.Role)
	}
	
	result, err := d.db.Exec(CreateUserWithRoleQuery, user.Username, user.Name, user.PasswordHash, role.ID)
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

func (d *Database) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	err := d.db.QueryRow(GetUserByUsernameQuery, username).Scan(
		&user.ID,
		&user.Username,
		&user.Name,
		&user.PasswordHash,
		&user.RoleName,
		&user.GroupName,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("пользователь с username: %s не найден", username)
	}
	return user, nil
}
