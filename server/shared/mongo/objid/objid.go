package objid

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"server/shared/id"
)

// FromID converts an id to object id
func FromID(id fmt.Stringer) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id.String())
}

func MustFromID(id fmt.Stringer) primitive.ObjectID {
	oid, err := FromID(id)
	if err != nil {
		fmt.Println("error")
		panic(err)
	}
	return oid
}

// ToAccountID converts object id to account id
func ToAccountID(oid primitive.ObjectID) id.AccountID {
	return id.AccountID(oid.Hex())
}

// ToTripID converts object id to trip id
func ToTripID(oid primitive.ObjectID) id.TripID {
	return id.TripID(oid.Hex())
}
