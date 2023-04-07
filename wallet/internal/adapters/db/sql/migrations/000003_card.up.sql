CREATE TABLE cards
(
  number     VARCHAR(16) PRIMARY KEY NOT NULL,
  account_id BIGINT                  NOT NULL
);

ALTER TABLE cards
  ADD FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE CASCADE;

-- This is a check constraint that ensures that the card number is a 16-digit
ALTER TABLE cards
  ADD CONSTRAINT cards_number_check CHECK (number ~ '^[0-9]{16}$');
