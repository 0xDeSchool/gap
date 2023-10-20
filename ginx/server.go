package ginx

import (
	"github.com/gin-contrib/logger"
	"strconv"

	"github.com/0xDeSchool/gap/log"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type ServerConfigureFunc func(s *Server) error

type ServerOptions struct {
	Port     int
	RootUrl  string
	LogLevel zerolog.Level
}

type Server struct {
	G       *gin.Engine
	Options *ServerOptions
}

func NewServer(g *gin.Engine, options *ServerOptions) *Server {
	return &Server{
		G:       g,
		Options: options,
	}
}

func (s *Server) Run() error {
	addr := ":" + strconv.Itoa(s.Options.Port)
	if s.Options.Port == 0 {
		addr = ":5000"
	}
	log.Info("********** Listening: " + addr + " ***********\n")
	return s.G.Run(addr)
}

func (s *Server) Use(handlers ...gin.HandlerFunc) *Server {
	s.G.Use(handlers...)
	return s
}

func (s *Server) UseDefaultHandlers() *Server {
	s.G.ContextWithFallback = true

	mid := logger.SetLogger(
		logger.WithDefaultLevel(s.Options.LogLevel),
	)
	s.G.Use(mid)
	s.G.Use(ErrorMiddleware)
	s.G.Use(UnitWorkMiddleware())
	return s
}
