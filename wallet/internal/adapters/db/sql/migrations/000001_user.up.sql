CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users
(
  id          BIGSERIAL PRIMARY KEY NOT NULL,
  external_id uuid UNIQUE        NOT NULL
);
