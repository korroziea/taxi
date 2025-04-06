package response

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/korroziea/taxi/user-service/internal/domain"
)

type ErrResponse struct {
	ErrCode int    `json:"error_code"`
	Msg     string `json:"message"`
}

func Error(c *gin.Context, err error) {
	fmt.Println(err)

	switch {
	case errors.Is(err, domain.ErrUserAlreadyExists):
		c.JSON(http.StatusBadRequest, ErrResponse{
			ErrCode: -40,
			Msg:     "user with the same phone already exists",
		})
	case errors.Is(err, domain.ErrUserNotFound):
		c.JSON(http.StatusBadRequest, ErrResponse{
			ErrCode: -44,
			Msg:     "phone or password is incorrect",
		})
	default:
		c.JSON(http.StatusInternalServerError, ErrResponse{
			ErrCode: -10,
			Msg:     "internal server error",
		})
	}
}
