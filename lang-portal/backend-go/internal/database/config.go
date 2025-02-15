package database

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

// Config holds the database configuration
type Config struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
}

// NewConfig creates a new database configuration from environment variables
func NewConfig() (*Config, error) {
	// Try to load .env from the project root
	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(filepath.Dir(filepath.Dir(filename))))
	if err := godotenv.Load(filepath.Join(projectRoot, ".env")); err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	return &Config{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		DBName:   os.Getenv("DB_NAME"),
	}, nil
}

// DSN returns the Data Source Name
func (c *Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", 
		c.User, c.Password, c.Host, c.Port, c.DBName)
}
