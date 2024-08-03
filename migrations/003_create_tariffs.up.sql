CREATE TABLE IF NOT EXISTS tariffs
(
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT      NOT NULL,
    price      INTEGER   NOT NULL,
    bandwidth  INTEGER   NOT NULL DEFAULT 0,
    server_id  BIGINT    NOT NULL,
    active     BOOLEAN   NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_server_id FOREIGN KEY (server_id) REFERENCES servers (id)
);

CREATE TRIGGER update_tariffs_updated_at
    BEFORE UPDATE
    ON public.tariffs
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();