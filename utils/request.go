package utils

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"podcast/types"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetCtxUser(c *gin.Context) (string, JwtPayload) {
	user, _ := c.Get("user")
	id := strconv.Itoa(int(user.(JwtPayload).ID))

	return id, user.(JwtPayload)
}

// Limit the request size in MB
func LimitRequestSize(c *gin.Context, maxMb int64) io.ReadCloser {
	return http.MaxBytesReader(c.Writer, c.Request.Body, int64(maxMb<<20))
}

func IsOwner(c *gin.Context, ownerId uint) bool {
	id, _ := GetCtxUser(c)
	return fmt.Sprint(ownerId) == id
}

func UserHasActiveStripeAccount(c *gin.Context, f func(string) (types.User, error)) (string, error) {
	id, _ := GetCtxUser(c)

	user, err := f(id)
	if err != nil {
		// utils.ErrorResponse(c, err, "User not found")
		return "", err
	}

	if user.StripeAccountId == "" || user.PayoutsEnabled == false {
		// utils.ErrorResponse(c, err, "No connect account was found, please start by creating one through the connect flow")
		return "", errors.New("user has no active stripe account id")
	}

	return user.StripeAccountId, nil
}
