CREATE TABLE "users"
(
  "id"              uuid PRIMARY KEY NOT NULL,
  "name"            varchar          NOT NULL,
  "username"        varchar UNIQUE   NOT NULL,
  "email"           varchar UNIQUE   NOT NULL,
  "hashed_password" varchar          NOT NULL,
  "is_verified"     bool             NOT NULL DEFAULT FALSE,
  "created_at"      timestamptz      NOT NULL DEFAULT (now())
);
