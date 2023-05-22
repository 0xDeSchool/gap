package x

import (
	"github.com/rs/xid"
)

type IdGenerator[T comparable] interface {
	Create() T
}

type UuidGenerator struct {
}

func NewUuidGenerator() IdGenerator[string] {
	return &UuidGenerator{}
}

func (u *UuidGenerator) Create() string {
	return xid.New().String()
}
