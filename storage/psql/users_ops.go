package psql

import (
	"context"
	"errors"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/vstdy/otus-highload/model"
	"github.com/vstdy/otus-highload/pkg"
	"github.com/vstdy/otus-highload/storage/psql/schema"
)

const userTableName = "user"

// CreateUser adds given url objects to storage
func (st *Storage) CreateUser(ctx context.Context, rawObj model.User) (model.User, error) {
	query := `
		INSERT INTO "user" ("first_name","second_name","age","biography","city","password")
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING *;
	`
	args := []interface{}{rawObj.FirstName, rawObj.SecondName, rawObj.Age, rawObj.Biography, rawObj.City, rawObj.Password}

	var dbObj schema.User
	err := pgxscan.Get(ctx, st.db, &dbObj, query, args...)
	if err != nil {
		return model.User{}, err
	}

	return dbObj.ToCanonical(), nil
}

// AuthenticateUser verifies the identity of credentials.
func (st *Storage) AuthenticateUser(ctx context.Context, rawObj model.User) (model.User, error) {
	query := `
		SELECT *
		FROM "user"
		WHERE id = $1;
	`
	args := []interface{}{rawObj.ID}

	var dbObj schema.User
	err := pgxscan.Get(ctx, st.db, &dbObj, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, pkg.ErrWrongCredentials
		}

		return model.User{}, err
	}

	return dbObj.ToCanonical(), nil
}

// GetUser returns user data.
func (st *Storage) GetUser(ctx context.Context, userUUID uuid.UUID) (model.User, error) {
	query := `
		SELECT *
		FROM "user"
		WHERE uuid = $1;
	`
	args := []interface{}{userUUID}

	var dbObj schema.User
	err := pgxscan.Get(ctx, st.db, &dbObj, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, pkg.ErrNotFound
		}

		return model.User{}, err
	}

	return dbObj.ToCanonical(), nil
}
