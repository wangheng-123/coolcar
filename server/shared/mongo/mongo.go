package mgo

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"server/shared/mongo/objid"
	"time"
)

//common field name
const (
	IDFieldName        = "_id"
	UpdatedAtFieldName = "updatedat"
)

type IDField struct {
	ID primitive.ObjectID `bson:"_id"`
}

// UpdatedAtField defines the updatedat field
type UpdatedAtField struct {
	UpdatedAt int64 `bson:"updatedat"`
}

// NewObjectID generates a new object id.
var NewObjectID = primitive.NewObjectID

// NewObjIDWithValue sets id for next objectID generation.
func NewObjIDWithValue(id fmt.Stringer) {
	NewObjectID = func() primitive.ObjectID {
		return objid.MustFromID(id)
	}
}

// UpdatedAt returns a value suitable for UpdatedAt field
var UpdatedAt = func() int64 {
	return time.Now().UnixNano()
}

func Set(v interface{}) bson.M {
	return bson.M{
		"$set": v,
	}
}

func SetOnInsert(v interface{}) bson.M {
	return bson.M{
		"$setOnInsert": v,
	}
}

// ZeroOrDoesNotExist generates a filter expression with
// field equal to zero or field does not exist.
func ZeroOrDoesNotExist(field string, zero interface{}) bson.M {
	return bson.M{
		"$or": []bson.M{
			{
				field: zero,
			},
			{
				field: bson.M{
					"$exists": false,
				},
			},
		},
	}
}
