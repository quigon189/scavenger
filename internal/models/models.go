package models

import "time"

type UserRole string

const (
	StudentRole UserRole = "student"
	AdminRole   UserRole = "admin"
)

type User struct {
	ID           int    `yaml:"id"`
	Username     string `yaml:"username"`
	Name         string `yaml:"name"`
	PasswordHash string `yaml:"password"`

	RoleID   int    `yaml:"role_id"`
	RoleName string `yaml:"role"`

	GroupID   int
	GroupName string
}

type Group struct {
	ID          int
	Name        string
	Disciplines []Discipline
	Students    []User
}

type Discipline struct {
	ID      int
	Name    string
	GroupID *int

	Group Group
	Labs  []Lab
}

type Lab struct {
	ID           string
	Name         string
	MDFileID     int
	Deadline     time.Time
	Description  string
	DisciplineID int

	MDFile      StoredFile
	StoredFiles []StoredFile
	Reports     []LabReport
}

type StoredFile struct {
	ID       int
	Filename string
	Path     string
	URL      string
	Size     int64
}

type LabReport struct {
	ID         string
	Student    string
	Group      string
	Discipline string
	LabName    string
	Path       string
	Comment    string
	UploadedAt time.Time
	Status     string
	Grade      int
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

type FSConfig struct {
	BasePath string `yaml:"base_path"`
	BaseURL  string `yaml:"base_url"`
}

type TestDataConfig struct {
	Roles struct {
		Admin   []User            `yaml:"admin"`
		Student map[string][]User `yaml:"student"`
	} `yaml:"roles"`
}

type AdminStats struct {
	TotalReports   int
	PendingReports int
	GradedReports  int
	TotalGroups    int
}

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Auth     AuthConfig     `yaml:"auth"`
	DB       DatebaseConfig `yaml:"database"`
	FS       FSConfig       `yaml:"filestorage"`
	TestData TestDataConfig `yaml:"test_data"`
}
