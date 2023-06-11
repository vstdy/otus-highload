package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"

	"github.com/vstdy/otus-highload/model"
)

const claimsField = "uuid"

// addJWTCookie adds a jwt cookie to the response.
func (h Handler) addJWTCookie(w http.ResponseWriter, obj model.User) error {
	claims := map[string]interface{}{
		claimsField: obj.UUID,
	}
	_, token, err := h.jwtAuth.Encode(claims)
	if err != nil {
		return fmt.Errorf("auth cookie: %v", err)
	}

	cookie := http.Cookie{
		Name:  "jwt",
		Value: token,
		Path:  "/",
	}
	http.SetCookie(w, &cookie)

	return nil
}

// getUserUUID retrieves the user UUID from the context.
func (h Handler) getUserUUID(ctx context.Context) (uuid.UUID, error) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	userUUID, err := uuid.Parse(claims[claimsField].(string))
	if err != nil {
		return uuid.Nil, err
	}

	return userUUID, nil
}

func (h Handler) getUUIDs(r *http.Request) (uuid.UUID, uuid.UUID, error) {
	userUUID, err := h.getUserUUID(r.Context())
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	paramUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	return userUUID, paramUUID, err
}

const defaultPageLimit = 10

func (h Handler) getPageParams(r *http.Request) (model.Page, error) {
	var err error
	page := model.Page{Limit: defaultPageLimit}

	offsetParam := r.URL.Query().Get("offset")
	if offsetParam != "" {
		page.Offset, err = strconv.Atoi(offsetParam)
		if err != nil {
			return model.Page{}, err
		}
	}

	limitParam := r.URL.Query().Get("limit")
	if limitParam != "" {
		page.Limit, err = strconv.Atoi(limitParam)
		if err != nil {
			return model.Page{}, err
		}
	}

	return page, nil
}
