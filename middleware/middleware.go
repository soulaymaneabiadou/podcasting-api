package middleware

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

func Compression() gin.HandlerFunc {
	return gzip.Gzip(gzip.DefaultCompression)
}

func RequestId() gin.HandlerFunc {
	return requestid.New()
}
