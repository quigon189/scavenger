package database

import (
	"database/sql"
	"errors"
	"fmt"
	"scavenger/internal/models"
)

func (d *Database) CreateUserWithRole(user *models.User, roleName string) error {
	role := struct{
		ID int
		Name string
		Description string
	}{}

	err := d.db.QueryRow(GetRoleByName, roleName).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
	)	
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("роль %s не найдена", roleName)
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
