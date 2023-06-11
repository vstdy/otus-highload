package rest

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/vstdy/otus-highload/api/rest/handler"
	"github.com/vstdy/otus-highload/service/project"
)

// NewRouter returns router.
func NewRouter(svc project.IService, config Config) (chi.Router, error) {
	jwtAuth := jwtauth.New(jwa.HS256.String(), []byte(config.SecretKey), nil)
	h := handler.NewHandler(svc, jwtAuth, config.LogLevel)
	r := chi.NewRouter()

	r.Handle("/metrics", promhttp.Handler())

	r.Route("/", func(r chi.Router) {
		r.Use(
			middleware.RequestID,
			middleware.RealIP,
			middleware.Logger,
			middleware.Recoverer,
			middleware.StripSlashes,
			middleware.Timeout(config.Timeout),
			middleware.AllowContentType("application/json"),
			addMetrics(),
		)

		// Public routes
		r.Group(func(r chi.Router) {
			r.Post("/login", h.Login)
			r.Post("/user/register", h.Register)
		})

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(
				jwtauth.Verifier(jwtAuth),
				jwtauth.Authenticator,
			)

			r.Route("/user", func(r chi.Router) {
				r.Get("/get/{id}", h.GetUser)
				r.Get("/search", h.SearchUsers)
			})

			r.Route("/friend", func(r chi.Router) {
				r.Put("/set/{id}", h.SetFriend)
				r.Put("/delete/{id}", h.DeleteFriend)
			})

			r.Route("/post", func(r chi.Router) {
				r.Post("/create", h.CreatePost)
				r.Put("/update", h.UpdatePost)
				r.Put("/delete/{id}", h.DeletePost)
				r.Get("/get/{id}", h.GetPost)
				r.Get("/feed", h.GetFeed)
			})

			r.Route("/dialog", func(r chi.Router) {
				r.Post("/{id}/send", h.SendDialog)
				r.Get("/{id}/list", h.ListDialog)
			})
		})
	})

	return r, nil
}
