package dao

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	rentalpb "server/rental/api/gen/v1"
	"server/shared/id"
	mgo "server/shared/mongo"
	"server/shared/mongo/objid"
)

const (
	tripField      = "trip"
	accountIDField = tripField + ".accountid"
	statusField    = tripField + ".status"
)

type Mongo struct {
	col *mongo.Collection
}

func NewMongo(db *mongo.Database) *Mongo {
	return &Mongo{
		col: db.Collection("trip"),
	}
}

//TripRecord defines a trip record in mongo
type TripRecord struct {
	mgo.IDField        `bson:"inline"`
	mgo.UpdatedAtField `bson:"inline"`
	Trip               *rentalpb.Trip `bson:"trip"`
}

//TODO:同一个account最多只能有应该进行中的trip
//TODO:强类型化tripID
//TODO:表格驱动测试

// CreateTrip creates a trip
func (m *Mongo) CreateTrip(c context.Context, trip *rentalpb.Trip) (*TripRecord, error) {
	r := &TripRecord{
		Trip: trip,
	}
	r.ID = mgo.NewObjectID()
	r.UpdatedAt = mgo.UpdatedAt()
	_, err := m.col.InsertOne(c, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (m *Mongo) GetTrip(c context.Context, id id.TripID, accountID id.AccountID) (*TripRecord, error) {
	objID, err := objid.FromID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id:%v", err)
	}
	res := m.col.FindOne(c, bson.M{
		mgo.IDFieldName: objID,
		accountIDField:  accountID,
	})

	if err := res.Err(); err != nil {
		return nil, err
	}
	var tr TripRecord
	err = res.Decode(&tr)
	if err != nil {
		return nil, fmt.Errorf("cannot decode:%v", err)
	}
	return &tr, nil
}

// GetTrips gets trips for the account by status.
// If status is not specified, gets all trips for the account.
func (m *Mongo) GetTrips(c context.Context, accountID id.AccountID, status rentalpb.TripStatus) ([]*TripRecord, error) {
	filter := bson.M{
		accountIDField: accountID.String(),
	}
	if status != rentalpb.TripStatus_TS_NOT_SPECIFIED {
		filter[statusField] = status
	}

	res, err := m.col.Find(c, filter, options.Find().SetSort(bson.M{
		mgo.IDFieldName: -1,
	}))
	if err != nil {
		return nil, err
	}

	var trips []*TripRecord
	for res.Next(c) {
		var trip TripRecord
		err := res.Decode(&trip)
		if err != nil {
			return nil, err
		}
		trips = append(trips, &trip)
	}
	return trips, nil
}

// UpdateTrip updates a trip.
func (m *Mongo) UpdateTrip(c context.Context, tid id.TripID, aid id.AccountID, updatedAt int64, trip *rentalpb.Trip) error {
	objID, err := objid.FromID(tid)
	if err != nil {
		return fmt.Errorf("invalid id: %v", err)
	}

	newUpdatedAt := mgo.UpdatedAt()
	res, err := m.col.UpdateOne(c, bson.M{
		mgo.IDFieldName:        objID,
		accountIDField:         aid.String(),
		mgo.UpdatedAtFieldName: updatedAt,
	}, mgo.Set(bson.M{
		tripField:              trip,
		mgo.UpdatedAtFieldName: newUpdatedAt,
	}))
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
