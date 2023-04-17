-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- User table
CREATE TABLE "user"
(
    "id"          BIGSERIAL PRIMARY KEY,
    "uuid"        UUID                 DEFAULT uuid_generate_v4(),
    "first_name"  VARCHAR(50) NOT NULL,
    "second_name" VARCHAR(50) NOT NULL,
    "age"         SMALLINT    NOT NULL,
    "biography"   TEXT        NOT NULL,
    "city"        VARCHAR     NOT NULL,
    "password"    VARCHAR     NOT NULL,
    "created_at"  TIMESTAMPTZ NOT NULL DEFAULT now(),
    "updated_at"  TIMESTAMPTZ NOT NULL DEFAULT now(),
    "deleted_at"  TIMESTAMPTZ
);

-- +goose Down
DROP TABLE "user";