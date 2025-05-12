package trip

import (
	"context"
	"html/template"
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
	Trips(ctx context.Context, userID string) ([]domain.Trip, error)
	Cost(ctx context.Context) (int64, error)
	Report(ctx context.Context) ([]byte, error)
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
	router.POST("/api/trips", h.middleware.VerifyUser, h.startTrip())
	router.PATCH("/api/trips/cancel", h.middleware.VerifyUser, h.cancelTrip())
	router.GET("/api/trips", h.middleware.VerifyUser, h.trips())

	router.GET("/api/trips/cost", h.middleware.VerifyUser, h.cost())

	router.POST("/api/report", h.middleware.VerifyUser, h.report())

	// templates
	router.GET("/trips", h.TripsRespTemp())
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

		trip := req.toDomain()
		trip.UserID = FromContext(ctx)

		err := h.service.StartTrip(ctx, trip)
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

func (h *Handler) trips() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := withKey(c)

		userID := FromContext(ctx)
		trips, err := h.service.Trips(ctx, userID)
		if err != nil {
			h.l.Error("service.Trips", zap.Error(err))

			response.TripError(c, err) // todo: error handling

			return
		}

		c.JSON(http.StatusOK, toTripsResp(trips))
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

func (h *Handler) report() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := withKey(c)

		report, err := h.service.Report(ctx)
		if err != nil {
			h.l.Error("service.Report", zap.Error(err))

			response.TripError(c, err) // todo: error handling

			return
		}

		c.Header("Content-Type", "application/pdf")
		c.Header("Content-Disposition", "attachment; filename=cars_report.pdf")
		c.Data(http.StatusOK, "application/pdf", report)
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

const rootFrontendPath = "internal/handler/frontend/"

func (h *Handler) TripsRespTemp() gin.HandlerFunc {
	return func(c *gin.Context) {
		temp := template.Must(template.ParseFiles(rootFrontendPath + "trips.html"))
		temp.Execute(c.Writer, nil)
	}
}
