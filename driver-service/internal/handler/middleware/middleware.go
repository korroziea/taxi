package middleware

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/korroziea/taxi/driver-service/internal/config"
	"github.com/korroziea/taxi/driver-service/internal/domain"
	"github.com/korroziea/taxi/driver-service/internal/handler/response"
)

const (
	authorizationHeader = "Authorization"
	bearerPrefix        = "Bearer "

	driverIDContextKey = "X-Driver-ID-Context-Key"
)

type Cache interface {
	GetToken(ctx context.Context, driverID string) (string, error)
}

type Middleware struct {
	cfg   config.Config
	cache Cache
}

func New(cfg config.Config, cache Cache) *Middleware {
	middleware := &Middleware{
		cfg:   cfg,
		cache: cache,
	}

	return middleware
}

func (m *Middleware) VerifyUser(c *gin.Context) {
	auth := c.GetHeader(authorizationHeader)
	if auth == "" {
		response.AbortUnauthorized(c, domain.ErrParseBearerToken)

		return
	}

	if len(auth) < len(bearerPrefix) {
		response.AbortUnauthorized(c, domain.ErrParseBearerToken)

		return
	}

	if auth[:len(bearerPrefix)] != bearerPrefix {
		response.AbortUnauthorized(c, domain.ErrParseBearerToken)

		return
	}

	tokenString := auth[len(bearerPrefix):]
	tokenString = strings.TrimSpace(tokenString)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.cfg.SecretKey), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		response.AbortUnauthorized(c, domain.ErrParseJWTToken)

		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		response.AbortUnauthorized(c, domain.ErrParseClaims)

		return
	}

	driverID, err := claims.GetSubject()
	if err != nil {
		response.AbortUnauthorized(c, domain.ErrFetchSub)

		return
	}

	val, err := m.cache.GetToken(c.Request.Context(), driverID)
	if err != nil || val != tokenString {
		response.AbortUnauthorized(c, domain.ErrTokenNotFound)

		return
	}

	withKey(c, driverID)

	c.Next()
}

func withKey(c *gin.Context, driverID string) {
	c.Set(driverIDContextKey, driverID)
}

func FromContext(c *gin.Context) string {
	val, _ := c.Get(driverIDContextKey)

	driverID, _ := val.(string)

	return driverID
}
