package eventbus

import (
	"context"

	"github.com/0xDeSchool/gap/log"
	"github.com/0xDeSchool/gap/x"
	"github.com/lileio/pubsub/v2"
	"github.com/lileio/pubsub/v2/middleware/recover"
)

const testError = "test publish error"

func InitEventBus(serviceName string) {
	pubsub.SetClient(&pubsub.Client{
		ServiceName: serviceName,
		Provider:    NewMemoryProvider(),
		Middleware: []pubsub.Middleware{
			recover.Middleware{},
		},
	})
	go ListenErrorEvents(0)
}

func Publish[T any](ctx context.Context, obj *T) {
	topic := x.TypeOf[T]().String()
	PublishJSON(ctx, topic, obj)
}

func PublishJSON(ctx context.Context, topic string, obj interface{}) {
	result := pubsub.PublishJSON(ctx, topic, obj)
	<-result.Ready
	if result.Err != nil {
		log.Warnf("publish message failed, topic: %s, error: %s", topic, result.Err.Error())
	}
}

func Subscribe[T any](handler func(ctx context.Context, msg *T) error) {
	if handler == nil {
		panic("handler is nil")
	}
	c := pubsub.GetClient()
	c.On(NewHandlerOption(handler))
}
