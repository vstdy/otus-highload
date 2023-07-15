package tarantool

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/tarantool/go-tarantool/v2"

	"github.com/vstdy/otus-highload/model"
)

const dialogTableName = "dialog"

// SendDialog sends message to dialog.
func (st *Storage) SendDialog(ctx context.Context, chatUUID, fromUUID, toUUID uuid.UUID, text string) error {
	args := []interface{}{chatUUID, fromUUID, toUUID, text, time.Now(), time.Now()}

	req := tarantool.NewCallRequest("send_dialog").Args(args)
	resp, err := st.conn.Do(req).Get()
	if err != nil {
		return err
	}
	_ = resp

	return nil
}

// ListDialog returns dialog messages.
func (st *Storage) ListDialog(ctx context.Context, chatID uuid.UUID, page model.Page) ([]model.Dialog, error) {
	args := []interface{}{chatID, page.Offset, page.Limit}

	var objs []model.Dialog
	req := tarantool.NewCallRequest("list_dialog").Args(args)
	resp, err := st.conn.Do(req).Get()
	if err != nil {
		return nil, err
	}
	_ = resp

	return objs, nil
}
