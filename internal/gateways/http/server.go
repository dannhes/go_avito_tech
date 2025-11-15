package http

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"go_avito_tech/api/gen"
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
}

func NewServer(cfg Config, useCases UseCases) *Server {
	e := echo.New()
	e.HideBanner = true
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
		c.JSON(code, map[string]string{"error": msg})
	}
	return s
}

func (s *Server) Run() error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	return s.echo.Start(addr)
}
