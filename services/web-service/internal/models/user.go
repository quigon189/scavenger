package models

type User struct {
	Username string
	Name     string
	Role     string
	Email    string
	Group    string
	Status   string
	GroupID  *int
}

type Group struct {
	ID   int
	Name string
}
