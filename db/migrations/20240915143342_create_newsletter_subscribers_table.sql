-- +goose Up
-- +goose StatementBegin
CREATE TABLE newsletter_subscribers (
    email TEXT NOT NULL,
    newsletter_id BIGINT NOT NULL REFERENCES newsletter(id),
    subscribed_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (email, newsletter_id)
);

CREATE INDEX idx_newsletter_id ON newsletter_subscribers (newsletter_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_newsletter_id;

DROP TABLE IF EXISTS newsletter_subscribers;
-- +goose StatementEnd
