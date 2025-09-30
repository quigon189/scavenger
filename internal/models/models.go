package models

import "time"

type UserRole string

const (
	StudentRole UserRole = "student"
	AdminRole   UserRole = "admin"
)

type User struct {
	ID           int `yaml:"id"`
	Username     string `yaml:"username"`
	Name         string `yaml:"name"`
	PasswordHash string `yaml:"password"`

	RoleID       int `yaml:"role_id"`
	RoleName string `yaml:"role"`

	GroupID int
	GroupName string

}

type Group struct {
	ID          int
	Name        string  
	Disciplines []string 
}

type Discipline struct {
	ID   string `yaml:"id"`
	Name string `yaml:"name"`
	Labs []Lab  `yaml:"labs"`
}

type Lab struct {
	ID       string    `yaml:"id"`
	Name     string    `yaml:"name"`
	Path     string    `yaml:"path"`
	Deadline time.Time `yaml:"deadline,omitempty"`
}

type LabReport struct {
	ID         string    `yaml:"id"`
	Student    string    `yaml:"student"`
	Group      string    `yaml:"group"`
	Discipline string    `yaml:"discipline"`
	LabName    string    `yaml:"lab_name"`
	Path       string    `yaml:"path"`
	Comment    string    `yaml:"comment"`
	UploadetAt time.Time `yaml:"uploadet_at"`
	Status     string    `yaml:"stastus"`
	Grade      int       `yaml:"grade"`
}

type DatebaseConfig struct {
	DataSource     string `yaml:"data_source"`
	MigrationsPath string `yaml:"migrations_path"`
}

type ServerConfig struct {
	Port       string `yaml:"port"`
	UploadPath string `yaml:"upload_path"`
}

type AuthConfig struct {
	SessionSecret string `yaml:"session_secret"`
}

type TestDataConfig struct {
	Roles struct {
		Admin []User `yaml:"admin"`
		Student map[string][]User `yaml:"student"`
	} `yaml:"roles"`
}

type Config struct {
	Server      ServerConfig   `yaml:"server"`
	Auth        AuthConfig     `yaml:"auth"`
	DB          DatebaseConfig `yaml:"database"`
	Users       []User         
	Groups      []Group        
	Disciplines []Discipline   `yaml:"disciplines"`
	LabReports  []LabReport    `yaml:"lab_reports,omitempty"`
	TestData    TestDataConfig `yaml:"test_data"`
}
