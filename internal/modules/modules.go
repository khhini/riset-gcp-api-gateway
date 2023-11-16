package modules

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	helloHandler "github.com/khhini/riset-gcp-api-gateway-auth/internal/modules/hello/handler"

	userHandler "github.com/khhini/riset-gcp-api-gateway-auth/internal/modules/user/handler"
	userRepository "github.com/khhini/riset-gcp-api-gateway-auth/internal/modules/user/repository"
	userService "github.com/khhini/riset-gcp-api-gateway-auth/internal/modules/user/service"
)

type Modules struct {
	PgPool          *pgxpool.Pool
	JwtSignatureKey []byte
}

func NewModules(pool *pgxpool.Pool, jwtSigatureKey []byte) *Modules {
	return &Modules{
		PgPool:          pool,
		JwtSignatureKey: jwtSigatureKey,
	}
}

func (m *Modules) SetupRouter() *gin.Engine {
	router := gin.Default()

	userRepo := userRepository.NewPostgreRepository(m.PgPool)
	userSvc := userService.NewService(userRepo)

	api := router.Group("/api")
	{
		helloHandler.NewHandler(api)
		userHandler.NewHandler(api, userSvc)
	}

	return router
}
