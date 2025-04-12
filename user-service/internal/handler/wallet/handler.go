package wallet

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/korroziea/taxi/user-service/internal/domain"
	"github.com/korroziea/taxi/user-service/internal/handler/middleware"
	"github.com/korroziea/taxi/user-service/internal/handler/response"
	"go.uber.org/zap"
)

const idURLParam = "id"

type Middleware interface {
	VerifyUser(c *gin.Context)
}

type Service interface {
	CreateWallet(ctx context.Context) (domain.ViewWallet, error)
	WalletList(ctx context.Context) ([]domain.ViewWallet, error)
	Wallet(ctx context.Context, walletID string) (domain.ViewWallet, error)
	ChangeType(ctx context.Context, walletID string) (domain.ViewWallet, error)
	Refill(ctx context.Context, amount int64) (domain.ViewWallet, error)
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
	router.POST("/wallets", h.middleware.VerifyUser, h.createWallet())
	router.GET("/wallets", h.middleware.VerifyUser, h.walletList())
	router.GET("/wallets/:id", h.middleware.VerifyUser, h.wallet())
}

func (h *Handler) createWallet() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := withKey(c)

		wallet, err := h.service.CreateWallet(ctx)
		if err != nil {
			h.l.Error("service.CreateWallet", zap.Error(err))

			response.UserError(c, err) // todo: error handling

			return
		}

		c.JSON(http.StatusCreated, toWalletView(wallet))
	}
}

func (h *Handler) walletList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := withKey(c)

		wallets, err := h.service.WalletList(ctx)
		if err != nil {
			h.l.Error("service.WalletList", zap.Error(err))

			response.UserError(c, err) // todo: error handling

			return
		}

		c.JSON(http.StatusOK, toWalletListView(wallets))
	}
}

func (h *Handler) wallet() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := withKey(c)

		wallet, err := h.service.Wallet(ctx, c.Param(idURLParam))
		if err != nil {
			h.l.Error("service.Wallet", zap.Error(err))

			response.UserError(c, err) // todo: error handling

			return
		}

		c.JSON(http.StatusOK, toWalletView(wallet))
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
