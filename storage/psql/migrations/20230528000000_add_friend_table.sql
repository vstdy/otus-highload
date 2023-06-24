-- +goose Up
-- Friend table
CREATE TABLE "friend"
(
    "user_id"   BIGINT,
    "friend_id" BIGINT,
    PRIMARY KEY ("user_id", "friend_id"),
    CHECK ("user_id" <> "friend_id")
);

CREATE INDEX "friend_friend_id_index" ON "friend" USING btree ("friend_id");

-- +goose Down
DROP TABLE "friend";