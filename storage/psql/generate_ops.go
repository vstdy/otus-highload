package psql

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/vstdy/otus-highload/model"
)

// CopyUsers copies users to storage
func (st *Storage) CopyUsers(ctx context.Context, objs []model.User) (int64, error) {
	tableName := pgx.Identifier{userTableName}
	columnNames := []string{"first_name", "second_name", "age", "city", "password"}
	rowSrc := pgx.CopyFromSlice(len(objs), func(i int) ([]interface{}, error) {
		return []interface{}{objs[i].FirstName, objs[i].SecondName, objs[i].Age, objs[i].City, objs[i].Password}, nil
	})

	count, err := st.masterConn.CopyFrom(ctx, tableName, columnNames, rowSrc)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// CopyFriends copies friends to storage
func (st *Storage) CopyFriends(ctx context.Context, objs []model.Friend) (int64, error) {
	tableName := pgx.Identifier{friendTableName}
	columnNames := []string{"user_id", "friend_id"}
	rowSrc := pgx.CopyFromSlice(len(objs), func(i int) ([]interface{}, error) {
		return []interface{}{objs[i].UserID, objs[i].FriendID}, nil
	})

	count, err := st.masterConn.CopyFrom(ctx, tableName, columnNames, rowSrc)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// CopyPosts copies posts to storage
func (st *Storage) CopyPosts(ctx context.Context, objs []model.Post) (int64, error) {
	tableName := pgx.Identifier{postTableName}
	columnNames := []string{"text", "author_id"}
	rowSrc := pgx.CopyFromSlice(len(objs), func(i int) ([]interface{}, error) {
		return []interface{}{objs[i].Text, objs[i].AuthorID}, nil
	})

	count, err := st.masterConn.CopyFrom(ctx, tableName, columnNames, rowSrc)
	if err != nil {
		return 0, err
	}

	return count, nil
}
