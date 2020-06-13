package util

import (
	"github.com/dgrijalva/jwt-go"
	"seabase/extend/conf"
	"time"
)

type Claims struct {
	UserName string `json:"userName"`
	RedisKey string `json:"redisKey"`
	jwt.StandardClaims
}

// GenerateToken generate tokens used for auth
func GenerateToken(userName, redisKey string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(time.Duration(conf.ServerConf.JWTExpire) * time.Hour)

	claims := Claims{
		EncodeMD5(userName),
		redisKey,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "seabase",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString([]byte(conf.ServerConf.JWTSecret))
	return token, err
}

// ParseToken parsing token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.ServerConf.JWTSecret), nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
