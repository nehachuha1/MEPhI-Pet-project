package config

import "os"

type databaseConfig struct {
	PgxUser     string
	PgxPassword string
	PgxAddress  string
	PgxPort     string
	PgxDB       string
	RedisURL    string
}

type gRPC struct {
	Port string
}

type Config struct {
	Database databaseConfig
	GRPC     gRPC
}

func NewConfig() *Config {
	return &Config{
		Database: databaseConfig{
			PgxUser:     getEnv("POSTGRES_USER", ""),
			PgxPassword: getEnv("POSTGRES_PASSWORD", ""),
			PgxAddress:  getEnv("POSTGRES_ADDRESS", "localhost"),
			PgxPort:     getEnv("POSTGRES_PORT", "5432"),
			PgxDB:       getEnv("POSTGRES_DB_NAME", ""),
			RedisURL:    getEnv("REDIS_URL", ""),
		},
		GRPC: gRPC{
			Port: getEnv("GRPC_PORT", ""),
		},
	}
}

func getEnv(name string, defaultValue string) string {
	if value, isExists := os.LookupEnv(name); isExists {
		return value
	}
	return defaultValue
}
