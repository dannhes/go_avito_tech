package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go_avito_tech/internal/logger"
)

func LoggingMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			latency := time.Since(start)
			status := c.Response().Status
			req := c.Request()
			if err != nil {
				logger.L.Error("handler error",
					zap.String("method", req.Method),
					zap.String("path", req.URL.Path),
					zap.Int("status", status),
					zap.Duration("latency", latency),
					zap.String("remote_ip", c.RealIP()),
					zap.Error(err),
				)
				return err
			}
			logger.L.Info("http request",
				zap.String("method", req.Method),
				zap.String("path", req.URL.Path),
				zap.Int("status", status),
				zap.Duration("latency", latency),
				zap.String("remote_ip", c.RealIP()),
			)
			return nil
		}
	}
}
