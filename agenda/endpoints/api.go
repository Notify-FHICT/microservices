package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Notify-FHICT/microservices/agenda/storage"
	"github.com/Notify-FHICT/microservices/agenda/storage/models"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// APIHandler handles incoming HTTP requests for agenda management
type APIHandler struct {
	c storage.DB
}

// NewAPIHandler creates a new APIHandler with the provided collection
func NewAPIHandler(collection storage.DB) APIHandler {
	return APIHandler{
		collection,
	}
}

// Define Prometheus metrics
var histogram = promauto.NewHistogram(prometheus.HistogramOpts{
	Name:    "request_duration_seconds",
	Help:    "Duration of the request.",
	Buckets: []float64{.01, .05, 0.1, 0.15, 0.2, 0.25, 0.3, 0.4, 0.5},
})

var (
	reqProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "processed_req_total",
		Help: "The total number of processed requests",
	})
)
var (
	reqSuccessfullyProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "successfully_processed_req_total",
		Help: "The total number of successfully processed requests",
	})
)

// Server starts the HTTP server for handling agenda management requests
func (api *APIHandler) Server() {
	// Expose Prometheus metrics endpoint
	http.Handle("/metrics", promhttp.Handler())

	// Handle "create" endpoint
	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		reqProcessed.Inc()
		switch r.Method {
		case http.MethodPost:
			// Create a new record.
			var event models.Event
			err := json.NewDecoder(r.Body).Decode(&event)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}
			api.c.CreateEvent(event)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(event.ID)
			reqSuccessfullyProcessed.Inc()
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		histogram.Observe(time.Since(now).Seconds())
	})

	// Handle "read" endpoint
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
			event, err := api.c.ReadEvent(oid)
			if err != nil {
				http.Error(w, err.Error(), http.StatusServiceUnavailable)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(event)
			reqSuccessfullyProcessed.Inc()
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		histogram.Observe(time.Since(now).Seconds())
	})

	// Handle "update" endpoint
	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		reqProcessed.Inc()
		switch r.Method {
		case http.MethodPut:
			// Update an existing record.
			var event models.Event
			err := json.NewDecoder(r.Body).Decode(&event)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}
			out, err := api.c.UpdateEvent(event)
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

	// Handle "link_noteID" endpoint
	http.HandleFunc("/link_noteID", func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		reqProcessed.Inc()
		switch r.Method {
		case http.MethodPut:
			// Update an existing record.
			var event models.UpdateNoteID
			err := json.NewDecoder(r.Body).Decode(&event)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}
			err = api.c.LinkNoteID(event)
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

	// Handle "link_tagID" endpoint
	http.HandleFunc("/link_tagID", func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		reqProcessed.Inc()
		switch r.Method {
		case http.MethodPut:
			// Update an existing record.
			var event models.UpdateTagID
			err := json.NewDecoder(r.Body).Decode(&event)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}
			err = api.c.LinkTagID(event)
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

	// Handle "delete" endpoint
	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		reqProcessed.Inc()
		switch r.Method {
		case http.MethodDelete:
			// Remove the record.
			var event models.Event
			err := json.NewDecoder(r.Body).Decode(&event)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}
			err = api.c.DeleteEvent(event.ID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusServiceUnavailable)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Entry successfully removed"))
			reqSuccessfullyProcessed.Inc()
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		histogram.Observe(time.Since(now).Seconds())
	})

	// Handle "dashboard" endpoint
	http.HandleFunc("/dashboard/", func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		reqProcessed.Inc()
		switch r.Method {
		case http.MethodGet:
			// Serve the resource
			id := strings.TrimPrefix(r.URL.Path, "/dashboard/")
			oid, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}
			out, err := api.c.ReadDashboard(oid)
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

	// Start HTTP server on port 3000
	http.ListenAndServe(":3000", nil)

	// Print a message indicating server closure
	fmt.Println("Server closed oh no!")
}
