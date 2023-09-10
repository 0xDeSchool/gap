package ginx

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type ServerBuilder struct {
	preRuns []ServerConfigureFunc
	initors []ServerConfigureFunc

	Options *ServerOptions
	// 只在程序启动过程中进行操作，以保证协程安全
	Items map[string]any
}

func NewServerBuilder() *ServerBuilder {
	return &ServerBuilder{
		preRuns: make([]ServerConfigureFunc, 0),
		initors: make([]ServerConfigureFunc, 0),
		Options: &ServerOptions{
			LogLevel: zerolog.InfoLevel,
		},
		Items: make(map[string]any),
	}
}

// PreConfigure 配置服务，该方法参数在App.Run中、gin.Run之前运行
func (b *ServerBuilder) PreConfigure(action ServerConfigureFunc) *ServerBuilder {
	b.preRuns = append(b.preRuns, action)
	return b
}

// Configure 配置服务，该方法参数在App.Run中、gin.Run之前运行
func (b *ServerBuilder) Configure(action ServerConfigureFunc) *ServerBuilder {
	b.initors = append(b.initors, action)
	return b
}

func (b *ServerBuilder) Add(fun func(b *ServerBuilder)) *ServerBuilder {
	fun(b)
	return b
}

func (b *ServerBuilder) Build() (*Server, error) {
	g := gin.Default()
	server := NewServer(g, b.Options)

	for _, action := range b.preRuns {
		err := action(server)
		if err != nil {
			return nil, err
		}
	}

	for _, action := range b.initors {
		err := action(server)
		if err != nil {
			return nil, err
		}
	}

	return server, nil
}
