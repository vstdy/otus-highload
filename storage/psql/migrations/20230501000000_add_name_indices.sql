-- +goose Up
CREATE INDEX IF NOT EXISTS "user_first_name_index"
    ON "user" USING btree (lower("first_name") COLLATE "C");

CREATE INDEX IF NOT EXISTS "user_second_name_index"
    ON "user" USING btree (lower("second_name") COLLATE "C");

-- +goose Down
DROP INDEX IF EXISTS "user_first_name_index";
DROP INDEX IF EXISTS "user_second_name_index";