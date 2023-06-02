package handler

import (
	"errors"
	"net/http"

	"github.com/vstdy/otus-highload/pkg"
)

// SetFriend adds friend to user.
func (h Handler) SetFriend(w http.ResponseWriter, r *http.Request) {
	userUUID, friendUUID, err := h.getUUIDs(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.SetFriend(r.Context(), userUUID, friendUUID)
	if err != nil {
		if errors.As(err, new(pkg.ErrSetFriendInvalidArgs)) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// DeleteFriend deletes user's friend.
func (h Handler) DeleteFriend(w http.ResponseWriter, r *http.Request) {
	userUUID, friendUUID, err := h.getUUIDs(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.DeleteFriend(r.Context(), userUUID, friendUUID)
	if err != nil {
		if errors.Is(err, pkg.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
