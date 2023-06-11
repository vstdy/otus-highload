package psql

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/vstdy/otus-highload/pkg"
)

const friendTableName = "friend"

// SetFriend adds friend to user.
func (st *Storage) SetFriend(ctx context.Context, userUUID, friendUUID uuid.UUID) error {
	query := `
		INSERT INTO "friend"
		VALUES ((SELECT id FROM "user" WHERE uuid = $1 AND deleted_at IS NULL),
				(SELECT id FROM "user" WHERE uuid = $2 AND deleted_at IS NULL))
		ON CONFLICT DO NOTHING;
	`
	args := []interface{}{userUUID, friendUUID}

	_, err := st.masterConn.Exec(ctx, query, args...)
	if err != nil {
		pgErr := new(pgconn.PgError)
		if !errors.As(err, &pgErr) {
			return err
		}
		if pgErr.Code == pkg.NotNullViolation {
			return pkg.ErrInvalidUserArgs{Column: pgErr.ColumnName}
		}
	}

	return err
}

// DeleteFriend deletes user's friend.
func (st *Storage) DeleteFriend(ctx context.Context, userUUID, friendUUID uuid.UUID) error {
	query := `
		DELETE
		FROM "friend"
		WHERE user_id = (SELECT id FROM "user" WHERE uuid = $1 AND deleted_at IS NULL)
		  AND friend_id = (SELECT id FROM "user" WHERE uuid = $2 AND deleted_at IS NULL);
	`
	args := []interface{}{userUUID, friendUUID}

	res, err := st.masterConn.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return pkg.ErrNotFound
	}

	return err
}
