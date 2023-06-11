package psql

import (
	"context"
	"errors"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/vstdy/otus-highload/model"
	"github.com/vstdy/otus-highload/pkg"
)

const dialogTableName = "dialog"

// SendDialog sends message to dialog.
func (st *Storage) SendDialog(ctx context.Context, chatUUID, fromUUID, toUUID uuid.UUID, text string) error {
	query := `
		INSERT INTO "dialog" ("chat_id", "from", "to", "text") VALUES ($1, $2, $3, $4);
	`
	args := []interface{}{chatUUID, fromUUID, toUUID, text}

	_, err := st.masterConn.Exec(ctx, query, args...)
	if err != nil {
		pgErr := new(pgconn.PgError)
		if !errors.As(err, &pgErr) {
			return err
		}
		if pgErr.Code == pkg.NotNullViolation {
			return pkg.ErrInvalidUserArgs{Column: pgErr.ColumnName}
		}

		return err
	}

	return nil
}

// ListDialog returns dialog messages.
func (st *Storage) ListDialog(ctx context.Context, chatID uuid.UUID, page model.Page) ([]model.Dialog, error) {
	query := `
		SELECT * FROM dialog WHERE chat_id = $1 ORDER BY created_at DESC OFFSET $2 LIMIT $3
	`
	args := []interface{}{chatID, page.Offset, page.Limit}

	var objs []model.Dialog
	err := pgxscan.Select(ctx, st.masterConn, &objs, query, args...)
	if err != nil {
		return nil, err
	}

	return objs, nil
}
