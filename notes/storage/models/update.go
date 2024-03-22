package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type UpdateEventID struct {
	ID      primitive.ObjectID `bson:"_id"`
	EventID primitive.ObjectID `bson:"eventID"`
}

type UpdateTagID struct {
	ID    primitive.ObjectID `bson:"_id"`
	TagID primitive.ObjectID `bson:"tagID"`
}
