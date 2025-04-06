package user

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/korroziea/taxi/user-service/internal/domain"
	"github.com/korroziea/taxi/user-service/internal/handler/response"
)

type Service interface {
	SignIn(ctx context.Context, user domain.SignInUser) error
	SignUp(ctx context.Context, user domain.SignUpUser) error
}

type Handler struct {
	service Service
}

func New(service Service) *Handler {
	handler := &Handler{
		service: service,
	}

	return handler
}

func (h *Handler) InitRoutes(router gin.IRouter) {
	router.POST("/sign-up", h.signUp())
	router.POST("/sign-in", h.signIn())
}

func (h *Handler) signUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		var req signUpReq
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, err)

			return
		}

		if err := h.service.SignUp(ctx, req.toDomain()); err != nil {
			response.Error(c, err)

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
			response.Error(c, err)

			return
		}

		if err := h.service.SignIn(ctx, req.toDomain()); err != nil {
			response.Error(c, err)

			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}
