CREATE TABLE accounts
(
  id      SERIAL PRIMARY KEY NOT NULL,
  user_id INTEGER            NOT NULL,
  name    varchar(20)        NOT NULL,
  balance NUMERIC(12, 4)     NOT NULL DEFAULT 0.0000
);

ALTER TABLE accounts
  ADD FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE;
