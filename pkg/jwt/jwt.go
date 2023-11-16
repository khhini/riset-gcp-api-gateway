package jwt

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/dgrijalva/jwt-go.v3"
)

type Claims struct {
	UserID uuid.UUID `json:"userId"`
	jwt.StandardClaims
}

func CreateToken(jwtSignatureKey []byte, userID uuid.UUID) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSignatureKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
