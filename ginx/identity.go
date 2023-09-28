package ginx

import (
	"context"
	"github.com/0xDeSchool/gap/multi_tenancy"
)

type userKeyType struct{}

var userKey = userKeyType{}

type CurrentUserInfo[TKey comparable] struct {
	ID       TKey
	UserName string
	Address  string
	Avatar   string
	AuthType string // 验证类型
	Data     map[string]any
}

func (u *CurrentUserInfo[TKey]) Authenticated() bool {
	var dv TKey
	return u.ID != dv
}

func (u *CurrentUserInfo[TKey]) Get(key string) any {
	if u.Data == nil {
		return nil
	}
	return u.Data[key]
}

func (u *CurrentUserInfo[TKey]) Set(key string, value any) {
	if u.Data == nil {
		u.Data = make(map[string]any)
	}
	u.Data[key] = value
}

func CurrentUser[TKey comparable](c context.Context) *CurrentUserInfo[TKey] {
	user, ok := c.Value(userKey).(*CurrentUserInfo[TKey])
	if !ok {
		return &CurrentUserInfo[TKey]{}
	}
	return user
}

func WithUserContext[TKey comparable](c context.Context, user *CurrentUserInfo[TKey]) context.Context {
	return context.WithValue(c, userKey, user)
}

func CurrentTenant(c context.Context) *multi_tenancy.TenantInfo {
	t := multi_tenancy.CurrentTenant(c)
	if t == nil {
		t = &multi_tenancy.TenantInfo{}
	}
	return t
}
