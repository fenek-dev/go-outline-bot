CREATE TABLE IF NOT EXISTS transactions
(
    id          uuid PRIMARY KEY   DEFAULT gen_random_uuid(),
    user_id     BIGINT    NOT NULL,
    amount      INTEGER   NOT NULL,
    type        TEXT      NOT NULL,
    status      TEXT      NOT NULL,
    external_id TEXT      NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE INDEX IF NOT EXISTS user_id_idx ON transactions (user_id);
CREATE INDEX IF NOT EXISTS external_id_idx ON transactions (external_id);
CREATE INDEX IF NOT EXISTS status_idx ON transactions (status);

CREATE TRIGGER update_transactions_updated_at
    BEFORE UPDATE
    ON public.transactions
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();