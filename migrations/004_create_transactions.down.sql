DROP TRIGGER IF EXISTS update_transactions_updated_at ON transactions CASCADE;
DROP TABLE IF EXISTS transactions CASCADE;
DROP TYPE IF EXISTS transaction_type CASCADE;