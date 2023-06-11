package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vstdy/otus-highload/api/rest/model"
	"github.com/vstdy/otus-highload/pkg"
)

// SendDialog sends message to dialog.
func (h Handler) SendDialog(w http.ResponseWriter, r *http.Request) {
	from, to, err := h.getUUIDs(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var dialog model.SendDialogBody
	err = json.NewDecoder(r.Body).Decode(&dialog)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = h.service.SendDialog(r.Context(), from, to, dialog.Text)
	if err != nil {
		if errors.As(err, new(pkg.ErrInvalidUserArgs)) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// ListDialog returns dialog messages.
func (h Handler) ListDialog(w http.ResponseWriter, r *http.Request) {
	from, to, err := h.getUUIDs(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	page, err := h.getPageParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	obj, err := h.service.ListDialog(r.Context(), from, to, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	getPostResponse := model.NewDialogListResponse(obj)
	res, err := json.Marshal(getPostResponse)
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
