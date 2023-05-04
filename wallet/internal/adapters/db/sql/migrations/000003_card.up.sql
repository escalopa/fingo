CREATE TABLE cards
(
  number     VARCHAR(64) PRIMARY KEY NOT NULL,
  account_id BIGINT                  NOT NULL
);

ALTER TABLE cards
  ADD FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE CASCADE;
