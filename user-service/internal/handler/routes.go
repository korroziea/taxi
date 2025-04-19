package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Router interface {
	InitRoutes(router gin.IRouter)
}

type Handler struct {
	User   Router
	Wallet Router
	Trip   Router
}

func New(
	user Router,
	wallet Router,
	trip Router,
) *Handler {
	handler := &Handler{
		User:   user,
		Wallet: wallet,
		Trip:   trip,
	}

	return handler
}

func (h *Handler) InitRoutes() http.Handler {
	r := gin.New()

	h.User.InitRoutes(r)
	h.Wallet.InitRoutes(r)
	h.Trip.InitRoutes(r)

	return r
}
