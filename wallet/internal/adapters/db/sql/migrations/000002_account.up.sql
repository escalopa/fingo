CREATE TABLE accounts
(
  id      BIGSERIAL PRIMARY KEY NOT NULL,
  user_id BIGINT                NOT NULL,
  name    varchar(20)           NOT NULL,
  balance DOUBLE PRECISION      NOT NULL DEFAULT 0.0000
);

ALTER TABLE accounts
  ADD FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE;
