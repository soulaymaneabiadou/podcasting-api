package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func healthRoutes(r *gin.RouterGroup) {
	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"deps": map[string]string{
				"database": "ok",
			},
		})
	})
}
