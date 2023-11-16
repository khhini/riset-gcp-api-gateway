package config

type listenConfig struct {
	Host string `yaml:"host" json:"host"`
	Port uint   `yaml:"port" json:"port"`
}

func defaultListenConfig() listenConfig {
	return listenConfig{
		Host: "0.0.0.0",
		Port: 8080,
	}
}

func (l *listenConfig) LoadFromEnv() {
	loadEnvStr("HOST", &l.Host)
	loadEnvUint("PORT", &l.Port)
}
