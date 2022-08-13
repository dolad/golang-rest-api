package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dolad/rest-api/internal/comment"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Handler struct {
	Router  *mux.Router
	Service *comment.Service
}

type Response struct {
	Message string
	Error   string
}

func NewHandler(service *comment.Service) *Handler {
	return &Handler{
		Service: service,
	}
}

func BasicAuth(original func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if user == "admin" && pass == "password" && ok {
			original(w, r)
		} else {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			sendErrorResponse(w, "not authorized", errors.New("not authorized"))
		}
	}
}

// LoggingMiddleware - a handy middleware function that logs out incoming requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(
			log.Fields{
				"Method": r.Method,
				"Path":   r.URL.Path,
			}).
			Info("handled request")

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) SetupRoutes() {
	log.Info("Setting up Routes")
	h.Router = mux.NewRouter()
	h.Router.Use(LoggingMiddleware)
	h.Router.HandleFunc("/api/comment/{id}", h.GetComment).Methods("GET")
	h.Router.HandleFunc("/api/comment/{id}", BasicAuth(h.UpdateComment)).Methods("PUT")
	h.Router.HandleFunc("/api/comment/{id}", BasicAuth(h.DeleteComment)).Methods("DELETE")
	h.Router.HandleFunc("/api/comment", BasicAuth(h.PostComment)).Methods("POST")
	h.Router.HandleFunc("/api/comment", h.GetAllComments).Methods("GET")
	h.Router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(Response{Message: "I am alive"}); err != nil {
			panic(err)
		}
	})
}

func sendErrorResponse(w http.ResponseWriter, message string, err error) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusInternalServerError)
	if err := json.NewEncoder(w).Encode(Response{Message: message, Error: err.Error()}); err != nil {
		panic(err)
	}
}

func sendOkResponse(w http.ResponseWriter, resp interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(resp)
}
