package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// UpdateEventID represents the update of an event ID.
type UpdateEventID struct {
	ID      primitive.ObjectID `bson:"_id"`
	EventID primitive.ObjectID `bson:"eventID"`
}

// Middle represents a middle object.
type Middle struct {
	ID     primitive.ObjectID `bson:"_id"`
	NoteID primitive.ObjectID `bson:"noteID"`
}

// UpdateTagID represents the update of a tag ID.
type UpdateTagID struct {
	ID    primitive.ObjectID `bson:"_id"`
	TagID primitive.ObjectID `bson:"tagID"`
}
