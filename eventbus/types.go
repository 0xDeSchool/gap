package eventbus

type EntityCreatedEvent[T any] struct {
	Entity *T `json:"entity"`
}

func Created[T any](entity *T) *EntityCreatedEvent[T] {
	return &EntityCreatedEvent[T]{
		Entity: entity,
	}
}

type EntityUpdatedEvent[T any] struct {
	Entity *T `json:"entity"`
}

func Updated[T any](entity *T) *EntityUpdatedEvent[T] {
	return &EntityUpdatedEvent[T]{
		Entity: entity,
	}
}

type EntityDeletedEvent[T any] struct {
	ID     string `json:"id"`
	Entity *T     `json:"entity"`
}

func Deleted[T any](entity *T) *EntityDeletedEvent[T] {
	return &EntityDeletedEvent[T]{
		Entity: entity,
	}
}
