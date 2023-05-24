package trip

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	rentalpb "server/rental/api/gen/v1"
	"server/rental/trip/client/poi"
	"server/rental/trip/dao"
	"server/shared/auth"
	"server/shared/id"
	mgo "server/shared/mongo"
	mongotesting "server/shared/mongo/testing"
	"server/shared/server"
	"testing"
)

func TestCreateTrip(t *testing.T) {
	c := auth.ContextWithAccountID(context.Background(), "account1")
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannot create mongo client:%v", err)
	}

	logger, err := server.NewZapLogger()
	if err != nil {
		t.Fatalf("cannot create logger:%v", err)
	}

	pm := &profileManager{}
	cm := &carManager{}
	//s := newService(c, t, pm, cm)

	s := &Service{
		ProfileManager: pm,
		CarManager:     cm,
		POIManager:     &poi.Manager{},
		Mongo:          dao.NewMongo(mc.Database("coolcar")),
		Logger:         logger,
	}

	//nowFunc = func() int64 {
	//	return 1605695246
	//}
	req := &rentalpb.CreateTripRequest{
		CarId: "car1",
		Start: &rentalpb.Location{
			Latitude:  32.123,
			Longitude: 114.2525,
		},
	}
	pm.iID = "identity1"
	golden := `{"account_id":"%q","car_id":"car1","start":{"location":{"latitude":32.123,"longitude":114.2525},"poi_name":"天安门","timestamp_sec":1605695246},"current":{"location":{"latitude":32.123,"longitude":114.2525},"poi_name":"天安门","timestamp_sec":1605695246},"status":1,"identity_id":"identity1"}`
	cases := []struct {
		name string
		//accountID    string
		tripID       string
		profileErr   error
		carVerifyErr error
		carUnlockErr error
		want         string
		wantErr      bool
	}{
		{
			name: "normal_create",
			//accountID: "account1",
			tripID: "5f8132eb12714bf629489054",
			//want:      fmt.Sprintf(golden, "account1"),
			want: golden, //??
		},
		{
			name: "profile_err",
			//accountID:  "account2",
			tripID:     "5f8132eb12714bf629489055",
			profileErr: fmt.Errorf("profile"),
			wantErr:    true,
		},
		{
			name: "car_verify_err",
			//accountID:    "account3",
			tripID:       "5f8132eb12714bf629489056",
			carVerifyErr: fmt.Errorf("verify"),
			wantErr:      true,
		},
		{
			name: "car_unlock_err",
			//accountID:    "account4",
			tripID:       "5f8132eb12714bf629489057",
			carUnlockErr: fmt.Errorf("unlock"),
			//want:         fmt.Sprintf(golden, "account4"),
			want: golden,
		},
	}

	for _, cc := range cases {
		t.Run(cc.name, func(t *testing.T) {
			mgo.NewObjIDWithValue(id.TripID(cc.tripID))
			pm.err = cc.profileErr
			cm.unlockErr = cc.carUnlockErr
			cm.verifyErr = cc.carVerifyErr
			//c := auth.ContextWithAccountID(context.Background(), id.AccountID(cc.accountID))
			//res, err := s.CreateTrip(c, req)

			res, err := s.CreateTrip(c, req)
			if cc.wantErr {
				if err == nil {
					t.Errorf("want error; got none")
				} else {
					return
				}
			}
			if err != nil {
				t.Errorf("error creating trip: %v", err)
				return
			}
			if res.Id != cc.tripID {
				t.Errorf("incorrect id; want %q, got %q", cc.tripID, res.Id)
			}
			b, err := json.Marshal(res.Trip)
			if err != nil {
				t.Errorf("cannot marshall response: %v", err)
			}
			got := string(b)
			if cc.want != got {
				t.Errorf("incorrect response: want %s, got %s", cc.want, got)
			}
		})
	}

}

type profileManager struct {
	iID id.IdentityID
	err error
}

func (p *profileManager) Verify(context.Context, id.AccountID) (id.IdentityID, error) {
	return p.iID, p.err
}

type carManager struct {
	verifyErr error
	unlockErr error
}

func (m *carManager) Verify(context.Context, id.CarID, *rentalpb.Location) error {
	return m.verifyErr
}

func (m *carManager) Unlock(c context.Context, cid id.CarID, aid id.AccountID, tid id.TripID, avatarURL string) error {
	return m.unlockErr
}

func (m *carManager) Lock(c context.Context, cid id.CarID) error {
	return nil
}
func TestMain(m *testing.M) {
	os.Exit(mongotesting.RunWithMongoInDocker(m))
}
