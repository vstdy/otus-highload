package psql

import (
	"context"
	"errors"

	"github.com/georgysavva/scany/v2/pgxscan"
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

// GetFriendSetters returns users who added given user to friends.
func (st *Storage) GetFriendSetters(ctx context.Context, friendID int64) ([]uuid.UUID, error) {
	query := `
		SELECT u.uuid
		FROM friend f
		INNER JOIN "user" u on u.id = f.user_id
		WHERE friend_id = $1;
	`
	args := []interface{}{friendID}

	var users []uuid.UUID
	err := pgxscan.Select(ctx, st.masterConn, &users, query, args...)
	if err != nil {
		return nil, err
	}

	return users, nil
}
