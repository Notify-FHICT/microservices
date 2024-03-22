package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type UpdateNoteID struct {
	ID     primitive.ObjectID `bson:"_id"`
	NoteID primitive.ObjectID `bson:"noteID"`
}

type UpdateTagID struct {
	ID    primitive.ObjectID `bson:"_id"`
	TagID primitive.ObjectID `bson:"tagID"`
}
