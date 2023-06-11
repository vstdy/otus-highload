-- +goose Up
-- Post table
CREATE TABLE "post"
(
    "uuid"       UUID PRIMARY KEY                         DEFAULT uuid_generate_v4(),
    "text"       TEXT CHECK (LENGTH("text") > 0) NOT NULL,
    "author_id"  BIGINT                          NOT NULL,
    "created_at" TIMESTAMPTZ                     NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMPTZ                     NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE "post";