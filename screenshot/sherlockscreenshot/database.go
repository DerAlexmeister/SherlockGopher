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
	mongoURI = "mongodb://0.0.0.0:27017"
)

type DB struct {
	Client *mongo.Client
}

func NewDB() *DB {
	return &DB{}
}

func (db *DB) GetMongoClient() *mongo.Client {
	return db.Client
}

func (db *DB) SetMongoClient(client *mongo.Client) {
	db.Client = client
}

func Connect() *DB {
	credential := options.Credential{
		Username: "root",
		Password: "rootpassword",
	}

	clientOpts := options.Client().ApplyURI(mongoURI).SetAuth(credential)
	client, err := mongo.NewClient(clientOpts)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	res := NewDB()
	res.SetMongoClient(client)

	return res
}

func (db *DB) Save(input *Screenshot) {
	collection := db.GetMongoClient().Database("dbscreenshots").Collection("screenshots")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, input)
	if err != nil {
		log.Fatal(err)
	}
}

func (db *DB) ReturnAllScreenshots() ([]*Screenshot, error) {
	collection := db.GetMongoClient().Database("dbscreenshots").Collection("screenshots")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	var screenshots []*Screenshot
	for cur.Next(ctx) {
		var screenshot *Screenshot
		err := cur.Decode(&screenshot)
		if err != nil {
			return nil, err
		}
		screenshots = append(screenshots, screenshot)
	}
	return screenshots, nil
}
