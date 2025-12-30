package config

import "os"

type Config struct {
	DBUser string
	DBPass string
	DBHost string
	DBPort string
	DBName string
}

func LoadConfig() Config {
	return Config{
		DBUser: getenv("DB_USER", "root"),
		DBPass: getenv("DB_PASS", "1234"),
		DBHost: getenv("DB_HOST", "127.0.0.1"),
		DBPort: getenv("DB_PORT", "3307"),
		DBName: getenv("DB_NAME", "digo"),
	}
}

func getenv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}
