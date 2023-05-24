package dao

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"server/shared/id"
	mgo "server/shared/mongo"
	"server/shared/mongo/objid"
)

const openIDField = "open_id"

type Mongo struct {
	col *mongo.Collection
}

func NewMongo(db *mongo.Database) *Mongo {
	return &Mongo{
		col: db.Collection("account"),
	}
}

// ResolveAccountID resolves the account id from open id.
func (m *Mongo) ResolveAccountID(c context.Context, openID string) (id.AccountID, error) {
	//m.col.InsertOne(c,bson.M{
	//	mgo.IDField: m.newObjID(),
	//	openIDField: openID,
	//})

	insertedID := mgo.NewObjectID()

	res := m.col.FindOneAndUpdate(c, bson.M{
		openIDField: openID,
	}, mgo.SetOnInsert(bson.M{
		mgo.IDFieldName: insertedID,
		openIDField:     openID,
	}), options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After))
	if err := res.Err(); err != nil {
		return "", fmt.Errorf("cannot findOneAndUpdate:%v", err)
	}
	var row mgo.IDField
	err := res.Decode(&row)
	if err != nil {
		return "", fmt.Errorf("cannot decode result:%v", err)
	}
	return objid.ToAccountID(row.ID), nil
}
