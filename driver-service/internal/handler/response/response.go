package response

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/korroziea/taxi/driver-service/internal/domain"
)

type ErrResponse struct {
	ErrCode int    `json:"error_code"`
	Msg     string `json:"message"`
}

func DriverError(c *gin.Context, err error) {
	fmt.Println("response.DriverError: ", err) // todo: remove

	switch {
	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrResponse{
			ErrCode: -10,
			Msg:     "internal server error",
		})
	}
}

func AbortUnauthorized(c *gin.Context, err error) {
	fmt.Println("response.AbortUnauthorized: ", err) // todo: remove

	switch {
	case errors.Is(err, domain.ErrParseBearerToken):
		c.AbortWithStatusJSON(http.StatusUnauthorized, ErrResponse{
			ErrCode: -41,
			Msg:     "invalid specified token. check token type and length",
		})
	case errors.Is(err, domain.ErrParseJWTToken) ||
		errors.Is(err, domain.ErrParseClaims) ||
		errors.Is(err, domain.ErrFetchSub):
		c.AbortWithStatusJSON(http.StatusUnauthorized, ErrResponse{
			ErrCode: -41,
			Msg:     "invalid specified token",
		})
	case errors.Is(err, domain.ErrTokenNotFound):
		c.AbortWithStatusJSON(http.StatusUnauthorized, ErrResponse{
			ErrCode: -41,
			Msg:     "token expired. sign in again",
		})
	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrResponse{
			ErrCode: -10,
			Msg:     "internal server error",
		})
	}
}
