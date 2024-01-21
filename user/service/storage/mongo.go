package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB interface {
	//CRUD:
	CreateUser(user User) error
	ReadUser(id primitive.ObjectID) (*User, error)
	UpdateUser(user User) (*User, error)
	DeleteUser(id primitive.ObjectID) error
}

type MongoDB struct {
	collection *mongo.Collection
}

func NewMongoDB() (DB, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://admin:adm1n@noteagenda.5tgrti6.mongodb.net/?retryWrites=true&w=majority").SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	collection := client.Database("Account").Collection("users")

	err = Ping(client)
	if err != nil {
		return nil, err
	}

	return &MongoDB{
		collection,
	}, nil
}

func Pong(t time.Time) {
	log.Printf("Pong! (took: %s)", time.Since(t))
}

func Ping(client *mongo.Client) error {

	defer Pong(time.Now())

	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		return err
	}

	return nil
}

func (db *MongoDB) CreateUser(user User) error {
	result, err := db.collection.InsertOne(context.TODO(), user)
	fmt.Printf("%s got pushed", result)
	return err
}

func (db MongoDB) ReadUser(id primitive.ObjectID) (*User, error) {
	filter := bson.D{{Key: "_id", Value: id}}
	var result User
	err := db.collection.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		return nil, err
	}
	return &result, nil

}
func (db MongoDB) UpdateUser(user User) (*User, error) {
	filter := bson.D{{Key: "_id", Value: user.ID}}
	var result User
	err := db.collection.FindOneAndReplace(context.TODO(), filter, user).Decode(&result)

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (db MongoDB) DeleteUser(id primitive.ObjectID) error {
	filter := bson.D{{Key: "_id", Value: id}}
	var result User
	err := db.collection.FindOneAndDelete(context.TODO(), filter).Decode(&result)
	if err != nil {
		return err
	}
	log.Printf("%s has been deleted", result)
	return err
}

// func insertone() {

// 	entry := Tag{UserID: 4, Title: "work", Color: "blue"}

// 	result, err := collection.InsertOne(context.TODO(), entry)
// 	if err != nil {
// 		panic(err)
// 	}user
// 	fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)
// }
