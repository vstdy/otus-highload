package handler

import (
	"context"
	"encoding/json"

	"github.com/go-chi/jwtauth/v5"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"

	"github.com/vstdy/otus-highload/api/rest/hub"
	"github.com/vstdy/otus-highload/api/rest/model"
	"github.com/vstdy/otus-highload/pkg/logging"
	"github.com/vstdy/otus-highload/service/project"
)

const (
	serviceName = "otus-project server"
)

// Handler keeps handler dependencies.
type Handler struct {
	service  project.IService
	upgrader websocket.Upgrader
	jwtAuth  *jwtauth.JWTAuth
	logLevel zerolog.Level
	hub      *hub.Hub
}

// NewHandler returns a new Handler instance.
func NewHandler(
	service project.IService,
	upgrader websocket.Upgrader,
	jwtAuth *jwtauth.JWTAuth,
	logLevel zerolog.Level,
	hub *hub.Hub,
) Handler {
	handler := Handler{service: service, upgrader: upgrader, jwtAuth: jwtAuth, logLevel: logLevel, hub: hub}
	go handler.broadcast()

	return handler
}

// Logger returns logger with service field set.
func (h Handler) Logger(ctx context.Context) (context.Context, zerolog.Logger) {
	ctx, logger := logging.GetCtxLogger(ctx, logging.WithLogLevel(h.logLevel))
	logger = logger.With().Str(logging.ServiceKey, serviceName).Logger()

	return ctx, logger
}

func (h *Handler) broadcast() {
	for {
		select {
		case conn := <-h.hub.Register:
			if userClients, ok := h.hub.Clients[conn.User]; ok {
				userClients[conn.Conn] = conn.Send
				continue
			}
			h.hub.Clients[conn.User] = map[*websocket.Conn]chan []byte{conn.Conn: conn.Send}
		case conn := <-h.hub.Unregister:
			if userClients, ok := h.hub.Clients[conn.User]; ok {
				if _, ok := userClients[conn.Conn]; ok {
					delete(userClients, conn.Conn)
					close(conn.Send)
					if len(userClients) == 0 {
						delete(h.hub.Clients, conn.User)
					}
				}
			}
		case postNtf := <-h.service.GetHub():
			postResponse := model.NewPostResponse(postNtf.PostExt)
			message, err := json.Marshal(postResponse)
			if err != nil {
				continue
			}
			for _, user := range postNtf.Users {
				clients, ok := h.hub.Clients[user.String()]
				if !ok {
					continue
				}
				for client, send := range clients {
					select {
					case send <- message:
					default:
						close(send)
						delete(clients, client)
						if len(clients) == 0 {
							delete(h.hub.Clients, user.String())
						}
					}
				}
			}
		}
	}
}
