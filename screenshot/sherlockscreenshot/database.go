package sherlockscreenshot

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoURI = "mongodb://localhost:27017"
)

type DB struct {
	Client *mongo.Client
}

func Connect() *DB {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return &DB{
		Client: client,
	}
}

func (db *DB) Save(input *Screenshot) {
	collection := db.Client.Database("dbscreenshots").Collection("screenshots")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, input)
	if err != nil {
		log.Fatal(err)
	}
}

func (db *DB) ReturnAllScreenshots() []*Screenshot {
	collection := db.Client.Database("dbscreenshots").Collection("screenshots")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	var screenshots []*Screenshot
	for cur.Next(ctx) {
		var screenshot *Screenshot
		err := cur.Decode(&screenshot)
		if err != nil {
			log.Fatal(err)
		}
		screenshots = append(screenshots, screenshot)
	}
	return screenshots
}
