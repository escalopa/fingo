CREATE TABLE currency
(
  id   BIGSERIAL PRIMARY KEY NOT NULL,
  name VARCHAR(3)              NOT NULL -- USD, EGP, EUR, GBP, RUB
);

INSERT INTO currency (name)
VALUES ('USD'),
       ('EGP'),
       ('EUR'),
       ('GBP'),
       ('RUB');

ALTER TABLE accounts
  ADD COLUMN currency_id BIGINT NOT NULL default 1;

ALTER TABLE accounts
  ADD FOREIGN KEY (currency_id) REFERENCES currency (id);
