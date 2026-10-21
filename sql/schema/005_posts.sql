-- +goose Up
CREATE TABLE posts (
    id UUID PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    title TEXT NOT NULL,
    url TEXT UNIQUE NOT NULL,
    description TEXT,
    published_at TIMESTAMPTZ,
    feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE
);

CREATE INDEX posts_feed_id_idx ON posts (feed_id);

-- +goose Down
DROP TABLE posts;
