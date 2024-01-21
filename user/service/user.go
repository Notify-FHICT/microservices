package service

import (
	"encoding/json"
	"net/http"

	"github.com/Notify-FHICT/microservices/service/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	c storage.DB
}

func NewUserService(collection storage.DB) *Service {

	return &Service{
		collection,
	}

}

func decodeIDs(content *http.Request) (*primitive.ObjectID, error) {
	var id primitive.ObjectID
	err := json.NewDecoder(content.Body).Decode(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func decodeUsers(content *http.Request) (*storage.User, error) {
	var usr storage.User
	err := json.NewDecoder(content.Body).Decode(&usr)
	if err != nil {
		return nil, err
	}
	return &usr, nil
}

func (s *Service) CreateUser(w http.ResponseWriter, r *http.Request) error {
	usr, err := decodeUsers(r)
	if err != nil {
		return err
	}
	return s.c.CreateUser(*usr)
}

func (s *Service) ReadUser(w http.ResponseWriter, r *http.Request) (*storage.User, error) {
	id, err := decodeIDs(r)
	if err != nil {
		return nil, err
	}
	return s.c.ReadUser(*id)
} //(id primitive.ObjectID) {

func (s *Service) UpdateUser(w http.ResponseWriter, r *http.Request) (*storage.User, error) {
	usr, err := decodeUsers(r)
	if err != nil {
		return nil, err
	}

	// result, err :=
	return s.c.UpdateUser(*usr)
	// log.Print(result)
	// return err
}

func (s *Service) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	id, err := decodeIDs(r)
	if err != nil {
		return err
	}
	return s.c.DeleteUser(*id)
	// log.Print(usr)
	// return err
}
