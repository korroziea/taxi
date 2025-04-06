package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/korroziea/taxi/user-service/internal/domain"
)

type ErrResponse struct {
	ErrCode int    `json:"error_code"`
	Msg     string `json:"message"`
}

func Error(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrUserAlreadyExists):
		c.JSON(http.StatusBadRequest, ErrResponse{
			ErrCode: -40,
			Msg:     "user with the same phone already exists",
		})
	default:
		c.JSON(http.StatusInternalServerError, ErrResponse{
			ErrCode: -10,
			Msg:     "internal server error",
		})
	}
}
