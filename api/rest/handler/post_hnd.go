package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"

	"github.com/vstdy/otus-highload/api/rest/hub"
	"github.com/vstdy/otus-highload/api/rest/model"
	"github.com/vstdy/otus-highload/pkg"
)

// CreatePost creates post.
func (h Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	userUUID, err := h.getUserUUID(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var body model.CreatePostBody
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	postUUID, err := h.service.CreatePost(r.Context(), userUUID, body.Text)
	if err != nil {
		if errors.Is(err, pkg.ErrWrongCredentials) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createPostResponse := model.NewCreatePostResponse(postUUID)
	res, err := json.Marshal(createPostResponse)
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

// UpdatePost updates post.
func (h Handler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	userUUID, err := h.getUserUUID(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var body model.UpdatePostBody
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = h.service.UpdatePost(r.Context(), userUUID, body.ToCanonical())
	if err != nil {
		if errors.Is(err, pkg.ErrWrongCredentials) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

// DeletePost deletes post.
func (h Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
	userUUID, postUUID, err := h.getUUIDs(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.DeletePost(r.Context(), userUUID, postUUID)
	if err != nil {
		if errors.Is(err, pkg.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetPost returns post.
func (h Handler) GetPost(w http.ResponseWriter, r *http.Request) {
	_, postUUID, err := h.getUUIDs(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	obj, err := h.service.GetPost(r.Context(), postUUID)
	if err != nil {
		if errors.Is(err, pkg.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	getPostResponse := model.NewPostResponse(obj)
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

// GetFeed returns friends most recent posts.
func (h Handler) GetFeed(w http.ResponseWriter, r *http.Request) {
	userUUID, err := h.getUserUUID(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	page, err := h.getPageParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	objs, err := h.service.PostsFeed(r.Context(), userUUID, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	feedResponse := model.NewPostListResponse(objs)
	res, err := json.Marshal(feedResponse)
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

// GetPostedFeed returns friends most recent posts.
func (h Handler) GetPostedFeed(w http.ResponseWriter, r *http.Request) {
	ctx, logger := h.Logger(r.Context())

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Warn().Err(err).Msg("upgrade")
		return
	}

	userUUID, err := h.getUserUUID(ctx)
	if err != nil {
		logger.Warn().Err(err).Msg("get user uuid")
		return
	}

	send := make(chan []byte)
	cnt := &hub.Connect{
		Conn: conn,
		Send: send,
		User: userUUID.String(),
	}
	h.hub.Register <- cnt

	go h.readWS(cnt, logger)
	go h.writeWS(cnt, logger)
}

func (h Handler) readWS(cnt *hub.Connect, logger zerolog.Logger) {
	defer cnt.Conn.Close()

	for {
		_, _, err := cnt.Conn.NextReader()
		if err == nil {
			continue
		}
		if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
			logger.Warn().Err(err).Msg("ws closing connection error")
		}
		h.hub.Unregister <- cnt
		return
	}
}

func (h Handler) writeWS(cnt *hub.Connect, logger zerolog.Logger) {
	defer cnt.Conn.Close()

	for {
		select {
		case message, ok := <-cnt.Send:
			if !ok {
				return
			}

			err := cnt.Conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				logger.Warn().Err(err).Msg("write message")
			}
		}
	}
}
