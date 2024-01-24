package storage

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"Username"`
	Country  string             `bson:"Country"`
}

func (m *User) UnmarshalJSON(data []byte) error {
	// Define a custom type to unmarshal the raw JSON data
	type Alias User

	// Create an instance of the Alias type to perform the unmarshaling
	aux := &struct {
		ID string `json:"id"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}

	// Unmarshal the raw JSON data into the auxiliary structure
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Convert the string ID to a primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex(aux.ID)
	if err != nil {
		return err
	}

	// Assign the converted ObjectID to the main struct
	m.ID = objectID

	return nil
}
