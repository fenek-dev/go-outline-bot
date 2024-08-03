CREATE TABLE IF NOT EXISTS subscriptions
(
    id            BIGSERIAL PRIMARY KEY,
    user_id       BIGINT    NOT NULL,
    server_id     BIGINT    NOT NULL,
    tariff_id     BIGINT    NOT NULL,
    initial_price INTEGER   NOT NULL,
    key           TEXT      NOT NULL,
    status        TEXT      NOT NULL,
    expired_at    TIMESTAMP NOT NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT fk_server_id FOREIGN KEY (server_id) REFERENCES servers (id),
    CONSTRAINT fk_tariff_id FOREIGN KEY (tariff_id) REFERENCES tariffs (id)

);

CREATE INDEX IF NOT EXISTS user_id_idx ON subscriptions (user_id);


CREATE TRIGGER update_subscriptions_updated_at
    BEFORE UPDATE
    ON public.subscriptions
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();