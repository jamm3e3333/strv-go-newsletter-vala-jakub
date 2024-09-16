-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE IF NOT EXISTS newsletter_public_id_seq
    START WITH 1100
    INCREMENT BY 1;

CREATE TABLE newsletter (
    id BIGSERIAL PRIMARY KEY,
    public_id BIGINT UNIQUE NOT NULL DEFAULT nextval('newsletter_public_id_seq'::regclass),
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
