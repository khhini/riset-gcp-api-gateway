package config

import "fmt"

type PostgresDBConfig struct {
	_dbConfig
}

func defaultPostgresDBConfig() dbConfigInterface {
	return &PostgresDBConfig{
		_dbConfig: _dbConfig{
			Host:     "127.0.0.1",
			Port:     5432,
			DBName:   "sample_app",
			SslMode:  "disable",
			Password: "defaultdbpassword",
		},
	}
}

func (cfg *PostgresDBConfig) LoadFromEnv() {
	loadEnvStr("POSTGRES_DB_HOST", &cfg.Host)
	loadEnvUint("POSTGRES_DB_PORT", &cfg.Port)
	loadEnvStr("POSTGRES_DB_NAME", &cfg.DBName)
	loadEnvStr("POSTGRES_DB_SSL_MODE", &cfg.SslMode)
	loadEnvStr("POSTGRES_DB_PASSWORD", &cfg.Password)
}

func (cfg *PostgresDBConfig) ConnStr() string {
	return fmt.Sprintf("host=%s port=%d user=postgres password=%s database=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Password, cfg.DBName, cfg.SslMode)
}

func (cfg *PostgresDBConfig) UnixConnStr() string {
	return fmt.Sprintf("user=postgres dbname=%s sslmode=%s host=%s", cfg.DBName, cfg.SslMode, cfg.Host)
}
