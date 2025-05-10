package wallet

import (
	"context"
	"net/http"
	"html/template"

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
	Refill(ctx context.Context, walletID string, amount int64) (domain.ViewWallet, error)
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
	router.POST("/api/wallets", h.middleware.VerifyUser, h.createWallet())
	router.GET("/api/wallets", h.middleware.VerifyUser, h.walletList())
	router.GET("/api/wallets/:id", h.middleware.VerifyUser, h.wallet())
	router.PUT("/api/wallets/:id/type", h.middleware.VerifyUser, h.changeType())
	router.PUT("/api/wallets/:id/refill", h.middleware.VerifyUser, h.refill())
	
	router.GET("/wallets", h.wallets())
	router.GET("/wallets/:id", h.oneWallet())
	router.GET("/wallets/:id/refill", h.refillWallet())
}

func (h *Handler) createWallet() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := withKey(c)

		wallet, err := h.service.CreateWallet(ctx)
		if err != nil {
			h.l.Error("service.CreateWallet", zap.Error(err))

			response.WalletError(c, err) // todo: error handling

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

			response.WalletError(c, err) // todo: error handling

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

			response.WalletError(c, err) // todo: error handling

			return
		}

		c.JSON(http.StatusOK, toWalletView(wallet))
	}
}

func (h *Handler) changeType() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := withKey(c)

		wallet, err := h.service.ChangeType(ctx, c.Param(idURLParam))
		if err != nil {
			h.l.Error("service.ChangeType", zap.Error(err))

			response.WalletError(c, err) // todo: error handling

			return
		}

		c.JSON(http.StatusOK, toWalletView(wallet))
	}
}

func (h *Handler) refill() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := withKey(c)

		var req refillReq
		if err := c.ShouldBindJSON(&req); err != nil {
			h.l.Error("ShouldBindJSON", zap.Error(err))

			response.WalletError(c, err) // todo: error handling

			return
		}

		wallet, err := h.service.Refill(ctx, c.Param(idURLParam), req.Amount)
		if err != nil {
			h.l.Error("service.Refill", zap.Error(err))

			response.WalletError(c, err) // todo: error handling

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

const rootFrontendPath = "internal/handler/frontend/"

func (h *Handler) wallets() gin.HandlerFunc {
	return func(c *gin.Context) {
		temp := template.Must(template.ParseFiles(rootFrontendPath + "wallets.html"))
		temp.Execute(c.Writer, nil)
	}
}

func (h *Handler) oneWallet() gin.HandlerFunc {
	return func(c *gin.Context) {
		temp := template.Must(template.ParseFiles(rootFrontendPath + "wallet.html"))
		temp.Execute(c.Writer, nil)
	}
}

func (h *Handler) refillWallet() gin.HandlerFunc {
	return func(c *gin.Context) {
		temp := template.Must(template.ParseFiles(rootFrontendPath + "refill.html"))
		temp.Execute(c.Writer, nil)
	}
}
