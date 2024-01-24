package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/Notify-FHICT/microservices/user/service/storage"
)

type Service struct {
	c storage.DB
}

func NewUserService(collection storage.DB) *Service {
	return &Service{
		collection,
	}
}

func decodeIDs(r *http.Request) (*storage.User, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var usr storage.User
	if err := usr.UnmarshalJSON(body); err != nil {
		return nil, err
	}

	return &usr, nil
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

func (s *Service) ReadUser(w http.ResponseWriter, r *http.Request) ([]byte, error) {

	usr, err := decodeIDs(r)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return nil, err
	}
	fmt.Println(reflect.TypeOf(usr))
	usr2, err2 := s.c.ReadUser(usr.ID)
	if err != nil {
		http.Error(w, "Database Error", http.StatusBadRequest)
		return nil, err2
	}
	return json.Marshal(usr2)

} //(id primitive.ObjectID) {

func (s *Service) UpdateUser(w http.ResponseWriter, r *http.Request) (*storage.User, error) {
	usr, err := decodeIDs(r)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return nil, err
	}

	// usr, err := decodeUsers(r)
	// if err != nil {
	// 	return nil, err
	// }

	return s.c.UpdateUser(*usr)
}

func (s *Service) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	usr, err := decodeIDs(r)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return err
	}
	// id, err := decodeIDs(r)
	// if err != nil {
	// 	return err
	// }
	return s.c.DeleteUser(usr.ID)
}
