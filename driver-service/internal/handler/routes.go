package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Router interface {
	InitRoutes(r gin.IRouter)
}

type Handler struct {
	driver Router
}

func New(driver Router) *Handler {
	handler := &Handler{
		driver: driver,
	}

	return handler
}

func (h *Handler) InitRoutes() http.Handler {
	r := gin.New()

	r.POST("/ping", ping)

	h.driver.InitRoutes(r)

	return r
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, "Pong")
}
