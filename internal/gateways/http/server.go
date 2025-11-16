package http

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go_avito_tech/api/gen"
	"go_avito_tech/internal/logger"
	"go_avito_tech/internal/middleware"
	"go_avito_tech/internal/repository"
)

type Config struct {
	Host string
	Port uint16
}

type Server struct {
	config Config
	echo   *echo.Echo
}

type UseCases struct {
	Users  repository.UserRepository
	Teams  repository.TeamRepository
	PullRs repository.PullRequestRepository
	Stats  repository.StatsRepository
}

func NewServer(cfg Config, useCases UseCases) *Server {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.LoggingMiddleware())
	h := NewHandler(useCases)
	gen.RegisterHandlers(e, h)
	s := &Server{
		config: cfg,
		echo:   e,
	}
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := 500
		msg := err.Error()
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			msg = fmt.Sprintf("%v", he.Message)
		}
		logger.L.Error("http_error",
			zap.String("method", c.Request().Method),
			zap.String("path", c.Path()),
			zap.Int("status", code),
			zap.String("error", msg),
		)
		err = c.JSON(code, map[string]string{"error": msg})
		if err != nil {
			return
		}
	}
	return s
}

func (s *Server) Run() error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	return s.echo.Start(addr)
}
