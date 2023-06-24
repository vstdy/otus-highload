-- +goose Up
-- Chat table
CREATE TABLE "chat"
(
    "uuid"          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "participant_1" BIGINT NOT NULL,
    "participant_2" BIGINT NOT NULL
);

CREATE INDEX "chat_participants_index" ON "chat" USING btree (participant_1, participant_2);
SELECT create_distributed_table('chat', 'uuid', shard_count := 15);

-- Dialog table
CREATE TABLE "dialog"
(
    "uuid"       UUID                                     DEFAULT uuid_generate_v4(),
    "chat_id"    UUID                            NOT NULL,
    "from"       UUID                            NOT NULL,
    "to"         UUID                            NOT NULL,
    "text"       TEXT CHECK (LENGTH("text") > 0) NOT NULL,
    "created_at" TIMESTAMPTZ                     NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMPTZ                     NOT NULL DEFAULT now(),
    PRIMARY KEY ("chat_id", "uuid")
);

SELECT create_distributed_table('dialog', 'chat_id', colocate_with => 'chat');

-- +goose Down
DROP TABLE "dialog";
DROP TABLE "chat";