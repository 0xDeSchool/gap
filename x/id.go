package x

import (
	"github.com/0xDeSchool/gap/app"
	"github.com/dineshappavoo/basex"
	"github.com/yitter/idgenerator-go/idgen"
	"math/big"
)

type IdGeneratorOptions struct {
	WorkerId uint16
}

type IdGenerator[T comparable] interface {
	Create() T
}

type StringIdGenerator struct {
}

func NewStringIdGenerator() *IdGenerator[string] {
	workerId := app.Get[IdGeneratorOptions]().WorkerId
	var options = idgen.NewIdGeneratorOptions(workerId)
	// 保存参数（务必调用，否则参数设置不生效）：
	idgen.SetIdGenerator(options)
	var g IdGenerator[string] = &StringIdGenerator{}
	return &g
}

func (u *StringIdGenerator) Create() string {
	id := idgen.NextId()
	idStr, err := basex.EncodeInt(big.NewInt(id))
	if err != nil {
		panic(err)
	}
	return idStr
}

type Int64IdGenerator struct {
}

func NewNumberIdGenerator() *IdGenerator[int64] {
	workerId := app.Get[IdGeneratorOptions]().WorkerId
	var options = idgen.NewIdGeneratorOptions(workerId)
	// 保存参数（务必调用，否则参数设置不生效）：
	idgen.SetIdGenerator(options)
	var g IdGenerator[int64] = &Int64IdGenerator{}
	return &g
}

func (u *Int64IdGenerator) Create() int64 {
	return idgen.NextId()
}

func DefaultIdGenerators(ab *app.AppBuilder) {
	ab.ConfigureServices(func() error {
		app.TryAddValue(&IdGeneratorOptions{WorkerId: 1})
		app.TryAddSingleton(NewStringIdGenerator)
		app.TryAddSingleton(NewNumberIdGenerator)
		return nil
	})
}
