package config

import (
	"os"
	"strconv"
)

func loadEnvStr(key string, result *string) {
	s, ok := os.LookupEnv(key)
	if !ok {
		return
	}

	*result = s
}

func loadEnvByte(key string, result *[]byte) {
	s, ok := os.LookupEnv(key)
	if !ok {
		return
	}

	*result = []byte(s)
}

func loadEnvUint(key string, result *uint) {
	s, ok := os.LookupEnv(key)
	if !ok {
		return
	}

	n, err := strconv.Atoi(s)

	if err != nil {
		return
	}

	*result = uint(n)
}

type _dbConfig struct {
	Host     string
	Port     uint
	DBName   string
	SslMode  string
	Password string
}

type dbConfigInterface interface {
	LoadFromEnv()
	ConnStr() string
	UnixConnStr() string
}
