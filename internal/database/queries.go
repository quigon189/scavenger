package database

const (
	// Roles
	CreateRoleQuery = `INSERT INTO roles (name, description) VALUES (?, ?)`
	GetRoleByName = `SELECT id, name, description FROM roles WHERE name = ?`
	GetAllRolesQuery = `SELECT name FROM roles`

	// Groups
	GetAllGroupsQuery = `SELECT id, name FROM groups`
	GetGroupByIDQuery = `SELECT id, name FROM groups WHERE id = ?`
	GetGroupByName = `SELECT id, name FROM groups WHERE name = ?`
	GetStudentGroupQuery = `
		SELECT g.name, g.id
		FROM students s
		LEFT JOIN groups g ON s.group_id = g.id
		WHERE s.user_id = ?
	`
	CreateGroupQuery = `INSERT INTO groups (number, name) VALUES (?, ?)`
	DeleteGroupByIDQuery = `DELETE FROM groups WHERE id = ?`
	UpdateGroupQuery = `UPDATE groups SET name = ? WHERE id = ?`

	// Users
	CreateUserQuery = `INSERT INTO users (username, name, password_hash, role_id) VALUES (?, ?, ?, ?)`
	GetUserByUsernameQuery = `
		SELECT u.id, u.username, u.name, u.password_hash, r.name AS role_name
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.username = ?
	`
	UpdateUserQuery = `UPDATE users SET username = ?, name = ?, password_hash = ? WHERE id = ?` 
	DeleteUserQuery =`DELETE FROM users WHERE id = ?` 

	// Students
	CreateStudentQuery = `INSERT INTO students (user_id, group_id) VALUES (?, ?)`
	GetStudentsByGroupIDQuery = `
		SELECT u.id, u.username, u.name, g.id, g.name
		FROM students s
		LEFT JOIN users u ON s.user_id = u.id
		LEFT JOIN groups g ON s.group_id = g.id
		WHERE s.group_id = ?
	`
	GetAllStudentsQuery = `
		SELECT u.id, u.username, u.name, g.id, g.name
		FROM students s
		LEFT JOIN users u ON s.user_id = u.id
		LEFT JOIN groups g ON s.group_id = g.id
	`
	GetStudentQuery = `
		SELECT u.id, u.username, u.name, g.name
		FROM students s
		LEFT JOIN users u ON s.user_id = u.id
		LEFT JOIN groups g ON s.group_id = g.id
		WHERE s.user_id = ?
	`
	GetStudentByUsernameQuery = `
		SELECT u.id, u.username, u.name, g.id, g.name
		FROM students s
		LEFT JOIN users u ON s.user_id = u.id
		LEFT JOIN groups g ON s.group_id = g.id
		WHERE u.username = ?
	`
	UpdateStudentQuery = `UPDATE students SET group_id = ? WHERE user_id = ?`
	DeleteStudentQuery = `DELETE FROM students WHERE user_id = ?`

	// Disciplines
	CreateDisciplineQuery = `INSERT INTO disciplines (name, group_id) VALUES (?, ?)`
	GetAllDisciplinesQuery = `SELECT id, name, group_id FROM disciplines`
	GetAllDisciplinesByGroupIDQuery = `SELECT id, name, group_id FROM disciplines WHERE group_id = ?`
	GetDisciplinesWithoutGroupQuery = `SELECT id, name, group_id FROM disciplines WHERE group_id IS NULL`
	GetDisciplineByIDQuery = `SELECT id, name, group_id FROM disciplines WHERE id = ?`
	UpdateDisciplineQuery = `UPDATE disciplines SET name = ?, group_id = ? WHERE id = ?`
	DeleteDisciplineQuery = `DELETE FROM disciplines WHERE id = ?`

	// Lab
	CreateDisciplineLabQuery = `INSERT INTO labs (name, description, md_path, deadline, discipline_id) VALUES (?, ?, ?, ?, ?)`
	CreateLabFilesQuery = `INSERT INTO lab_files (lab_id, file_path) VALUES (?, ?)`
	GetDisciplineLabsQuery = `SELECT id, name, description, md_path, deadline FROM labs WHERE discipline_id = ?`
	GetLabFilesQuery = `SELECT file_path FROM lab_files WHERE lab_id = ?`
	GetLabByIDQuery = `SELECT id, name, description, md_path, deadline, discipline_id FROM labs WHERE id = ?`
	UpdateLabQuery = `UPDATE labs SET name = ?, description = ?, md_path = ?, deadline = ? WHERE id = ?`
	DeleteLabQuery = `DELETE FROM labs WHERE id = ?`
	DeleteLabFilesQuery = `DELETE FROM lab_files WHERE lab_id = ? AND file_path = ?`
)

