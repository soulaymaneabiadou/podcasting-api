package middleware

import (
	"log/slog"
	"net/http"
	"podcast/logger"
	"time"

	"github.com/gin-gonic/gin"
)

type Logger func(string, ...any)

func HttpLogger() gin.HandlerFunc {
	logger := logger.NewLogger()

	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		l := createHttpLogger(logger, c)
		p := createRequestLogParams(start, c)
		logRequestParams(l, p)
	}
}

func createHttpLogger(logger *slog.Logger, c *gin.Context) Logger {
	var log = logger.Debug

	if c.Writer.Status() >= http.StatusInternalServerError {
		log = logger.Error
	} else if c.Writer.Status() >= http.StatusBadRequest {
		log = logger.Warn
	} else {
		log = logger.Info
	}

	return log
}

func createRequestLogParams(start time.Time, c *gin.Context) gin.LogFormatterParams {
	path := c.Request.URL.Path
	query := c.Request.URL.RawQuery
	params := gin.LogFormatterParams{}

	params.ClientIP = c.ClientIP()
	params.TimeStamp = time.Now()
	params.Latency = params.TimeStamp.Sub(start)
	if params.Latency > time.Minute {
		params.Latency = params.Latency.Truncate(time.Second)
	}

	params.Method = c.Request.Method
	params.StatusCode = c.Writer.Status()
	params.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
	params.BodySize = c.Writer.Size()

	if query != "" {
		path = path + "?" + query
	}
	params.Path = path

	return params
}

func logRequestParams(log Logger, params gin.LogFormatterParams) {

	log(params.ErrorMessage,
		slog.String("protocol", "http"), // params.Request.Proto,
		slog.String("method", params.Method),
		slog.String("ip", params.ClientIP),
		slog.String("path", params.Path),
		slog.Int("status", params.StatusCode),
		slog.String("latency", params.Latency.String()),
		slog.Int("body_size", params.BodySize),
		// "user_agent", params.Request.UserAgent(),
	)
}
