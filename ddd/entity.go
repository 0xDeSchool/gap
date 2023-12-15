package ddd

import (
	"encoding/json"
	"strconv"

	"github.com/0xDeSchool/gap/errx"
)

type Entity[TKey comparable] interface {
	GetId() TKey
	SetId(id TKey)
}

type OrderEntity interface {
	GetOrder() float64
}

type OrderEntityBase struct {
	Order float64 `bson:"order"`
}

func (e OrderEntityBase) GetOrder() float64 {
	return e.Order
}

type EntityBase[TKey comparable] struct {
	ID TKey `bson:"_id,omitempty"`
}

func (e EntityBase[TKey]) GetId() TKey {
	return e.ID
}

func (e *EntityBase[TKey]) SetId(id TKey) {
	e.ID = id
}

type WithExtraEntity[TKey comparable] struct {
	Extra map[string]string `bson:"extra,omitempty"`
}

func (m *WithExtraEntity[TKey]) SetValue(key string, v interface{}) {
	if m.Extra == nil {
		m.Extra = make(map[string]string)
	}
	content, err := json.Marshal(v)
	if err == nil {
		m.Extra[key] = string(content)
	}
}

func (m *WithExtraEntity[TKey]) GetValue(key string, v interface{}) (bool, error) {
	if m.Extra == nil {
		return false, nil
	}
	content, ok := m.Extra[key]
	if ok {
		err := json.Unmarshal([]byte(content), v)
		return true, err
	}
	return false, nil
}

func (m *WithExtraEntity[TKey]) GetString(key string) string {
	if m.Extra == nil {
		return ""
	}
	content, ok := m.Extra[key]
	if ok {
		return content
	}
	return ""
}

func (m *WithExtraEntity[TKey]) SetString(key string, value string) {
	if m.Extra == nil {
		m.Extra = make(map[string]string)
	}
	m.Extra[key] = value
}

func (m *WithExtraEntity[TKey]) GetInt(key string) int {
	if m.Extra == nil {
		return 0
	}
	content, ok := m.Extra[key]
	if ok {
		v, err := strconv.Atoi(content)
		errx.CheckError(err)
		return v
	}
	return 0
}

func (m *WithExtraEntity[TKey]) SetInt(key string, value int) {
	if m.Extra == nil {
		m.Extra = make(map[string]string)
	}
	m.Extra[key] = strconv.Itoa(value)
}

func (m *WithExtraEntity[TKey]) GetFloat32(key string) float32 {
	if m.Extra == nil {
		return 0
	}
	content, ok := m.Extra[key]
	if ok {
		v, err := strconv.ParseFloat(content, 32)
		errx.CheckError(err)
		return float32(v)
	}
	return 0
}

func (m *WithExtraEntity[TKey]) SetFloat32(key string, value float32) {
	if m.Extra == nil {
		m.Extra = make(map[string]string)
	}
	m.Extra[key] = strconv.FormatFloat(float64(value), 'E', -1, 32)
}

type WithExtraArrayEntity[TKey comparable] struct {
	Data []string `bson:"data,omitempty"`
}

func (a *WithExtraArrayEntity[TKey]) Append(str string) {
	if a.Data == nil {
		a.Data = make([]string, 0)
	}
	a.Data = append(a.Data, str)
}
