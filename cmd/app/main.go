package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/khhini/riset-gcp-api-gateway-auth/config"
	"github.com/khhini/riset-gcp-api-gateway-auth/internal/modules"
	"github.com/khhini/riset-gcp-api-gateway-auth/pkg/postgres"
)

var cfg *config.Config
var pool *pgxpool.Pool
var ctx context.Context
var err error

func init() {
	ctx = context.Background()
	cfg = config.DefaultConfig("sample-app", "v1")
	cfg.LoadFromEnv()

	pool, err = postgres.ConnectPostgresDB(ctx, cfg.DBConfig.ConnStr())
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("DB Connection is alive")
}

func main() {
	mod := modules.NewModules(pool, cfg.JwtSignatureKey)
	router := mod.SetupRouter()

	router.Run()
}
