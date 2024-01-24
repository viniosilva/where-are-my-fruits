package middlewares

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var skipPaths []string = []string{"/api/swagger/", "/api/healthcheck"}

func JSONLogMiddleware(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		for _, p := range skipPaths {
			if strings.HasPrefix(path, p) {
				return
			}
		}

		start := time.Now()
		query := c.Request.URL.RawQuery

		c.Next()

		duration := time.Now().Sub(start)
		requestID := c.Writer.Header().Get("Request-Id")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		loggerFn := logger.Infow
		if c.Writer.Status() >= http.StatusBadRequest {
			loggerFn = logger.Warnw
		}
		if c.Writer.Status() >= http.StatusInternalServerError {
			loggerFn = logger.Errorw
		}

		loggerFn("request",
			"client_ip", c.ClientIP(),
			"duration", duration.Milliseconds(),
			"method", c.Request.Method,
			"path", c.Request.RequestURI,
			"query", query,
			"status", c.Writer.Status(),
			"referrer", c.Request.Referer(),
			"request_id", requestID,
		)
	}
}
