package utils

import (
	"io"
	"net/http"
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
