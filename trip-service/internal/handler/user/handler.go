package user

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/korroziea/taxi/trip-service/internal/config"
	"github.com/korroziea/taxi/trip-service/internal/domain"
	"go.uber.org/zap"
)

type Service interface {
	Trips(ctx context.Context, userID string) ([]domain.Trip, error)
}

type Handler struct {
	l   *zap.Logger
	cfg config.Config

	service Service
}

func New(
	l *zap.Logger,
	cfg config.Config,
	service Service,
) *Handler {
	handler := &Handler{
		l:   l, // todo: fix logger
		cfg: cfg,

		service: service,
	}

	return handler
}

func (h *Handler) InitRoutes(router gin.IRouter) {
	router.GET("/api/trips", h.trips())
}

func (h *Handler) trips() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		var req tripsReq
		if err := c.ShouldBindJSON(&req); err != nil {
			h.l.Error("ShouldBindJSON", zap.Error(err))

			return
		}

		trips, err := h.service.Trips(ctx, req.UserID)
		if err != nil {
			h.l.Error("service.Trips", zap.Error(err))

			return
		}

		c.JSON(http.StatusOK, toTripsResp(trips))
	}
}
