package eventbus

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/0xDeSchool/gap/app"
	"sync"
	"time"

	"github.com/0xDeSchool/gap/errx"
	"github.com/0xDeSchool/gap/log"
	"github.com/0xDeSchool/gap/x"
	"github.com/lileio/pubsub/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type handlerOption[T any] struct {
	topic      string
	MsgHandler func(ctx context.Context, msg *T) error
}

func NewHandlerOption[T any](handler func(ctx context.Context, msg *T) error) pubsub.HandlerOptions {
	op := handlerOption[T]{
		MsgHandler: handler,
	}
	op.topic = x.TypeOf[T]().String()
	return pubsub.HandlerOptions{
		Topic:   x.TypeOf[T]().String(),
		Name:    x.TypeOf[T]().String(),
		Handler: op.handle,
		AutoAck: true,
		JSON:    true,
	}
}

func (s *handlerOption[T]) handle(ctx context.Context, msg *T, pm *pubsub.Msg) error {
	err := s.MsgHandler(ctx, msg)
	if err != nil {
		if store, ok := app.GetOptional[EventStore](); ok {
			msgStr, _ := json.Marshal(msg)
			saveErr := (*store).Save(ctx, &EventMsg{
				Topic:       s.topic,
				Data:        string(msgStr),
				EventID:     pm.ID,
				PublishTime: pm.PublishTime,
				Metadata:    pm.Metadata,
				ErrorMsg:    err.Error(),
			})
			if saveErr != nil {
				log.Warn("save event failed", saveErr)
			}
		}
	}
	return err
}

type MemoryProvider struct {
	mutex    sync.RWMutex
	handlers map[string][]pubsub.MsgHandler
}

func NewMemoryProvider() *MemoryProvider {
	return &MemoryProvider{
		handlers: make(map[string][]pubsub.MsgHandler),
	}
}

func (mp *MemoryProvider) Publish(ctx context.Context, topic string, m *pubsub.Msg) error {
	mp.mutex.Lock()
	defer mp.mutex.Unlock()
	errs := make([]error, 0)
	if m.PublishTime == nil {
		m.PublishTime = x.Ptr(time.Now())
	}
	if m.ID == "" {
		m.ID = primitive.NewObjectID().Hex()
	}
	for _, h := range mp.handlers[topic] {
		err := h(ctx, *m)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errx.Errors(fmt.Sprintf("handle memory MQ message occurs %d error", len(errs)), errs...)
	}
	return nil
}

func (mp *MemoryProvider) Subscribe(opts pubsub.HandlerOptions, h pubsub.MsgHandler) {
	mp.mutex.RLock()
	defer mp.mutex.RUnlock()

	mp.handlers[opts.Topic] = append(mp.handlers[opts.Topic], h)
}

func (mp *MemoryProvider) Shutdown() {
	return
}
