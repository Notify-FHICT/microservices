package endpoints

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Notify-FHICT/microservices/user/service"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type APIHandler struct {
	s *service.Service
}

func NewAPIHandler(service *service.Service) APIHandler {
	return APIHandler{
		s: service,
	}
}

var (
	demoGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "testing trying to autoscale!!!",
		Help: "testing module",
	})
)

func (api *APIHandler) Server() {

	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, User!")
		fmt.Fprintf(w, "API path: %s", r.RequestURI)
	})

	http.HandleFunc("/plus", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "added 1")
		demoGauge.Inc()

		fmt.Fprintf(w, "API path: %s", r.RequestURI)
	})

	http.HandleFunc("/minus", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "removed 1")
		demoGauge.Dec()

		fmt.Fprintf(w, "API path: %s", r.RequestURI)
	})

	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(strings.ToLower(r.Header.Get("Content-Type")), "json") {
			err := api.s.CreateUser(w, r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			} else {
				fmt.Println(w, "Creation successful")
			}
		}

	})

	http.HandleFunc("/read", func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(strings.ToLower(r.Header.Get("Content-Type")), "json") {
			usr, err := api.s.ReadUser(w, r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			fmt.Fprintf(w, "Succesfully Read; %s", usr)
		}
	})

	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(strings.ToLower(r.Header.Get("Content-Type")), "json") {
			usr, err := api.s.UpdateUser(w, r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			fmt.Fprintf(w, "Succesfully Updated to; %s", usr)
		}
	})

	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(strings.ToLower(r.Header.Get("Content-Type")), "json") {
			err := api.s.DeleteUser(w, r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			} else {
				fmt.Println(w, "Succesfully deleted")
			}
		}
	})

	http.ListenAndServe(":3000", nil)

	fmt.Println("I died")

}
