CREATE TABLE IF NOT EXISTS users
(
    id         BIGSERIAL PRIMARY KEY,
    first_name TEXT      NOT NULL,
    username   TEXT,
    balance    INTEGER   NOT NULL DEFAULT 0,
    phone      TEXT,
    partner_id BIGINT,
    bonus_used BOOLEAN   NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT valid_balance CHECK (balance >= 0)
);

CREATE TRIGGER update_users_updated_at
    BEFORE
        UPDATE
    ON public.users
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();