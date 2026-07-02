CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS accounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    balance BIGINT NOT NULL DEFAULT 0,              
    currency VARCHAR(3) NOT NULL,                   
    is_locked BOOLEAN NOT NULL DEFAULT FALSE,       
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE             
);


CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id UUID NOT NULL,
    to_account_id UUID,                             
    amount BIGINT NOT NULL,                         
    transaction_type VARCHAR(20) NOT NULL,          
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,

    CONSTRAINT fk_account FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE,
    CONSTRAINT fk_to_account FOREIGN KEY (to_account_id) REFERENCES accounts(id) ON DELETE SET NULL
);


CREATE INDEX IF NOT EXISTS idx_accounts_deleted_at ON accounts(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_transactions_account_id ON transactions(account_id);