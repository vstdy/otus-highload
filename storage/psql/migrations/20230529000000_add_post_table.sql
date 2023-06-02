-- +goose Up
-- Post table
CREATE TABLE "post"
(
    "id"         BIGSERIAL PRIMARY KEY,
    "uuid"       UUID                                     DEFAULT uuid_generate_v4(),
    "text"       TEXT CHECK (LENGTH("text") > 0) NOT NULL,
    "author_id"  BIGSERIAL                       NOT NULL,
    "created_at" TIMESTAMPTZ                     NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMPTZ                     NOT NULL DEFAULT now()
);

CREATE INDEX "post_uuid_index" ON "post" USING btree ("uuid");

-- +goose Down
DROP TABLE "post";