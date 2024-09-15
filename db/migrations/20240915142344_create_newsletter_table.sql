-- +goose Up
-- +goose StatementBegin
CREATE TABLE newsletter (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT REFERENCES client(id) NOT NULL,
    name TEXT NOT NULL,
    description TEXT DEFAULT '',
    create_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE OR REPLACE FUNCTION on_update_timestamp ()
    RETURNS TRIGGER
AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$
    LANGUAGE plpgsql;

CREATE TRIGGER newsletter_updated_at
    BEFORE UPDATE ON newsletter
    FOR EACH ROW
EXECUTE PROCEDURE on_update_timestamp ();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS newsletter;
-- +goose StatementEnd
