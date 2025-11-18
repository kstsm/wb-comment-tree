-- +goose Up
CREATE TABLE IF NOT EXISTS comments
(
    id         BIGSERIAL PRIMARY KEY,
    parent_id  BIGINT REFERENCES comments (id) ON DELETE CASCADE,
    comment    TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_comments_parent_id ON comments (parent_id);
CREATE INDEX IF NOT EXISTS idx_comments_created_at ON comments (created_at DESC);

CREATE INDEX IF NOT EXISTS idx_comments_comment_gin ON comments USING gin (to_tsvector('simple', comment));


-- +goose Down
DROP TABLE comments;



