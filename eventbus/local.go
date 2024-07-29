package eventbus

import (
	"context"
	"github.com/0xDeSchool/gap/app"
)

type LocalEventHandler[T any] interface {
	HandleEvent(ctx context.Context, eventData T) error
}

func PublishLocal[T any](ctx context.Context, eventData T) error {
	handlers := app.GetArray[LocalEventHandler[T]]()
	for _, h := range handlers {
		err := h.HandleEvent(ctx, eventData)
		if err != nil {
			return err
		}
	}
	return nil
}
