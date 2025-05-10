package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Router interface {
	InitRoutes(router gin.IRouter)
}

type Handler struct {
	User Router
}

func New(
	user Router,
) *Handler {
	handler := &Handler{
		User: user,
	}

	return handler
}

func (h *Handler) InitRoutes() http.Handler {
	r := gin.New()

	h.User.InitRoutes(r)

	return r
}
