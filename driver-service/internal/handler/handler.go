package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
}

func New() *Handler {
	handler := &Handler{}

	return handler
}

func (h *Handler) InitRoutes() http.Handler {
	r := gin.New()

	r.POST("/ping", ping)

	return r
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, "Pong")
}
