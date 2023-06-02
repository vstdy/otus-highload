package psql

import (
	"context"
	"errors"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/vstdy/otus-highload/model"
	"github.com/vstdy/otus-highload/pkg"
)

const postTableName = "post"

// CreatePost creates post.
func (st *Storage) CreatePost(ctx context.Context, userUUID uuid.UUID, text string) (uuid.UUID, error) {
	query := `
		INSERT INTO "post" (author_id, text)
		VALUES ((SELECT id FROM "user" WHERE uuid = $1 AND deleted_at IS NULL), $2)
		RETURNING uuid;
	`
	args := []interface{}{userUUID, text}

	var postUUID uuid.UUID
	err := pgxscan.Get(ctx, st.masterConn, &postUUID, query, args...)
	if err != nil {
		pgErr := new(pgconn.PgError)
		if !errors.As(err, &pgErr) {
			return uuid.Nil, err
		}
		if pgErr.Code == pkg.NotNullViolation {
			return uuid.Nil, pkg.ErrWrongCredentials
		}

		return uuid.Nil, err
	}

	return postUUID, nil
}

// UpdatePost updates post.
func (st *Storage) UpdatePost(ctx context.Context, userUUID uuid.UUID, post model.Post) error {
	query := `
		UPDATE "post"
		SET "text" = $3
		WHERE uuid = $2
		  AND author_id = (SELECT id FROM "user" WHERE uuid = $1 AND deleted_at IS NULL);
	`
	args := []interface{}{userUUID, post.UUID, post.Text}

	res, err := st.masterConn.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return pkg.ErrNotFound
	}

	return nil
}

// DeletePost deletes post.
func (st *Storage) DeletePost(ctx context.Context, userUUID, postUUID uuid.UUID) error {
	query := `
		DELETE FROM "post"
		WHERE uuid = $2
		  AND author_id = (SELECT id FROM "user" WHERE uuid = $1 AND deleted_at IS NULL);
	`
	args := []interface{}{userUUID, postUUID}

	res, err := st.masterConn.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return pkg.ErrNotFound
	}

	return nil
}

// GetPost returns post.
func (st *Storage) GetPost(ctx context.Context, postUUID uuid.UUID) (model.PostExt, error) {
	query := `
		SELECT
			p.*,
			u.uuid author_uuid
		FROM "post" p
		LEFT JOIN "user" u on p.author_id = u.id
		WHERE p.uuid = $1;
	`
	args := []interface{}{postUUID}

	var obj model.PostExt
	err := pgxscan.Get(ctx, st.masterConn, &obj, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.PostExt{}, pkg.ErrNotFound
		}

		return model.PostExt{}, err
	}

	return obj, nil
}

// PostsFeed returns friends most recent posts.
func (st *Storage) PostsFeed(ctx context.Context, userUUID uuid.UUID, page model.Page) ([]model.PostExt, error) {
	query := `
		SELECT p.*, pu.uuid author_uuid
		FROM "post" p
			LEFT JOIN "user" pu ON p.author_id = pu.id
		WHERE p.author_id IN (
			SELECT friend_id FROM "friend"
				LEFT JOIN "user" fu ON fu.id = friend.user_id
			WHERE fu.uuid = $1)
		ORDER BY p.created_at DESC
		OFFSET $2
		LIMIT $3;
	`
	args := []interface{}{userUUID, page.Offset, page.Limit}

	var objs []model.PostExt
	err := pgxscan.Select(ctx, st.masterConn, &objs, query, args...)
	if err != nil {
		return nil, err
	}

	return objs, nil
}
