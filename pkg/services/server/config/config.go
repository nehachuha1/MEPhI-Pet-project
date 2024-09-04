package config

import "os"

type sessionConfig struct {
	JwtKey     string
	SessionKey string
}

type databaseConfig struct {
	PgxUser     string
	PgxPassword string
	PgxAddress  string
	PgxPort     string
	PgxDB       string
	RedisURL    string
}

type Config struct {
	Sess     sessionConfig
	Database databaseConfig
	GrpcPort string
}

func NewConfig() *Config {
	return &Config{
		Sess: sessionConfig{
			JwtKey:     getEnv("JWT_SECRET_KEY", ""),
			SessionKey: getEnv("SESSION_KEY", ""),
		},
		Database: databaseConfig{
			PgxUser:     getEnv("POSTGRES_USER", ""),
			PgxPassword: getEnv("POSTGRES_PASSWORD", ""),
			PgxAddress:  getEnv("POSTGRES_ADDRESS", "localhost"),
			PgxPort:     getEnv("POSTGRES_PORT", "5432"),
			PgxDB:       getEnv("POSTGRES_DB_NAME", ""),
			RedisURL:    getEnv("REDIS_URL", ""),
		},
		GrpcPort: getEnv("GRPC_PORT", ""),
	}
}

func getEnv(name string, defaultValue string) string {
	if value, isExists := os.LookupEnv(name); isExists {
		return value
	}
	return defaultValue
}
