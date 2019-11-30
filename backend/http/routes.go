package http

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (s *server) routes() {
	s.router.Get("/version", s.handleVersion)
	s.router.Get("/hc", s.handleHealthCheck)

	s.router.Get("/", s.handleIndex)

	// RESTy routes for "clients" resource
	s.router.Route("/clients", func(r chi.Router) {
		r.With(paginate).Get("/", s.handleListClients)
		r.Post("/", s.handleCreateClient()) // POST /clients

		// Subrouters:
		r.Route("/{id}", func(r chi.Router) {
			r.Use(s.ClientCtx)
			r.Get("/", s.handleGetClient) // GET /client/123
		})
	})
}

// paginate is a stub, but very possible to implement middleware logic
// to handle the request params for handling a paginated request
func paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// just a stub.. some ideas are to look at URL query params for something like
		// the page number, or the limit, and send a query cursor down the chain
		next.ServeHTTP(w, r)
	})
}
