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

func (d *Database) GetRoles() ([]string, error) {
	var roles []string

	row, err := d.db.Query(GetAllRolesQuery)
	if err != nil {
		return roles, err
	}

	for row.Next() {
		var role string
		err := row.Scan(&role)
		if err == nil {
			roles = append(roles, role)
		}
	}

	return roles, nil
}

func (d *Database) GetAllStudents() ([]models.User, error) {
	students := []models.User{}
	row, err := d.db.Query(GetAllStudentsQuery)
	if err != nil {
		return students, err
	}

	for row.Next() {
		var stud models.User
		err := row.Scan(&stud.ID, &stud.Username, &stud.Name, &stud.GroupID, &stud.GroupName)
		if err == nil {
			stud.RoleName = string(models.StudentRole)
			students = append(students, stud)
		}
	}

	return students, nil
}

func (d *Database) GetStudentByGroupID(gID int) ([]models.User, error) {
	students := []models.User{}

	row, err := d.db.Query(GetStudentsByGroupIDQuery, gID)
	if err != nil {
		return students, err
	}

	for row.Next() {
		var stud models.User
		err := row.Scan(&stud.ID, &stud.Username, &stud.Name, &stud.GroupID, &stud.GroupName)
		if err == nil {
			stud.RoleName = string(models.StudentRole)
			students = append(students, stud)
		}
	}

	return students, nil
}

func (d *Database) GetStudentGroup(student *models.User) error {
	return d.db.QueryRow(GetStudentGroupQuery, student.ID).Scan(&student.GroupName)
}
