package models

type DatebaseConfig struct {
	DataSource     string `yaml:"data_source"`
	MigrationsPath string `yaml:"migrations_path"`
}

type ServerConfig struct {
	Port       string `yaml:"port"`
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

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Auth     AuthConfig     `yaml:"auth"`
	DB       DatebaseConfig `yaml:"database"`
	FS       FSConfig       `yaml:"filestorage"`
	TestData TestDataConfig `yaml:"test_data"`
}
