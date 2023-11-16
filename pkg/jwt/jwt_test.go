package jwt

import (
	"testing"

	"github.com/google/uuid"
	"github.com/khhini/riset-gcp-api-gateway-auth/config"
	"github.com/stretchr/testify/assert"
)

func TestJwtCreateToken(t *testing.T) {
	userID := uuid.New()
	cfg := config.DefaultConfig("test", "v1")
	cfg.LoadFromEnv()

	token, err := CreateToken(cfg.JwtSignatureKey, userID)
	assert.Nil(t, err)
	assert.NotNil(t, token)
}
