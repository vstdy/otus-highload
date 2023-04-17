package rest

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwa"

	"github.com/vstdy/otus-highload/service/project"
)

// NewRouter returns router.
func NewRouter(svc project.Service, config Config) (chi.Router, error) {
	jwtAuth := jwtauth.New(jwa.HS256.String(), []byte(config.SecretKey), nil)
	h := NewHandler(svc, jwtAuth, config.LogLevel)
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Use(
			middleware.RequestID,
			middleware.RealIP,
			middleware.Logger,
			middleware.Recoverer,
			middleware.StripSlashes,
			middleware.Timeout(config.Timeout),
			middleware.AllowContentType("application/json"),
		)

		// Public routes
		r.Group(func(r chi.Router) {
			r.Post("/login", h.login)
			r.Post("/user/register", h.register)
		})

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(
				jwtauth.Verifier(jwtAuth),
				jwtauth.Authenticator,
			)

			r.Route("/user", func(r chi.Router) {
				r.Get("/get/{id}", h.getUser)
			})
		})
	})

	return r, nil
}
