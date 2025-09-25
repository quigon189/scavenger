package database

const (
	// Roles
	GetRoleByName = `SELECT id, name, description FROM roles WHERE name = ?`
	// Groups
	GetAllGroupsQuery = `SELECT id, name FROM groups`

	// Users
	CreateUserWithRoleQuery = `INSERT INTO users (username, name, password_hash, role_id) VALUES (?, ?, ?, ?)`
	CreateUserWithRoleAndGroupQuery = `INSERT INTO users (username, name, password_hash, role_id, group_id) VALUES (?, ?, ?, ?, ?)`
	GetUserByUsernameQuery = `
		SELECT u.id, u.username, u.name, u.password_hash, r.name AS role_name, g.name AS group_name
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		LEFT JOIN groups g ON u.group_id = g.id
		WHERE u.username = ?
	`
)

