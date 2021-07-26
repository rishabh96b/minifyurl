package config

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	DB *DBConfig
}

type DBConfig struct {
	Driver   string
	Host     string
	Port     string
	Username string
	Password string
	Name     string
}

func GetMySQlDBConfig() *Config {
	return &Config{
		DB: &DBConfig{
			Driver:   "mysql",
			Host:     os.Getenv("DB_HOST"),
			Port:     "3306",
			Username: os.Getenv("DB_USERNAME"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
		},
	}
}
