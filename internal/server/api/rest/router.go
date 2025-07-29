package resthandler

import (
	"github.com/go-chi/chi/v5"
)

// Route - setting service routes
func Route(s *Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/api", func(r chi.Router) {
		r.Get("/", s.UserListHandler)
		r.Route("/user", func(r chi.Router) {
			r.Post("/block", s.UserBlock)
			r.Post("/unblock", s.UserUnblock)
		})
		r.Route("/token", func(r chi.Router) {
			r.Get("/{username}", s.TokenList)
			r.Post("/block", s.TokenBlock)
		})
	})
	return r
}
