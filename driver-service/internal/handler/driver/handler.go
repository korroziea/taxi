package driver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/korroziea/taxi/driver-service/internal/config"
	"github.com/korroziea/taxi/driver-service/internal/domain"
	"github.com/korroziea/taxi/driver-service/internal/handler/response"
	"go.uber.org/zap"
)

type Cache interface {
	SetToken(ctx context.Context, driverID, token string) error
}

type Service interface {
	SignUp(ctx context.Context, driver domain.SignUpDriver) error
	SignIn(ctx context.Context, driver domain.SignInDriver) (string, error)
}

type Handler struct {
	l   *zap.Logger
	cfg config.Config

	cache Cache

	service Service
}

func New(
	l *zap.Logger,
	cfg config.Config,
	cache Cache,
	service Service,
) *Handler {
	handler := &Handler{
		l:   l,
		cfg: cfg,

		cache: cache,

		service: service,
	}

	return handler
}

func (h *Handler) InitRoutes(r gin.IRouter) {
	r.POST("/sign-up", h.signUp())
	r.POST("/sign-in", h.signIn())
}

func (h *Handler) signUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req signUpReq
		if err := c.ShouldBindJSON(&req); err != nil {
			h.l.Error("ShouldBindJSON", zap.Error(err))

			response.DriverError(c, err)

			return
		}

		if err := h.service.SignUp(c.Request.Context(), req.toDomain()); err != nil {
			h.l.Error("service.SignUp", zap.Error(err))

			response.DriverError(c, err)

			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

func (h *Handler) signIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		var req signInReq
		if err := c.ShouldBindJSON(&req); err != nil {
			h.l.Error("ShouldBindJSON", zap.Error(err))

			response.DriverError(c, err)

			return
		}

		driverID, err := h.service.SignIn(ctx, req.toDomain())
		if err != nil {
			h.l.Error("service.SignIn", zap.Error(err))

			response.DriverError(c, err)

			return
		}

		if err = h.genToken(ctx, driverID); err != nil {
			h.l.Error("genToken", zap.Error(err))

			response.DriverError(c, err)

			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

func (h *Handler) genToken(ctx context.Context, driverID string) error {
	payload := jwt.MapClaims{
		"sub": driverID,
		"exp": time.Now().Add(3600 * time.Second).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	signedToken, err := token.SignedString([]byte(h.cfg.SecretKey))
	if err != nil {
		return fmt.Errorf("token.SignedString: %w", err)
	}

	if err = h.cache.SetToken(ctx, driverID, signedToken); err != nil {
		return fmt.Errorf("cache.SetToken: %w", err)
	}

	return nil
}
