package http

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLoggerMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		status := c.Writer.Status()
		attrs := []any{
			slog.String("method", c.Request.Method),
			slog.String("path", c.FullPath()),
			slog.Int("status", status),
			slog.Duration("latency", time.Since(start)),
		}

		if len(c.Errors) > 0 {
			attrs = append(attrs, slog.String("error", c.Errors.Last().Error()))
		}

		switch {
		case status >= 500:
			logger.Error("request failed", attrs...)
		case status >= 400:
			logger.Warn("request rejected", attrs...)
		default:
			logger.Info("request handled", attrs...)
		}
	}
}
