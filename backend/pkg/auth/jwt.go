package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/hellodeveye/report/internal/config"
	"github.com/hellodeveye/report/internal/models"
)

// Claims JWT声明
type Claims struct {
	OpenID string `json:"open_id"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT token
func GenerateToken(user models.User) (string, error) {
	claims := &Claims{
		OpenID: user.OpenID,
		Name:   user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24小时过期
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "feishu-report-assistant",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.GetJWTSecret()))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken 验证JWT token
func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetJWTSecret()), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.NewValidationError("invalid token", jwt.ValidationErrorClaimsInvalid)
}
