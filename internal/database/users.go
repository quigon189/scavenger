package database

import "scavenger/internal/models"

func (d *Database) GetAllUsers() []models.User {
	rows, err := d.db.Query()
}
