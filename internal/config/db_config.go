package config

import (
	"os"
	"strconv"
)

type DB struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

func LoadDB() DB {
	return DB{
		Host:     getEnv("HOST", "localhost"), // контейнер БД в docker-compose
		Port:     getEnvInt("DB_PORT", 5432),
		User:     getEnv("DB_USER", "user"),
		Password: getEnv("DB_PASSWORD", "pass"),
		Name:     getEnv("DB_NAME", "mydb"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

func getEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}
