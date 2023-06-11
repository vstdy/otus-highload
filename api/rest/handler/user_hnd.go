package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/vstdy/otus-highload/api/rest/model"
	"github.com/vstdy/otus-highload/pkg"
)

// Login authorizes user.
func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
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

// Register registers user.
func (h Handler) Register(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

// GetUser returns user data.
func (h Handler) GetUser(w http.ResponseWriter, r *http.Request) {
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

// SearchUsers searches users.
func (h Handler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	firstName := r.URL.Query().Get("first_name")
	lastName := r.URL.Query().Get("last_name")
	if firstName == "" && lastName == "" {
		http.Error(w, pkg.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}

	objs, err := h.service.SearchUsers(r.Context(), firstName, lastName)
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
