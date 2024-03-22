package models

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Note struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	UserID  primitive.ObjectID `bson:"userID"`
	TagID   primitive.ObjectID `bson:"tagID,omitempty"`
	Title   string             `bson:"title"`
	Content string             `bson:"content"`
}

func (m *Note) UnmarshalJSON(data []byte) error {
	// Define a custom type to unmarshal the raw JSON data
	type Alias Note

	// Create an instance of the Alias type to perform the unmarshaling
	aux := &struct {
		ID     string `json:"id,omitempty"`
		UserID string `json:"userid"`
		TagID  string `json:"tagid,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}

	// Unmarshal the raw JSON data into the auxiliary structure
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	null, err := primitive.ObjectIDFromHex("000000000000000000000000")
	if err != nil {
		return err
	}

	// Convert the string ID to a primitive.ObjectID
	objNoteID, err := primitive.ObjectIDFromHex(aux.ID)
	if err != nil {
		objNoteID = primitive.NewObjectID()
	}
	objNoteUserID, err := primitive.ObjectIDFromHex(aux.UserID)
	if err != nil {
		return err
	}
	objNoteTagID, err := primitive.ObjectIDFromHex(aux.TagID)
	if err != nil {
		objNoteTagID = null
	}

	// Assign the converted ObjectID to the main struct
	m.ID = objNoteID
	m.UserID = objNoteUserID
	m.TagID = objNoteTagID

	return nil
}
