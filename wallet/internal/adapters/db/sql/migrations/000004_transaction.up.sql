CREATE TYPE transaction_type AS ENUM ('deposit', 'withdrawal', 'transfer');

CREATE TABLE transactions
(
  id                     SERIAL PRIMARY KEY,
  type                   transaction_type NOT NULL,
  amount                 NUMERIC(10, 2)   NOT NULL,
  source_account_id      INTEGER REFERENCES accounts (id),
  destination_account_id INTEGER REFERENCES accounts (id),
  created_at             TIMESTAMP        NOT NULL DEFAULT NOW()
);
