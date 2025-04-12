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
	fmt.Println("response.Error: ", err) // todo: remove

	switch {
	case errors.Is(err, domain.ErrUserAlreadyExists):
		c.JSON(http.StatusBadRequest, ErrResponse{
			ErrCode: -40,
			Msg:     "user with the same phone already exists",
		})
	// todo: error handling
	case errors.Is(err, domain.ErrUserNotFound) || errors.Is(err, domain.ErrWrongPassword):
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
