package psql

import (
	"context"
	"errors"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/vstdy/otus-highload/pkg"
)

const chatTableName = "chat"

// AddChat adds new chat.
func (st *Storage) AddChat(ctx context.Context, user1, user2 int64) (res uuid.UUID, err error) {
	query := `
		INSERT INTO "chat" (participant_1, participant_2)
		VALUES ($1, $2)
		RETURNING uuid;
	`
	args := []interface{}{user1, user2}

	err = pgxscan.Get(ctx, st.masterConn, &res, query, args...)

	return res, err
}

// GetChat returns chat.
func (st *Storage) GetChat(ctx context.Context, user1, user2 int64) (res uuid.UUID, err error) {
	query := `
		SELECT uuid
		FROM "chat"
		WHERE participant_1 = $1
			AND participant_2 = $2
	`
	args := []interface{}{user1, user2}

	err = pgxscan.Get(ctx, st.masterConn, &res, query, args...)
	if errors.Is(err, pgx.ErrNoRows) {
		return uuid.Nil, pkg.ErrNotFound
	}

	return res, err
}
