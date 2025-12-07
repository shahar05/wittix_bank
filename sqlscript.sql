CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    currency TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE journal_entries (
    entry_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    external_id TEXT NOT NULL UNIQUE,
    description TEXT,
    posted_at TIMESTAMP NOT NULL DEFAULT NOW(),
    reversal_of UUID REFERENCES journal_entries(entry_id),
    idempotency_key TEXT NOT NULL UNIQUE
);

CREATE TYPE side_enum AS ENUM ('debit', 'credit');
CREATE TABLE journal_lines (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    entry_id UUID NOT NULL REFERENCES journal_entries(entry_id) ON DELETE CASCADE,
    account_id UUID NOT NULL REFERENCES accounts(id), -- FIXED
    side side_enum NOT NULL,
    amount NUMERIC(20,8) NOT NULL CHECK (amount > 0)
);
