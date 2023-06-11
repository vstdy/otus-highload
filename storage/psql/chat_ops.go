package psql

import (
	"context"
	"errors"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/vstdy/otus-highload/pkg"
)

const chatTableName = "chat"

// AddChat adds new chat.
func (st *Storage) AddChat(ctx context.Context, userUUID1, userUUID2 uuid.UUID) (uuid.UUID, error) {
	query := `
		WITH
			user_1 AS (SELECT (SELECT id FROM "user" WHERE uuid = $1 AND deleted_at IS NULL) id_1),
			user_2 AS (SELECT (SELECT id FROM "user" WHERE uuid = $2 AND deleted_at IS NULL) id_2)
		INSERT INTO "chat" (participant_1, participant_2)
		SELECT 
			CASE WHEN id_1 > id_2 THEN id_2 WHEN id_1 < id_2 THEN id_1 END,
			CASE WHEN id_1 > id_2 THEN id_1 WHEN id_1 < id_2 THEN id_2 END
		FROM user_1, user_2
		RETURNING uuid;
	`
	args := []interface{}{userUUID1, userUUID2}

	var res uuid.UUID
	err := pgxscan.Get(ctx, st.masterConn, &res, query, args...)
	if err != nil {
		pgErr := new(pgconn.PgError)
		if !errors.As(err, &pgErr) {
			return uuid.Nil, err
		}
		if pgErr.Code == pkg.NotNullViolation {
			return uuid.Nil, pkg.ErrUserNotFound
		}
	}

	return res, nil
}

// GetChat returns chat.
func (st *Storage) GetChat(ctx context.Context, userUUID1, userUUID2 uuid.UUID) (uuid.UUID, error) {
	query := `
		WITH
			user_1 AS (SELECT (SELECT id FROM "user" WHERE uuid = $1 AND deleted_at IS NULL) id_1),
			user_2 AS (SELECT (SELECT id FROM "user" WHERE uuid = $2 AND deleted_at IS NULL) id_2)
		SELECT chat.uuid
		FROM "chat", user_1, user_2
		WHERE participant_1 = CASE WHEN id_1 > id_2 THEN id_2 ELSE id_1 END
			AND participant_2 = CASE WHEN id_1 > id_2 THEN id_1 ELSE id_2 END
	`
	args := []interface{}{userUUID1, userUUID2}

	var res uuid.UUID
	err := pgxscan.Get(ctx, st.masterConn, &res, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, pkg.ErrNotFound
		}

		return uuid.Nil, err
	}

	return res, nil
}
