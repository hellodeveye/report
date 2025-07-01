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
func GenerateToken(openID, name string) (string, int64, error) {
	expireTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		OpenID: openID,
		Name:   name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime), // 24小时过期
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "report-assistant",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.GetJWTSecret()))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expireTime.Unix(), nil
}

// GenerateTokenFromUser 从用户对象生成JWT token (保持向后兼容)
func GenerateTokenFromUser(user models.User) (string, error) {
	token, _, err := GenerateToken(user.OpenID, user.Name)
	return token, err
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
