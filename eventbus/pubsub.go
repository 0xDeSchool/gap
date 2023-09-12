package eventbus

import (
	"context"
	"strings"
	"time"

	"github.com/0xDeSchool/gap/x"
	"github.com/lileio/pubsub/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type handlerOption[T any] struct {
	topic      string
	MsgHandler func(ctx context.Context, msg *T) error
}

func NewHandlerOption[T any](isSync bool, handler func(ctx context.Context, msg *T) error) pubsub.HandlerOptions {
	op := handlerOption[T]{
		MsgHandler: handler,
	}
	op.topic = x.TypeOf[T]().String()
	opts := pubsub.HandlerOptions{
		Topic:   x.TypeOf[T]().String(),
		Name:    x.TypeOf[T]().String(),
		Handler: op.handle,
		AutoAck: true,
		JSON:    true,
	}
	if isSync {
		opts.Name = "__sync__" + opts.Name
	}
	return opts
}

func (s *handlerOption[T]) handle(ctx context.Context, msg *T, pm *pubsub.Msg) error {
	err := s.MsgHandler(ctx, msg)
	return err
}

type msgHandler struct {
	IsSync  bool
	Handler pubsub.MsgHandler
}
type MemoryProvider struct {
	handlers     map[string][]msgHandler
	ErrorHandler func(ctx context.Context, msg *pubsub.Msg, err error)
}

func NewMemoryProvider() *MemoryProvider {
	return &MemoryProvider{
		handlers: make(map[string][]msgHandler),
	}
}

func (mp *MemoryProvider) Publish(ctx context.Context, topic string, m *pubsub.Msg) error {
	mp.handle(ctx, topic, m)
	return nil
}

func (mp *MemoryProvider) handle(ctx context.Context, topic string, m *pubsub.Msg) {
	if m.PublishTime == nil {
		m.PublishTime = x.Ptr(time.Now())
	}
	if m.ID == "" {
		m.ID = primitive.NewObjectID().Hex()
	}
	for _, h := range mp.handlers[topic] {
		if h.IsSync {
			mp.internalHandle(ctx, h.Handler, m)
		} else {
			go mp.internalHandle(ctx, h.Handler, m)
		}
	}
}

func (mp *MemoryProvider) internalHandle(ctx context.Context, mh pubsub.MsgHandler, m *pubsub.Msg) {
	err := mh(ctx, *m)
	if err != nil {
		if mp.ErrorHandler != nil {
			mp.ErrorHandler(ctx, m, err)
		}
	}
}

func (mp *MemoryProvider) Subscribe(opts pubsub.HandlerOptions, h pubsub.MsgHandler) {
	mh := msgHandler{
		IsSync:  strings.HasPrefix(opts.Name, "__sync__"),
		Handler: h,
	}
	mp.handlers[opts.Topic] = append(mp.handlers[opts.Topic], mh)
}

func (mp *MemoryProvider) Shutdown() {
	return
}
