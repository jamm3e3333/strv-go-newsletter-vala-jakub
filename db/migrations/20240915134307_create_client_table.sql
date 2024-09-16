-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE IF NOT EXISTS client_public_id_seq
    START WITH 1100
    INCREMENT BY 1;

CREATE TABLE client (
    id BIGSERIAL PRIMARY KEY,
    public_id BIGINT UNIQUE NOT NULL DEFAULT nextval('client_public_id_seq'::regclass),
    email TEXT UNIQUE NOT NULL,
    hashed_password TEXT NOT NULL,
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

CREATE TRIGGER client_updated_at
    BEFORE UPDATE ON client
    FOR EACH ROW
EXECUTE PROCEDURE on_update_timestamp ();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS client;
-- +goose StatementEnd
