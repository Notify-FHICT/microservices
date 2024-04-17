package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// UpdateNoteID represents an update to link or unlink a note ID
type UpdateNoteID struct {
	ID     primitive.ObjectID `bson:"_id"`
	NoteID primitive.ObjectID `bson:"noteID"`
}

// UpdateTagID represents an update to link or unlink a tag ID
type UpdateTagID struct {
	ID    primitive.ObjectID `bson:"_id"`
	TagID primitive.ObjectID `bson:"tagID"`
}
