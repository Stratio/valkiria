package routes

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"net/http"
)

type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) *HttpError

type HttpError struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      int    `json:"-"`
}

func NewRouter(ctx context.Context, routes map[string]map[string]Handler) *mux.Router {
	r := mux.NewRouter()

	for method, mappings := range routes {
		for route, fct := range mappings {
			localFct := fct
			wrap := func(w http.ResponseWriter, r *http.Request) {
				log.WithFields(log.Fields{"method": r.Method, "uri": r.RequestURI}).Info("HTTP request received")

				err := localFct(ctx, w, r)
				if err != nil {
					log.WithFields(log.Fields{"method": r.Method, "uri": r.RequestURI}).Info(err.Description)
					w.Header().Set("Content-Type", "application/json; charset=utf-8")
					w.Header().Set("X-Content-Type-Options", "nosniff")
					w.WriteHeader(err.Status)
					enc := json.NewEncoder(w)
					enc.Encode(err)
					return
				}
			}

			r.Path(route).Methods(method).HandlerFunc(wrap)
		}
	}
	return r
}
