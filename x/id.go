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
	WorkerId          uint16
	WorkerIdBitLength byte // 机器码位长，默认值6，取值范围 [1, 15]（要求：序列数位长+机器码位长不超过22）
	SeqBitLength      byte // 序列数位长，默认值6，取值范围 [3, 21]（要求：序列数位长+机器码位长不超过22）
}

func (io *IdGeneratorOptions) toOptions() *idgen.IdGeneratorOptions {
	opts := idgen.NewIdGeneratorOptions(1)
	if io.WorkerId != 0 {
		opts.WorkerId = io.WorkerId
	}
	if io.WorkerIdBitLength != 0 {
		opts.WorkerIdBitLength = io.WorkerIdBitLength
	}
	if io.SeqBitLength != 0 {
		opts.SeqBitLength = io.SeqBitLength
	}
	return opts
}

type IdGenerator[T comparable] interface {
	Create() T
}

type StringIdGenerator[T ~string] struct {
}

func NewStringIdGenerator[T ~string]() *IdGenerator[T] {
	options := app.Get[IdGeneratorOptions]().toOptions()
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

type Int64IdGenerator[T ~int64] struct {
}

func NewNumberIdGenerator[T ~int64]() *IdGenerator[T] {
	options := app.Get[IdGeneratorOptions]().toOptions()
	// 保存参数（务必调用，否则参数设置不生效）：
	idgen.SetIdGenerator(options)
	var g IdGenerator[T] = &Int64IdGenerator[T]{}
	return &g
}

func (u *Int64IdGenerator[T]) Create() T {
	return T(idgen.NextId())
}
