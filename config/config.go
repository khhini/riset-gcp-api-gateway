package config

import "time"

type Config struct {
	AppName         string
	Version         string
	JwtSignatureKey []byte
	LoginExpiration time.Duration
	Listen          listenConfig
	DBConfig        dbConfigInterface
}

func DefaultConfig(appName string, version string) *Config {
	return &Config{
		AppName:         appName,
		Version:         version,
		JwtSignatureKey: []byte("super-duper-secret-signature-key"),
		LoginExpiration: time.Duration(1) * time.Hour,
		Listen:          defaultListenConfig(),
		DBConfig:        defaultPostgresDBConfig(),
	}
}

func (c *Config) LoadFromEnv() {
	loadEnvByte("JWT_SIGNATURE_KEY", &c.JwtSignatureKey)
	c.DBConfig.LoadFromEnv()
	c.Listen.LoadFromEnv()
}
