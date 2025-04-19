package trip

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/korroziea/taxi/user-service/internal/domain"
	"github.com/korroziea/taxi/user-service/internal/handler/middleware"
	"github.com/korroziea/taxi/user-service/internal/handler/response"
	"go.uber.org/zap"
)

type Middleware interface {
	VerifyUser(c *gin.Context)
}

type Service interface {
	StartTrip(ctx context.Context, trip domain.StartTrip) error
	CancelTrip(ctx context.Context) error
	Cost(ctx context.Context) (int64, error)
}

type Handler struct {
	l          *zap.Logger
	middleware Middleware
	service    Service
}

func New(
	l *zap.Logger,
	middleware Middleware,
	service Service,
) *Handler {
	handler := &Handler{
		l:          l,
		middleware: middleware,
		service:    service,
	}

	return handler
}

func (h *Handler) InitRoutes(router gin.IRouter) {
	router.POST("/trips", h.middleware.VerifyUser, h.startTrip())
	router.PATCH("/trips/cancel", h.middleware.VerifyUser, h.cancelTrip())

	router.GET("/trips/cost", h.middleware.VerifyUser, h.cost())
}

func (h *Handler) startTrip() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := withKey(c)

		var req startTripReq // todo:
		if err := c.ShouldBindJSON(&req); err != nil {
			h.l.Error("ShouldBindJSON", zap.Error(err))

			response.TripError(c, err) // todo: error handling

			return
		}

		err := h.service.StartTrip(ctx, req.toDomain())
		if err != nil {
			h.l.Error("service.StartTrip", zap.Error(err))

			response.TripError(c, err) // todo: error handling

			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

func (h *Handler) cancelTrip() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := withKey(c)

		if err := h.service.CancelTrip(ctx); err != nil {
			h.l.Error("service.CancelTrip", zap.Error(err))

			response.TripError(c, err) // todo: error handling

			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

func (h *Handler) cost() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := withKey(c)

		cost, err := h.service.Cost(ctx)
		if err != nil {
			h.l.Error("service.CancelTrip", zap.Error(err))

			response.TripError(c, err) // todo: error handling

			return
		}

		c.JSON(http.StatusNoContent, cost) // todo
	}
}

type userIDContextKey struct{}

func withKey(c *gin.Context) context.Context {
	userID := middleware.FromContext(c)

	return context.WithValue(c, userIDContextKey{}, userID)
}

func FromContext(ctx context.Context) string {
	return ctx.Value(userIDContextKey{}).(string)
}
