package service

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

//JWTManager is json web token manager
type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

//UserClaims is a custom JWT claims that contains some user information
type UserClaims struct {
	jwt.StandardClaims
	User string `json:"username"`
	Role string `json:"role"`
	DB   string `json:"db"`
}

func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{secretKey, tokenDuration}
}
func (manager *JWTManager) Generate(user *User) (string, error) {
	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(manager.tokenDuration).Unix(),
		},
		User: user.User,
		Role: user.Role,
		DB:   user.DB,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(manager.secretKey))
}

//Verify verifies the acces token string an return a use claim if it es valid
func (manager *JWTManager) Verify(accesToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accesToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				//to check that signing method is no modified
				return nil, fmt.Errorf("unexpected token signing method")
			}
			return []byte(manager.secretKey), nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}
