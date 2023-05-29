package ddd

import (
	"encoding/json"
	"strconv"

	"github.com/0xDeSchool/gap/errx"
)

type Entity interface {
	GetId() string
	SetId(id string)
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

type EntityBase struct {
	ID string `bson:"_id,omitempty"`
}

func (e *EntityBase) GetId() string {
	return e.ID
}

func (e *EntityBase) SetId(id string) {
	e.ID = id
}

type WithExtraEntity struct {
	Extra map[string]string `bson:"extra,omitempty"`
}

func (m *WithExtraEntity) SetValue(key string, v interface{}) {
	if m.Extra == nil {
		m.Extra = make(map[string]string)
	}
	content, err := json.Marshal(v)
	if err == nil {
		m.Extra[key] = string(content)
	}
}

func (m *WithExtraEntity) GetValue(key string, v interface{}) (bool, error) {
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

func (m *WithExtraEntity) GetString(key string) string {
	if m.Extra == nil {
		return ""
	}
	content, ok := m.Extra[key]
	if ok {
		return content
	}
	return ""
}

func (m *WithExtraEntity) SetString(key string, value string) {
	if m.Extra == nil {
		m.Extra = make(map[string]string)
	}
	m.Extra[key] = value
}

func (m *WithExtraEntity) GetInt(key string) int {
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

func (m *WithExtraEntity) SetInt(key string, value int) {
	if m.Extra == nil {
		m.Extra = make(map[string]string)
	}
	m.Extra[key] = strconv.Itoa(value)
}

func (m *WithExtraEntity) GetFloat32(key string) float32 {
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

func (m *WithExtraEntity) SetFloat32(key string, value float32) {
	if m.Extra == nil {
		m.Extra = make(map[string]string)
	}
	m.Extra[key] = strconv.FormatFloat(float64(value), 'E', -1, 32)
}

type WithExtraArrayEntity struct {
	Data []string `bson:"data,omitempty"`
}

func (a *WithExtraArrayEntity) Append(str string) {
	if a.Data == nil {
		a.Data = make([]string, 0)
	}
	a.Data = append(a.Data, str)
}
