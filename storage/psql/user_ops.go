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

// CreateUser adds given user object to storage
func (st *Storage) CreateUser(ctx context.Context, rawObj model.User) (model.User, error) {
	query := `
		INSERT INTO "user" ("first_name","second_name","age","biography","city","password")
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING *;
	`
	args := []interface{}{rawObj.FirstName, rawObj.SecondName, rawObj.Age, rawObj.Biography, rawObj.City, rawObj.Password}

	var dbObj schema.User
	err := pgxscan.Get(ctx, st.masterConn, &dbObj, query, args...)
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
		WHERE id = $1
			AND deleted_at IS NULL;
	`
	args := []interface{}{rawObj.ID}

	var dbObj schema.User
	err := pgxscan.Get(ctx, st.masterConn, &dbObj, query, args...)
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
	conn := st.masterConn
	if st.asyncReplicaConn != nil {
		conn = st.asyncReplicaConn
	}

	query := `
		SELECT *
		FROM "user"
		WHERE uuid = $1
			AND deleted_at IS NULL;
	`
	args := []interface{}{userUUID}

	var dbObj schema.User
	err := pgxscan.Get(ctx, conn, &dbObj, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, pkg.ErrNotFound
		}

		return model.User{}, err
	}

	return dbObj.ToCanonical(), nil
}

// SearchUsers searches users.
func (st *Storage) SearchUsers(ctx context.Context, searchParams model.SearchUser) ([]model.User, error) {
	conn := st.masterConn
	if st.asyncReplicaConn != nil {
		conn = st.asyncReplicaConn
	}

	query := `
		SELECT *
		FROM "user"
		WHERE lower(first_name) ^@ lower($1)
			AND lower(second_name) ^@ lower($2)
			AND deleted_at IS NULL
		ORDER BY id;
	`
	args := []interface{}{searchParams.FirstName, searchParams.LastName}

	var dbObjs schema.Users
	err := pgxscan.Select(ctx, conn, &dbObjs, query, args...)
	if err != nil {
		return nil, err
	}

	return dbObjs.ToCanonical(), nil
}
