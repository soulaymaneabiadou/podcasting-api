package middleware

import (
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
)

func Secure() gin.HandlerFunc {
	config := secure.DefaultConfig()
	config.IsDevelopment = os.Getenv("GIN_MODE") != "release"

	return secure.New(config)
}

func Cors() gin.HandlerFunc {
	config := cors.Config{
		// AllowAllOrigins:        os.Getenv("GIN_MODE") != "release",
		AllowOrigins:           []string{os.Getenv("PUBLIC_URL")},
		AllowMethods:           []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:           []string{"Origin", "Content-Length", "Content-Type"},
		ExposeHeaders:          []string{"Content-Length"},
		AllowCredentials:       true,
		MaxAge:                 12 * time.Hour,
		AllowWildcard:          false,
		AllowBrowserExtensions: false,
		AllowWebSockets:        false,
		AllowFiles:             false,
	}

	return cors.New(config)
}
