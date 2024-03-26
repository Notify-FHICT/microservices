package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Notify-FHICT/microservices/notes/storage"
	"github.com/Notify-FHICT/microservices/notes/storage/models"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type APIHandler struct {
	c storage.DB
}

func NewAPIHandler(collection storage.DB) APIHandler {
	return APIHandler{
		collection,
	}
}

var (
	demoGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "Gauge_IncDec",
		Help: "testing module",
	})
)

func (api *APIHandler) Server() {
	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			// Create a new record.
			var note models.Note
			err := json.NewDecoder(r.Body).Decode(&note)
			// err := note.UnmarshalJSON([]byte(r.Body.Close().Error()))
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			}
			api.c.CreateNote(note)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(note.ID)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/read/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// Serve the resource
			id := strings.TrimPrefix(r.URL.Path, "/read/")
			oid, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			}
			note, err := api.c.ReadNote(oid)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(note)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			// Update an existing record.
			var note models.Note
			err := json.NewDecoder(r.Body).Decode(&note)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			}
			out, err := api.c.UpdateNote(note)
			if err != nil {
				http.Error(w, err.Error(), http.StatusServiceUnavailable)
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(out)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			// Remove the record.
			var note models.Note
			err := json.NewDecoder(r.Body).Decode(&note)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			}
			err = api.c.DeleteNote(note.ID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusServiceUnavailable)
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Entry successfully removed"))
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/link_tagID", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			// Update an existing record.
			var note models.UpdateTagID
			err := json.NewDecoder(r.Body).Decode(&note)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			}
			err = api.c.LinkTagID(note)
			if err != nil {
				http.Error(w, err.Error(), http.StatusServiceUnavailable)
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Entry successfully modified"))
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/link_event", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			// Update an existing record.
			var note models.UpdateEventID
			err := json.NewDecoder(r.Body).Decode(&note)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			}

			var send models.Middle
			send.NoteID = note.ID
			send.ID = note.EventID
			fmt.Println(send)
			LinkEvent(send)
			if err != nil {
				http.Error(w, err.Error(), http.StatusServiceUnavailable)
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Operation Queued"))
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/update_content", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			// Update an existing record.
			var note models.Note
			err := json.NewDecoder(r.Body).Decode(&note)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			}
			err = api.c.UpdateContent(note)
			if err != nil {
				http.Error(w, err.Error(), http.StatusServiceUnavailable)
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Entry successfully modified"))
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	null, _ := primitive.ObjectIDFromHex("000000000000000000000000")
	var note models.Middle
	note.NoteID = null
	note.ID = null
	LinkEvent(note)

	http.ListenAndServe(":3000", nil)

	fmt.Println("Server closed oh no!")

}
