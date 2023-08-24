package ginx

import (
	"github.com/0xDeSchool/gap/app"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"sort"
)

type serverHandler struct {
	Order int // 优先级
	Func  gin.HandlerFunc
}

type ServerBuilder struct {
	preRuns     []ServerConfigureFunc
	initors     []ServerConfigureFunc
	middlewares []serverHandler

	App     *app.AppBuilder
	Options *ServerOptions
	// 只在程序启动过程中进行操作，以保证协程安全
	Items map[string]any
}

func NewServerBuilder(builder *app.AppBuilder) *ServerBuilder {
	return &ServerBuilder{
		App:     builder,
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

func (b *ServerBuilder) Use(handlers ...gin.HandlerFunc) *ServerBuilder {
	for _, handler := range handlers {
		b.middlewares = append(b.middlewares, serverHandler{Order: 0, Func: handler})
	}
	return b
}

func (b *ServerBuilder) OrderUse(order int, handlers ...gin.HandlerFunc) *ServerBuilder {
	for _, handler := range handlers {
		b.middlewares = append(b.middlewares, serverHandler{Order: order, Func: handler})
	}
	return b
}

func (b *ServerBuilder) Build() (*Server, error) {
	g := gin.Default()
	server := NewServer(g, b.Options)

	sort.SliceStable(b.middlewares, func(i, j int) bool {
		return b.middlewares[i].Order < b.middlewares[j].Order
	})
	for _, handler := range b.middlewares {
		g.Use(handler.Func)
	}

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
