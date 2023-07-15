package tarantool

import (
	"context"

	"github.com/google/uuid"
	"github.com/tarantool/go-tarantool/v2"

	"github.com/vstdy/otus-highload/pkg"
)

const chatTableName = "chat"

// AddChat adds new chat.
func (st *Storage) AddChat(ctx context.Context, user1, user2 int64) (uuid.UUID, error) {
	args := []interface{}{user1, user2}

	var res uuid.UUID
	req := tarantool.NewCallRequest("add_chat").Args(args)
	resp, err := st.conn.Do(req).Get()
	if err != nil {
		return uuid.Nil, err
	}
	_ = resp

	return res, nil
}

// GetChat returns chat.
func (st *Storage) GetChat(ctx context.Context, user1, user2 int64) (uuid.UUID, error) {
	//args := []interface{}{[]interface{}{user1, user2}}
	//args := []interface{}{user1, user2}

	var res uuid.UUID
	//req := tarantool.NewCallRequest("box.space._user:select")
	//req := tarantool.NewCallRequest("box.session.user")
	//req := tarantool.NewCallRequest("box.func.get_chat:call").Args(args)
	//req := tarantool.NewEvalRequest("box.func.get_chat:call({%d,%d})").Args(args)
	req := tarantool.NewEvalRequest("get_chat(1,2)")
	//req := tarantool.NewCallRequest("get_chat").Args(args)
	resp, err := st.conn.Do(req).Get()
	if err != nil {
		return uuid.Nil, err
	}
	if len(resp.Data) == 0 {
		return uuid.Nil, pkg.ErrNotFound
	}

	return res, nil
}
