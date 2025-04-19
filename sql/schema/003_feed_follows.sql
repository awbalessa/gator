-- +goose Up
CREATE TABLE feed_follows (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    feed_id UUID NOT NULL REFERENCES feeds (id) ON DELETE CASCADE,
    UNIQUE (user_id, feed_id)
);

-- +goose Down
ALTER TABLE feed_follows
DROP CONSTRAINT feed_follows_user_id_fkey;

ALTER TABLE feed_follows
DROP CONSTRAINT feed_follows_feed_id_fkey;

DROP TABLE feed_follows;
