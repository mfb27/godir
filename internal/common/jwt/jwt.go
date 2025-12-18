package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	secretKey []byte
	tokenExp  time.Duration
)

type Claims struct {
	UserID   uint   `json:"userId"`
	Username string `json:"username"`
	Policy   string `json:"policy"`
	jwt.RegisteredClaims
}

// Init 初始化JWT配置
func Init(secretKeyStr string, tokenExpStr string) error {
	if secretKeyStr == "" {
		return fmt.Errorf("JWT SecretKey不能为空")
	}
	secretKey = []byte(secretKeyStr)

	if tokenExpStr == "" {
		tokenExp = 24 * time.Hour // 默认24小时
	} else {
		var err error
		tokenExp, err = time.ParseDuration(tokenExpStr)
		if err != nil {
			return fmt.Errorf("JWT TokenExp格式错误: %w", err)
		}
	}

	return nil
}

// GenerateToken 生成JWT token
func GenerateToken(userID uint, username string) (string, error) {
	if secretKey == nil {
		return "", fmt.Errorf("JWT未初始化，请先调用jwt.Init")
	}

	claims := Claims{
		UserID:   userID,
		Username: username,
		Policy:   "readwrite",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// ParseToken 解析JWT token
func ParseToken(tokenString string) (*Claims, error) {
	if secretKey == nil {
		return nil, fmt.Errorf("JWT未初始化，请先调用jwt.Init")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
