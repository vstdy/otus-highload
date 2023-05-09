package rest

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/vstdy/otus-highload/api/rest/model"
	canonical "github.com/vstdy/otus-highload/model"
	"github.com/vstdy/otus-highload/pkg"
	"github.com/vstdy/otus-highload/pkg/logging"
	"github.com/vstdy/otus-highload/service/project"
)

const (
	serviceName = "otus-project server"
)

// Handler keeps handler dependencies.
type Handler struct {
	service  project.IService
	jwtAuth  *jwtauth.JWTAuth
	logLevel zerolog.Level
}

// NewHandler returns a new Handler instance.
func NewHandler(service project.IService, jwtAuth *jwtauth.JWTAuth, logLevel zerolog.Level) Handler {
	return Handler{service: service, jwtAuth: jwtAuth, logLevel: logLevel}
}

// Logger returns logger with service field set.
func (h Handler) Logger(ctx context.Context) (context.Context, zerolog.Logger) {
	ctx, logger := logging.GetCtxLogger(ctx, logging.WithLogLevel(h.logLevel))
	logger = logger.With().Str(logging.ServiceKey, serviceName).Logger()

	return ctx, logger
}

// login authorizes user.
func (h Handler) login(w http.ResponseWriter, r *http.Request) {
	var bodyObj model.LoginBody
	err := json.NewDecoder(r.Body).Decode(&bodyObj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	rawObj, err := bodyObj.ToCanonical()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	obj, err := h.service.AuthenticateUser(r.Context(), rawObj)
	if err != nil {
		if errors.Is(err, pkg.ErrWrongCredentials) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.addJWTCookie(w, obj); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	loginResponse := model.NewLoginResponse(obj)
	res, err := json.Marshal(loginResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// register registers user.
func (h Handler) register(w http.ResponseWriter, r *http.Request) {
	var bodyObj model.RegisterBody
	err := json.NewDecoder(r.Body).Decode(&bodyObj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	rawObj := bodyObj.ToCanonical()

	obj, err := h.service.CreateUser(r.Context(), rawObj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.addJWTCookie(w, obj); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	registerResponse := model.NewRegisterResponse(obj)
	res, err := json.Marshal(registerResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// getUser returns user data.
func (h Handler) getUser(w http.ResponseWriter, r *http.Request) {
	userUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	obj, err := h.service.GetUser(r.Context(), userUUID)
	if err != nil {
		if errors.Is(err, pkg.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	getUserResponse := model.NewGetUserResponse(obj)
	res, err := json.Marshal(getUserResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// searchUser searches users.
func (h Handler) searchUsers(w http.ResponseWriter, r *http.Request) {
	firstName := r.URL.Query().Get("first_name")
	lastName := r.URL.Query().Get("last_name")
	if firstName == "" && lastName == "" {
		http.Error(w, pkg.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}

	searchParams := canonical.SearchUser{FirstName: firstName, LastName: lastName}
	objs, err := h.service.SearchUsers(r.Context(), searchParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	searchUsersResponse := model.NewSearchUsersResponse(objs)
	res, err := json.Marshal(searchUsersResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
