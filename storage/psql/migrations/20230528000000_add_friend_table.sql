-- +goose Up
-- Friend table
CREATE TABLE "friend"
(
    "user_id"   BIGSERIAL,
    "friend_id" BIGSERIAL,
    PRIMARY KEY ("user_id", "friend_id"),
    CHECK ("user_id" <> "friend_id")
);

-- +goose Down
DROP TABLE "friend";