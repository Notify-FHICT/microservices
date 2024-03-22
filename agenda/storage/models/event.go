package models

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	UserID primitive.ObjectID `bson:"userID"`
	TagID  primitive.ObjectID `bson:"tagID,omitempty"`
	NoteID primitive.ObjectID `bson:"noteID,omitempty"`
	Time   primitive.DateTime `bson:"time"`
	Title  string             `bson:"title"`
}

func (m *Event) UnmarshalJSON(data []byte) error {
	// Define a custom type to unmarshal the raw JSON data
	type Alias Event

	// Create an instance of the Alias type to perform the unmarshaling
	aux := &struct {
		ID     string `json:"id,omitempty"`
		UserID string `json:"userid"`
		TagID  string `json:"tagid,omitempty"`
		NoteID string `json:"noteid,omitempty"`
		Time   string `json:"time"`
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
	objEventID, err := primitive.ObjectIDFromHex(aux.ID)
	if err != nil {
		objEventID = primitive.NewObjectID()
	}
	objEventUserID, err := primitive.ObjectIDFromHex(aux.UserID)
	if err != nil {
		return err
	}
	objEventTagID, err := primitive.ObjectIDFromHex(aux.TagID)
	if err != nil {
		objEventTagID = null
	}

	objEventNoteID, err := primitive.ObjectIDFromHex(aux.NoteID)
	if err != nil {
		objEventNoteID = null
	}
	datetime, err := time.Parse("02/01/2006 15:04:05", aux.Time)
	if err != nil {
		return err
	}
	objEventTime := primitive.NewDateTimeFromTime(datetime)

	// Assign the converted ObjectID to the main struct
	m.ID = objEventID
	m.UserID = objEventUserID
	m.TagID = objEventTagID
	m.NoteID = objEventNoteID
	m.Time = objEventTime

	return nil
}
