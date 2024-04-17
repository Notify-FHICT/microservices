package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Notify-FHICT/microservices/notes/storage"
	"github.com/Notify-FHICT/microservices/notes/storage/models"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// APIHandler handles HTTP requests for notes.
type APIHandler struct {
	c storage.DB
}

// NewAPIHandler initializes a new APIHandler with a given database collection.
func NewAPIHandler(collection storage.DB) APIHandler {
	return APIHandler{
		collection,
	}
}

// Histogram for tracking request durations.
var histogram = promauto.NewHistogram(prometheus.HistogramOpts{
	Name:    "request_duration_seconds",
	Help:    "Duration of the request.",
	Buckets: []float64{.01, .025, .05, 0.1, 0.15, 0.2, 0.25, 0.3, 0.5},
})

// Counters for tracking processed requests.
var (
	reqProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "processed_req_total",
		Help: "The total number of processed requests",
	})
	reqSuccessfullyProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "successfully_processed_req_total",
		Help: "The total number of successfully processed requests",
	})
)

// Server starts the HTTP server for handling note operations.
func (api *APIHandler) Server() {
	// Prometheus metrics endpoint
	http.Handle("/metrics", promhttp.Handler())

	// Handler for creating a new note
	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		reqProcessed.Inc()
		switch r.Method {
		case http.MethodPost:
			// Create a new record.
			var note models.Note
			err := json.NewDecoder(r.Body).Decode(&note)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}
			api.c.CreateNote(note)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(note.ID)
			reqSuccessfullyProcessed.Inc()
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		histogram.Observe(time.Since(now).Seconds())
	})

	// Handler for reading a note
	http.HandleFunc("/read/", func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		reqProcessed.Inc()
		switch r.Method {
		case http.MethodGet:
			// Serve the resource
			id := strings.TrimPrefix(r.URL.Path, "/read/")
			oid, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}
			note, err := api.c.ReadNote(oid)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(note)
			reqSuccessfullyProcessed.Inc()
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		histogram.Observe(time.Since(now).Seconds())
	})

	// Handler for updating a note
	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		reqProcessed.Inc()
		switch r.Method {
		case http.MethodPut:
			// Update an existing record.
			var note models.Note
			err := json.NewDecoder(r.Body).Decode(&note)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}
			out, err := api.c.UpdateNote(note)
			if err != nil {
				http.Error(w, err.Error(), http.StatusServiceUnavailable)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(out)
			reqSuccessfullyProcessed.Inc()
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		histogram.Observe(time.Since(now).Seconds())
	})

	// Handler for deleting a note
	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		reqProcessed.Inc()
		switch r.Method {
		case http.MethodDelete:
			// Remove the record.
			var note models.Note
			err := json.NewDecoder(r.Body).Decode(&note)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}
			err = api.c.DeleteNote(note.ID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusServiceUnavailable)
				return
			}
			var send models.Middle
			send.NoteID = note.ID
			tmp, err := primitive.ObjectIDFromHex("000000000000000000000000")
			send.ID = tmp
			fmt.Println(send)
			LinkEvent(send)
			if err != nil {
				http.Error(w, err.Error(), http.StatusServiceUnavailable)
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Entry successfully removed"))
			reqSuccessfullyProcessed.Inc()
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		histogram.Observe(time.Since(now).Seconds())
	})

	// Handler for linking a tag ID to a note
	http.HandleFunc("/link_tagID", func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		reqProcessed.Inc()
		switch r.Method {
		case http.MethodPut:
			// Update an existing record.
			var note models.UpdateTagID
			err := json.NewDecoder(r.Body).Decode(&note)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}
			err = api.c.LinkTagID(note)
			if err != nil {
				http.Error(w, err.Error(), http.StatusServiceUnavailable)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Entry successfully modified"))
			reqSuccessfullyProcessed.Inc()
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		histogram.Observe(time.Since(now).Seconds())
	})

	// Handler for linking an event ID to a note
	http.HandleFunc("/link_event", func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		reqProcessed.Inc()
		switch r.Method {
		case http.MethodPut:
			// Update an existing record.
			var note models.UpdateEventID
			err := json.NewDecoder(r.Body).Decode(&note)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}

			var send models.Middle
			send.NoteID = note.ID
			send.ID = note.EventID
			fmt.Println(send)
			LinkEvent(send)
			if err != nil {
				http.Error(w, err.Error(), http.StatusServiceUnavailable)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Operation Queued"))
			reqSuccessfullyProcessed.Inc()
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		histogram.Observe(time.Since(now).Seconds())
	})

	// Handler for updating content of a note
	http.HandleFunc("/update_content", func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		reqProcessed.Inc()
		switch r.Method {
		case http.MethodPut:
			// Update an existing record.
			var note models.Note
			err := json.NewDecoder(r.Body).Decode(&note)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}
			err = api.c.UpdateContent(note)
			if err != nil {
				http.Error(w, err.Error(), http.StatusServiceUnavailable)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Entry successfully modified"))
			reqSuccessfullyProcessed.Inc()
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		histogram.Observe(time.Since(now).Seconds())
	})

	// Initialize a placeholder note and link an event
	null, _ := primitive.ObjectIDFromHex("000000000000000000000000")
	var note models.Middle
	note.NoteID = null
	note.ID = null
	LinkEvent(note)

	// Start HTTP server
	http.ListenAndServe(":3000", nil)

	fmt.Println("Server closed oh no!")
}
