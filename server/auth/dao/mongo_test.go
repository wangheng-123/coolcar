package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	mgo "server/shared/mongo"
	mongotesting "server/shared/mongo/testing"
	"testing"
)

var mongoURI string

func TestResolveAccountID(t *testing.T) {

	c := context.Background()
	mc, err := mongo.Connect(c, options.Client().ApplyURI(mongoURI))
	if err != nil {
		t.Fatalf("cannot connect mongodb:%v", err)
	}
	m := NewMongo(mc.Database("test"))

	_, err = m.col.InsertMany(c, []interface{}{
		bson.M{
			mgo.IDField: mustObjectID("643d35ef330e000039003a39"),
			openIDField: "openid_1",
		},
		bson.M{
			mgo.IDField: mustObjectID("643d35ef330e000039003a40"),
			openIDField: "openid_2",
		},
	})
	if err != nil {
		t.Fatalf("cannot insert initial values: %v", err)
	}

	m.newObjID = func() primitive.ObjectID {
		return mustObjectID("643d35ef330e000039003a41")
	}

	cases := []struct {
		name   string
		openID string
		want   string
	}{
		{
			name:   "existing_user",
			openID: "openid_1",
			want:   "643d35ef330e000039003a39",
		},
		{
			name:   "another_existing_user",
			openID: "openid_2",
			want:   "643d35ef330e000039003a40",
		},
		{
			name:   "new_user",
			openID: "openid_3",
			want:   "643d35ef330e000039003a41",
		},
	}

	for _, cc := range cases {
		t.Run(cc.name, func(t *testing.T) {
			id, err := m.ResolveAccountID(context.Background(), cc.openID)
			if err != nil {
				t.Errorf("failed resolve account id for %q:%v", cc.openID, err)
			}
			if id != cc.want {
				t.Errorf("resolve account id:want:%q,got:%q", cc.want, id)
			}
		})
	}
}

func mustObjectID(hex string) primitive.ObjectID {
	objID, err := primitive.ObjectIDFromHex(hex)
	if err != nil {
		panic(err)
	}
	return objID
}

func TestMain(m *testing.M) {
	os.Exit(mongotesting.RunWithMongoInDocker(m, &mongoURI))
}
