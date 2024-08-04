CREATE TABLE IF NOT EXISTS servers
(
    id              BIGSERIAL PRIMARY KEY,
    name            TEXT       NOT NULL,
    country_code    VARCHAR(3) NOT NULL,
    ping            INTEGER    NOT NULL DEFAULT 0,
    ip              TEXT       NOT NULL,
    port            TEXT       NOT NULL,
    api_key         TEXT       NOT NULL,
    certificate     TEXT       NOT NULL,
    max_connections INTEGER    NOT NULL DEFAULT 0,
    total_bandwidth INTEGER    NOT NULL DEFAULT 0,
    active          BOOLEAN    NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMP  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP  NOT NULL DEFAULT NOW()
);

CREATE TRIGGER update_servers_updated_at
    BEFORE UPDATE
    ON public.servers
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();