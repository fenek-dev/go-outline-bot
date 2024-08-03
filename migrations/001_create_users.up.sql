CREATE TABLE IF NOT EXISTS users
(
    id         BIGSERIAL PRIMARY KEY,
    first_name TEXT      NOT NULL,
    username   TEXT,
    balance    INTEGER   NOT NULL DEFAULT 0,
    phone      TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TRIGGER update_users_updated_at BEFORE
UPDATE ON public.users FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column ();