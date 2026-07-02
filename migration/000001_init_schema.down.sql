DROP INDEX IF EXISTS idx_transactions_account_id;
DROP INDEX IF EXISTS idx_accounts_deleted_at;

DROP TABLE IF EXISTS transactions;

DROP TABLE IF EXISTS accounts;

DROP EXTENSION IF EXISTS "uuid-ossp";