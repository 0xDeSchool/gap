package x

import (
	"github.com/dineshappavoo/basex"
	"github.com/yitter/idgenerator-go/idgen"
	"math/big"
)

type IdGenerator[T comparable] interface {
	Create() T
}

type DefaultIdGenerator struct {
}

func NewDefaultIdGenerator() IdGenerator[string] {
	var options = idgen.NewIdGeneratorOptions(1)
	// 保存参数（务必调用，否则参数设置不生效）：
	idgen.SetIdGenerator(options)
	return &DefaultIdGenerator{}
}

func (u *DefaultIdGenerator) Create() string {
	id := idgen.NextId()
	idStr, err := basex.EncodeInt(big.NewInt(id))
	if err != nil {
		panic(err)
	}
	return idStr
}
