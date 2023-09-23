package middleware

import (
	"fmt"
	"log"
	"strings"

	"podcast/database"
	"podcast/repositories"
	"podcast/services"
	"podcast/types"
	"podcast/utils"

	"slices"

	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string

		authHeader := c.Request.Header.Get("Authorization")
		cookieToken, cookieErr := c.Cookie("access_token")

		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		} else if cookieErr == nil {
			token = cookieToken
		}

		if token == "" {
			utils.UnauthorizedResponse(c)
			return
		}

		payload, err := utils.ParseAccessJWT(token)
		if err != nil {
			log.Println("user access_token unuthorized")
			utils.UnauthorizedResponse(c)
			return
		}

		c.Set("user", payload)

		c.Next()
	}
}

func AuthenticateRefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string

		authHeader := c.Request.Header.Get("X-REFRESH-TOKEN")
		cookieToken, cookieErr := c.Cookie("refresh_token")

		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		} else if cookieErr == nil {
			token = cookieToken
		}

		if token == "" {
			utils.UnauthorizedResponse(c)
			return
		}

		payload, err := utils.ParseRefreshJWT(token)

		if err != nil {
			log.Println("user refresh_token unuthorized")
			utils.UnauthorizedResponse(c)
			return
		}

		c.Set("user", payload)

		c.Next()
	}
}

func Authorize(roles []types.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		u, _ := c.Get("user")
		p := u.(utils.JwtPayload)

		// TODO: DI
		as := services.NewAuthService(repositories.NewUsersRepository(database.DB))

		user, err := as.GetUser(fmt.Sprintf("%d", p.ID))

		if err != nil || !slices.Contains(roles, types.Role(user.Role)) {
			utils.ForbiddenResponse(c, "You lack the necesseary role to proceed")
			return
		}

		c.Next()
	}
}
