package connctordb

import (
	"context"
	"fmt"
	models "rest_api_server/internal/Models"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Save(item models.Item) (result *mongo.InsertOneResult) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		panic(err)
	}
	c := client.Database("user").Collection("item")
	result, err = c.InsertOne(context.TODO(), item)
	return result
}

func Get(id int) models.Item {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		panic(err)
	}

	fmt.Print("TEST")
	c := client.Database("user").Collection("item")
	cursor, err := c.Find(context.TODO(), id)
	if err != nil {
		panic(err)
	}
	defer cursor.Close(context.TODO())
	var item models.Item
	if cursor.Next(context.TODO()) {
		err := cursor.Decode(&item)
		if err != nil {
			panic(err)
		}
	}
	return item
}
