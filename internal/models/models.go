package models

import "time"

type User struct {
	Username string `yaml:"username"`
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
	Role     string `yaml:"role"`
	Group    string `yaml:"group,omitempty"`
}

type Group struct {
	Name        string   `yaml:"name"`
	Disciplines []string `yaml:"disciplines"`
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

type ServerConfig struct {
	Port       string `yaml:"port"`
	UploadPath string `yaml:"upload_path"`
}

type AuthConfig struct {
	SessionSecret string `yaml:"session_secret"`
}

type Config struct {
	Server      ServerConfig `yaml:"server"`
	Auth        AuthConfig   `yaml:"auth"`
	Users       []User       `yaml:"users"`
	Groups      []Group      `yaml:"groups"`
	Disciplines []Discipline `yaml:"disciplines"`
	LabReports  []LabReport  `yaml:"lab_reports,omitempty"`
}
