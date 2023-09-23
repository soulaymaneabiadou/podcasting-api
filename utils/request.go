package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetCtxUser(c *gin.Context) (string, JwtPayload) {
	user, _ := c.Get("user")
	id := strconv.Itoa(int(user.(JwtPayload).ID))

	return id, user.(JwtPayload)
}
