CREATE TABLE IF NOT EXISTS coupons
(
    id             BIGSERIAL PRIMARY KEY,
    owner_id       BIGINT    NOT NULL,
    receiver_id    BIGINT    NOT NULL,
    transaction_id UUID      NOT NULL,
    amount         INTEGER   NOT NULL,
    code           TEXT      NOT NULL,
    status         TEXT      NOT NULL,
    expired_at     TIMESTAMP NOT NULL,
    created_at     TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMP NOT NULL DEFAULT NOW(),


    CONSTRAINT fk_owner_id FOREIGN KEY (owner_id) REFERENCES users (id),
    CONSTRAINT fk_receiver_id FOREIGN KEY (receiver_id) REFERENCES users (id),
    CONSTRAINT fk_transaction_id FOREIGN KEY (transaction_id) REFERENCES transactions (id)
);

CREATE INDEX IF NOT EXISTS owner_id_idx ON coupons (owner_id);

CREATE INDEX IF NOT EXISTS receiver_id_idx ON coupons (receiver_id);

CREATE INDEX IF NOT EXISTS code_idx ON coupons (code);

CREATE TRIGGER update_coupons_updated_at BEFORE
UPDATE ON public.coupons FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column ();