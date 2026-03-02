package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
}

var (
	config    *Config
	alpsToken string
)

func InitConfig() (err error) {
	if err = godotenv.Load(); err != nil {
		err = errors.Wrapf(err, "failed to load config: %+v", err)
		return
	}
	config = &Config{
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBName:     os.Getenv("DB_NAME"),
	}
	return
}

func (c *Config) GetDsn() (dsn string) {
	Assert(c.DBUser != "", "config db user must not be empty")
	Assert(c.DBPassword != "", "config db password must not be empty")
	Assert(c.DBHost != "", "config db host must not be empty")
	Assert(c.DBPort != "", "config db port must not be empty")
	Assert(c.DBName != "", "config db name must not be empty")

	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
	return
}
