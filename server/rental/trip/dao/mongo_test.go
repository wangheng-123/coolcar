package dao

import (
	"context"
	"github.com/google/go-cmp/cmp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/testing/protocmp"
	"os"
	rentalpb "server/rental/api/gen/v1"
	"server/shared/id"
	mgo "server/shared/mongo"
	"server/shared/mongo/objid"
	mongotesting "server/shared/mongo/testing"
	"testing"
)

func TestCreateTrip(t *testing.T) {
	c := context.Background()
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannot connect mongodb:%v", err)
	}
	db := mc.Database("test")
	err = mongotesting.SetupIndexes(c, db)
	if err != nil {
		t.Fatalf("cannot setup indexes:%v", err)
	}
	m := NewMongo(db)

	cases := []struct {
		name       string
		tripID     string
		accountID  string
		tripStatus rentalpb.TripStatus
		wantErr    bool
	}{
		{
			name:       "finished",
			tripID:     "644a190d372ba4ffb380ce5c",
			accountID:  "account1",
			tripStatus: rentalpb.TripStatus_FINISHED,
			wantErr:    false,
		},
		{
			name:       "another_finished",
			tripID:     "6449f94c772b58303041294c",
			accountID:  "account1",
			tripStatus: rentalpb.TripStatus_FINISHED,
			wantErr:    false,
		},
		{
			name:       "in_progress",
			tripID:     "6449f94c772b58303041294e",
			accountID:  "account1",
			tripStatus: rentalpb.TripStatus_IN_PROGRESS,
			wantErr:    false,
		},
		{
			name:       "another_in_progress",
			tripID:     "6449f94c772b58303041294f",
			accountID:  "account1",
			tripStatus: rentalpb.TripStatus_IN_PROGRESS,
			wantErr:    true,
		},
		{
			name:       "in_progress_by_another_account",
			tripID:     "6449f94c772b58303041294g",
			accountID:  "account2",
			tripStatus: rentalpb.TripStatus_IN_PROGRESS,
			wantErr:    false,
		},
	}

	for _, cc := range cases {
		mgo.NewObjectID = func() primitive.ObjectID {
			return objid.MustFromID(id.TripID(cc.tripID))
		}

		tr, err := m.CreateTrip(c, &rentalpb.Trip{
			AccountId: cc.accountID,
			Status:    cc.tripStatus,
		})

		if cc.wantErr {
			if err == nil {
				t.Errorf("%s :error expected; got none", cc.name)
			}
			continue
		}

		if err != nil {
			t.Errorf("%s: error createTrip:%v", cc.name, err)
			continue
		}

		if tr.ID.Hex() != cc.tripID {
			t.Errorf("%s: incorrect trip id;want:%q; got:%q", cc.name, cc.tripID, tr.ID.Hex())
		}
	}
}

func TestGetTrip(t *testing.T) {
	c := context.Background()
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannot connect mongodb:%v", err)
	}
	m := NewMongo(mc.Database("test"))

	acct := id.AccountID("account1")
	mgo.NewObjectID = primitive.NewObjectID
	trip, err := m.CreateTrip(c, &rentalpb.Trip{
		AccountId: "account1",
		CarId:     "car1",
		Start: &rentalpb.LocationStatus{
			Location: &rentalpb.Location{
				Latitude:  30,
				Longitude: 120,
			},
			FeeCent:  0,
			KmDriven: 0,
			PoiName:  "start_point",
		},
		End: &rentalpb.LocationStatus{
			Location: &rentalpb.Location{
				Latitude:  35,
				Longitude: 115,
			},
			FeeCent:  100,
			KmDriven: 35,
			PoiName:  "endpoint",
		},
		Status: rentalpb.TripStatus_IN_PROGRESS,
	})
	if err != nil {
		t.Fatalf("cannot create trip:%v", err)
	}

	got, err := m.GetTrip(c, objid.ToTripID(trip.ID), acct)
	if err != nil {
		t.Errorf("cannot get trip:%v", err)
	}

	//got.Trip.Start.PoiName = "badstart"

	if diff := cmp.Diff(trip, got, protocmp.Transform()); diff != "" {
		t.Errorf("result differs; -want +got:%s", diff)
	}

}

func TestGetTrips(t *testing.T) {
	rows := []struct {
		id        string
		accountID string
		status    rentalpb.TripStatus
	}{
		{
			id:        "5f8132eb10714bf629489051",
			accountID: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_FINISHED,
		},
		{
			id:        "5f8132eb10714bf629489052",
			accountID: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_FINISHED,
		},
		{
			id:        "5f8132eb10714bf629489053",
			accountID: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_FINISHED,
		},
		{
			id:        "5f8132eb10714bf629489054",
			accountID: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_IN_PROGRESS,
		},
		{
			id:        "5f8132eb10714bf629489055",
			accountID: "account_id_for_get_trips_1",
			status:    rentalpb.TripStatus_IN_PROGRESS,
		},
	}

	c := context.Background()
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannot connect mongodb: %v", err)
	}

	m := NewMongo(mc.Database("coolcar"))

	for _, r := range rows {
		mgo.NewObjIDWithValue(id.TripID(r.id))
		_, err := m.CreateTrip(c, &rentalpb.Trip{
			AccountId: r.accountID,
			Status:    r.status,
		})
		if err != nil {
			t.Fatalf("cannot create rows: %v", err)
		}
	}

	cases := []struct {
		name       string
		accountID  string
		status     rentalpb.TripStatus
		wantCount  int
		wantOnlyID string
	}{
		{
			name:      "get_all",
			accountID: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_TS_NOT_SPECIFIED,
			wantCount: 4,
		},
		{
			name:       "get_in_progress",
			accountID:  "account_id_for_get_trips",
			status:     rentalpb.TripStatus_IN_PROGRESS,
			wantCount:  1,
			wantOnlyID: "5f8132eb10714bf629489054",
		},
	}

	for _, cc := range cases {
		t.Run(cc.name, func(t *testing.T) {
			res, err := m.GetTrips(context.Background(),
				id.AccountID(cc.accountID),
				cc.status)
			if err != nil {
				t.Errorf("cannot get trips: %v", err)
			}

			if cc.wantCount != len(res) {
				t.Errorf("incorrect result count; want: %d, got: %d",
					cc.wantCount, len(res))
			}

			if cc.wantOnlyID != "" && len(res) > 0 {
				if cc.wantOnlyID != res[0].ID.Hex() {
					t.Errorf("only_id incorrect; want: %q, got %q",
						cc.wantOnlyID, res[0].ID.Hex())
				}
			}
		})
	}
}

func TestUpdateTrip(t *testing.T) {
	c := context.Background()
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannot connect mongodb: %v", err)
	}

	m := NewMongo(mc.Database("coolcar"))
	tid := id.TripID("5f8132eb12714bf629489054")
	aid := id.AccountID("account_for_update")

	var now int64 = 10000
	mgo.NewObjIDWithValue(tid)
	mgo.UpdatedAt = func() int64 {
		return now
	}

	tr, err := m.CreateTrip(c, &rentalpb.Trip{
		AccountId: aid.String(),
		Status:    rentalpb.TripStatus_IN_PROGRESS,
		Start: &rentalpb.LocationStatus{
			PoiName: "start_poi",
		},
	})
	if err != nil {
		t.Fatalf("cannot create trip: %v", err)
	}
	if tr.UpdatedAt != 10000 {
		t.Fatalf("wrong updatedat; want: 10000, got: %d", tr.UpdatedAt)
	}

	update := &rentalpb.Trip{
		AccountId: aid.String(),
		Status:    rentalpb.TripStatus_IN_PROGRESS,
		Start: &rentalpb.LocationStatus{
			PoiName: "start_poi_updated",
		},
	}
	cases := []struct {
		name          string
		now           int64
		withUpdatedAt int64
		wantErr       bool
	}{
		{
			name:          "normal_update",
			now:           20000,
			withUpdatedAt: 10000,
		},
		{
			name:          "update_with_stale_timestamp",
			now:           30000,
			withUpdatedAt: 10000,
			wantErr:       true,
		},
		{
			name:          "update_with_refetch",
			now:           40000,
			withUpdatedAt: 20000,
		},
	}

	for _, cc := range cases {
		now = cc.now
		err := m.UpdateTrip(c, tid, aid, cc.withUpdatedAt, update)
		if cc.wantErr {
			if err == nil {
				t.Errorf("%s: want error; got none", cc.name)
			} else {
				continue
			}
		} else {
			if err != nil {
				t.Errorf("%s: cannot update: %v", cc.name, err)
			}
		}
		updatedTrip, err := m.GetTrip(c, tid, aid)
		if err != nil {
			t.Errorf("%s: cannot get trip after udpate: %v", cc.name, err)
		}
		if cc.now != updatedTrip.UpdatedAt {
			t.Errorf("%s: incorrect updatedat: want %d, got %d",
				cc.name, cc.now, updatedTrip.UpdatedAt)
		}
	}
}

func TestMain(m *testing.M) {
	os.Exit(mongotesting.RunWithMongoInDocker(m))
}
