package x

import (
	"math/rand"
	"time"
)

const keys = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var rd = rand.New(rand.NewSource(time.Now().Unix()))

func Intn(n int) int {
	return rd.Intn(n)
}

func Float64() float64 {
	return rd.Float64()
}

func Float32() float32 {
	return rd.Float32()
}

func Int63() int64 {
	return rd.Int63()
}

func Letters(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = keys[rd.Intn(len(keys))]
	}
	return string(b)
}
