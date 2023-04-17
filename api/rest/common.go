package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"

	canonical "github.com/vstdy/otus-highload/model"
)

const claimsField = "uuid"

// addJWTCookie adds a jwt cookie to the response.
func (h Handler) addJWTCookie(w http.ResponseWriter, obj canonical.User) error {
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

// getUserID retrieves the user ID from the context.
func (h Handler) getUserID(ctx context.Context) (uuid.UUID, error) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	userID, err := uuid.Parse(claims[claimsField].(string))
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}
