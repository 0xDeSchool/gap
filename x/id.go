package x

import (
	"github.com/0xDeSchool/gap/app"
	"github.com/dineshappavoo/basex"
	"github.com/yitter/idgenerator-go/idgen"
	"math/big"
)

type ID = int64
type StringID = string

type IdGeneratorOptions struct {
	WorkerId uint16
}

type IdGenerator[T comparable] interface {
	Create() T
}

type StringIdGenerator[T string] struct {
}

func NewStringIdGenerator[T string]() *IdGenerator[T] {
	workerId := app.Get[IdGeneratorOptions]().WorkerId
	var options = idgen.NewIdGeneratorOptions(workerId)
	// 保存参数（务必调用，否则参数设置不生效）：
	idgen.SetIdGenerator(options)
	var g IdGenerator[T] = &StringIdGenerator[T]{}
	return &g
}

func (u *StringIdGenerator[T]) Create() T {
	id := idgen.NextId()
	idStr, err := basex.EncodeInt(big.NewInt(id))
	if err != nil {
		panic(err)
	}
	return T(idStr)
}

type Int64IdGenerator[T int64] struct {
}

func NewNumberIdGenerator[T int64]() *IdGenerator[T] {
	workerId := app.Get[IdGeneratorOptions]().WorkerId
	var options = idgen.NewIdGeneratorOptions(workerId)
	// 保存参数（务必调用，否则参数设置不生效）：
	idgen.SetIdGenerator(options)
	var g IdGenerator[T] = &Int64IdGenerator[T]{}
	return &g
}

func (u *Int64IdGenerator[T]) Create() T {
	return T(idgen.NextId())
}

func DefaultIdGenerators(ab *app.AppBuilder) {
	ab.ConfigureServices(func() error {
		app.TryAddValue(&IdGeneratorOptions{WorkerId: 1})
		app.TryAddSingleton(NewStringIdGenerator[string])
		app.TryAddSingleton(NewNumberIdGenerator[int64])
		return nil
	})
}
