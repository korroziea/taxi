package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Router interface {
	InitRoutes(router gin.IRouter)
}

type Handler struct {
	Auth Router
}

func New(auth Router) *Handler {
	handler := &Handler{
		Auth: auth,
	}

	return handler
}

func (h *Handler) InitRoutes() http.Handler {
	r := gin.New()

	h.Auth.InitRoutes(r)

	return r
}
