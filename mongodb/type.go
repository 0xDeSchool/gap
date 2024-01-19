package mongodb

import (
	"github.com/0xDeSchool/gap/errx"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func IDFromHex(id string) primitive.ObjectID {
	if id == "" {
		return primitive.NilObjectID
	}
	result, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		errx.PanicValidatition(err.Error())
	}
	return result
}

func IDString(id primitive.ObjectID) string {
	if id.IsZero() {
		return ""
	}
	return id.Hex()
}

func MustToId(id *string) primitive.ObjectID {
	if id == nil || *id == "" {
		errx.PanicValidatition("string id must not be empty")
	}
	result, err := primitive.ObjectIDFromHex(*id)
	if err != nil {
		errx.PanicValidatition(err.Error())
	}
	return result
}

func ToID(id *string) primitive.ObjectID {
	if id == nil || *id == "" {
		return primitive.NilObjectID
	}
	result, err := primitive.ObjectIDFromHex(*id)
	if err != nil {
		errx.PanicValidatition(err.Error())
	}
	return result
}
