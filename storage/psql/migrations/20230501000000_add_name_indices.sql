-- +goose Up
CREATE INDEX IF NOT EXISTS "first_name_index"
    ON "user" USING btree (lower(first_name) COLLATE "C");

CREATE INDEX IF NOT EXISTS "second_name_index"
    ON "user" USING btree (lower(second_name) COLLATE "C");

-- +goose Down
DROP INDEX IF EXISTS "first_name_index";
DROP INDEX IF EXISTS "second_name_index";