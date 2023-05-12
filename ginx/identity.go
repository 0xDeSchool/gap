package ginx

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CurrentUserInfo struct {
	ID       primitive.ObjectID
	UserName string
	Address  string
	Avatar   string
	Data     map[string]any
}

func (u *CurrentUserInfo) Authenticated() bool {
	return !u.ID.IsZero()
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
			ID: primitive.NilObjectID,
		}
	}
	return user
}
