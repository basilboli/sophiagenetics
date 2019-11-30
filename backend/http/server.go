package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sg-api/db"
	"strconv"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"

	"github.com/go-chi/chi"
)

// server structure represents contains database client and http router
type server struct {
	Db     *db.DB
	router *chi.Mux
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func NewServer() *server {
	s := &server{}
	s.router = chi.NewRouter()

	// A good base middleware stack goes here
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.Recoverer)

	// Adding CORS middleware
	_cors := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	s.router.Use(_cors.Handler)

	s.routes()
	return s
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, data interface{}, status int) error {
	w.WriteHeader(status)
	if data != nil {
		return json.NewEncoder(w).Encode(data)
	}
	return nil
}

func (s *server) decode(w http.ResponseWriter, r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func (s *server) handleIndex(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Hello SOPHiA GENETICS!"))
}

// handleHealthCheck handles endpoint to be used by external health check system
func (s *server) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	_, _ = io.WriteString(w, `{"alive": true}`)
}

// handleVersion handles endpoint exposing application version
func (s *server) handleVersion(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(fmt.Sprintf("Commit hash: %s\nBuild time: %s", s.Db.Opts.CommitHash, s.Db.Opts.BuildTime)))
	if err != nil {
		log.Printf("[WARN] Problem writing response, %s", err)
	}
}

// handleListClients handles endpoint showing all clients data
func (s *server) handleListClients(w http.ResponseWriter, r *http.Request) {
	polls, err := s.Db.ListClients()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = s.respond(w, r, polls, http.StatusOK)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleCreateClient handles endpoint to create new client
func (s *server) handleCreateClient() http.HandlerFunc {

	type response struct {
		Status string `json:"status"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var clientRequest *db.Client
		err := s.decode(w, r, &clientRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		_, err = s.Db.NewClient(clientRequest)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		greeting := &response{Status: "created"}
		err = s.respond(w, r, greeting, http.StatusCreated)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

}

// handleCreateClient handles endpoint to get client data by saleforceId
func (s *server) handleGetClient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	value := "client"
	client, ok := ctx.Value(value).(*db.Client)
	if !ok {
		log.Println("ERR. Cannot take element from http Context:", value)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err := s.respond(w, r, client, http.StatusOK)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ClientCtx gets and stores client data in context
func (s *server) ClientCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		param := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(param, 10, 64)
		if err != nil {
			http.Error(w, "Wrong client id", 400)
			return
		}
		client, err := s.Db.GetClient(id)
		if err != nil {
			http.Error(w, err.Error(), 404)
			return
		}
		ctx := context.WithValue(r.Context(), "client", client)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
