package ginx

import (
	"context"
	"github.com/0xDeSchool/gap/multi_tenancy"
)

type CurrentUserInfo struct {
	ID       string
	UserName string
	Address  string
	Avatar   string
	Data     map[string]any
}

func (u *CurrentUserInfo) Authenticated() bool {
	return u.ID != ""
}

func (u *CurrentUserInfo) Get(key string) any {
	if u.Data == nil {
		return nil
	}
	return u.Data[key]
}

func (u *CurrentUserInfo) Set(key string, value any) {
	if u.Data == nil {
		u.Data = make(map[string]any)
	}
	u.Data[key] = value
}

func CurrentUser(c context.Context) *CurrentUserInfo {
	user, ok := c.Value("Login.User").(*CurrentUserInfo)
	if !ok {
		return &CurrentUserInfo{
			ID: "",
		}
	}
	return user
}

func CurrentTenant(c context.Context) *multi_tenancy.TenantInfo {
	t := multi_tenancy.CurrentTenant(c)
	if t == nil {
		t = &multi_tenancy.TenantInfo{}
	}
	return t
}
