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

type serverHandler struct {
	Func gin.HandlerFunc
}

type Server struct {
	G           *gin.Engine
	Options     *ServerOptions
	middlewares []serverHandler
}

func NewServer(g *gin.Engine, options *ServerOptions) *Server {
	return &Server{
		G:       g,
		Options: options,
	}
}

func (s *Server) Use(middleware gin.HandlerFunc) {
	s.middlewares = append(s.middlewares, serverHandler{
		Func: middleware,
	})
}
func (s *Server) Run() error {
	s.useDefaultHandlers()
	for _, m := range s.middlewares {
		s.G.Use(m.Func)
	}
	addr := ":" + strconv.Itoa(s.Options.Port)
	if s.Options.Port == 0 {
		addr = ":5000"
	}
	log.Info("********** Listening: " + addr + " ***********\n")
	return s.G.Run(addr)
}

func (s *Server) useDefaultHandlers() {
	s.G.ContextWithFallback = true

	mid := logger.SetLogger(
		logger.WithDefaultLevel(s.Options.LogLevel),
	)
	s.Use(mid)
	s.Use(ErrorMiddleware)
	s.Use(UnitWorkMiddleware())
}
