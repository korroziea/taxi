package auth

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Service interface {
	SignIn(ctx context.Context) error
	SignUp(ctx context.Context) error
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
		c.JSON(http.StatusNotImplemented, nil)
	}
}

func (h *Handler) signIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, nil)
	}
}
