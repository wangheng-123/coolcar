package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	c := context.Background()
	connect, err := mongo.Connect(c, options.Client().ApplyURI("mongodb://root:root@localhost:27017"))
	if err != nil {
		panic(err)
	}
	collection := connect.Database("test").Collection("account")

	findRows(c, collection)

}

func findRows(c context.Context, col *mongo.Collection) {
	//res := col.FindOne(c, bson.M{
	//	"open_id": "123",
	//})
	//fmt.Printf("%+v\n", res)
	//var row struct {
	//	ID     primitive.ObjectID `bson:"_id"`
	//	OpenID string             `bson:"open_id"`
	//}
	//err := res.Decode(&row)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("%+v\n", row)

	cur, err := col.Find(c, bson.M{})
	if err != nil {
		panic(err)
	}
	for cur.Next(c) {
		var row struct {
			ID     primitive.ObjectID `bson:"_id"`
			OpenID string             `bson:"open_id"`
		}
		err = cur.Decode(&row)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%+v\n", row)
	}
}

func insertRows(c context.Context, col *mongo.Collection) {
	result, err := col.InsertMany(c, []interface{}{
		bson.M{
			"open_id": "125",
		},
		bson.M{
			"open_id": "458",
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", result)
}
