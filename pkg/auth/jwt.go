package auth

import (
	"context"
	"errors"
	"time"

	"github.com/anggitaputra13/mcs-auth/pkg/redis"
	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	secret     string
	expiration time.Duration
	redis      *redis.Client
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func NewJWT(secret string, expiration time.Duration, redisClient *redis.Client) *JWT {
	return &JWT{
		secret:     secret,
		expiration: expiration,
		redis:      redisClient,
	}
}

func (j *JWT) GenerateToken(userID string) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}

func (j *JWT) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token")
	}

	// Check if token is blacklisted in Redis
	_, err = j.redis.Get(context.Background(), "blacklist:"+tokenString).Result()
	if err == nil {
		return "", errors.New("token is blacklisted")
	}

	return claims.UserID, nil
}

func (j *JWT) BlacklistToken(token string) error {
	userID, err := j.ValidateToken(token)
	if err != nil {
		return err
	}

	// Get expiration time from token
	parsedToken, _ := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})
	expTime := parsedToken.Claims.(*Claims).ExpiresAt.Time

	// Set blacklist with TTL
	ttl := time.Until(expTime)
	return j.redis.Set(context.Background(), "blacklist:"+token, userID, ttl).Err()
}
