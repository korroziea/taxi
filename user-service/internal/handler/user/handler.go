package user

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/korroziea/taxi/user-service/internal/config"
	"github.com/korroziea/taxi/user-service/internal/domain"
	"github.com/korroziea/taxi/user-service/internal/handler/response"
	"go.uber.org/zap"
)

type Cache interface {
	SetToken(ctx context.Context, userID, token string) error
}

type Service interface {
	SignUp(ctx context.Context, user domain.SignUpUser) error
	SignIn(ctx context.Context, user domain.SignInUser) (string, error)
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
		l:   l, // todo: fix logger
		cfg: cfg,

		cache: cache,

		service: service,
	}

	return handler
}

func (h *Handler) InitRoutes(router gin.IRouter) {
	router.POST("/api/sign-up", h.signUp())
	router.POST("/api/sign-in", h.signIn())
}

func (h *Handler) signUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		var req signUpReq
		if err := c.ShouldBindJSON(&req); err != nil {
			h.l.Error("ShouldBindJSON", zap.Error(err))

			response.UserError(c, err)

			return
		}

		if err := h.service.SignUp(ctx, req.toDomain()); err != nil {
			h.l.Error("service.SignUp", zap.Error(err))

			response.UserError(c, err)

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
			response.UserError(c, err)

			return
		}

		userID, err := h.service.SignIn(ctx, req.toDomain())
		if err != nil {
			response.UserError(c, err)

			return
		}

		if err = h.genToken(ctx, userID); err != nil {
			response.UserError(c, err)

			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

func (h *Handler) genToken(ctx context.Context, userID string) error {
	payload := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(3600 * time.Second).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	signedToken, err := token.SignedString([]byte(h.cfg.SecretKey))
	if err != nil {
		return fmt.Errorf("token.SignedString: %w", err)
	}

	if err = h.cache.SetToken(ctx, userID, signedToken); err != nil {
		return fmt.Errorf("cache.SetToken: %w", err)
	}

	return nil
}
