package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Notify-FHICT/microservices/agenda/storage"
	"github.com/Notify-FHICT/microservices/agenda/storage/models"

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
			var event models.Event
			err := json.NewDecoder(r.Body).Decode(&event)
			// err := event.UnmarshalJSON([]byte(r.Body.Close().Error()))
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			}
			api.c.CreateEvent(event)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(event.ID)
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
			event, err := api.c.ReadEvent(oid)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusFound)
			json.NewEncoder(w).Encode(event)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			// Update an existing record.
			var event models.Event
			err := json.NewDecoder(r.Body).Decode(&event)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			}
			out, err := api.c.UpdateEvent(event)
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

	http.HandleFunc("/link_noteID", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			// Update an existing record.
			var event models.UpdateNoteID
			err := json.NewDecoder(r.Body).Decode(&event)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			}
			err = api.c.LinkNoteID(event)
			if err != nil {
				http.Error(w, err.Error(), http.StatusServiceUnavailable)
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Entry successfully modified"))
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/link_tagID", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			// Update an existing record.
			var event models.UpdateTagID
			err := json.NewDecoder(r.Body).Decode(&event)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			}
			err = api.c.LinkTagID(event)
			if err != nil {
				http.Error(w, err.Error(), http.StatusServiceUnavailable)
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Entry successfully modified"))
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			// Remove the record.
			var event models.Event
			err := json.NewDecoder(r.Body).Decode(&event)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			}
			err = api.c.DeleteEvent(event.ID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusServiceUnavailable)
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Entry successfully removed"))

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/dashboard/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// Serve the resource
			id := strings.TrimPrefix(r.URL.Path, "/dashboard/")
			oid, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			}
			out, err := api.c.ReadDashboard(oid)
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

	http.ListenAndServe(":3000", nil)

	fmt.Println("Server closed oh no!")

}
