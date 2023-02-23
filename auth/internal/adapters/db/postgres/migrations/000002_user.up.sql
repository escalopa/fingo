CREATE TABLE "users"
(
  "id"                  uuid PRIMARY KEY NOT NULL,
  "username"            varchar UNIQUE   NOT NULL,
  "email"               varchar UNIQUE   NOT NULL,
  "hashed_password"     varchar          NOT NULL,
  "full_name"           varchar          NOT NULL,
  "password_changed_at" timestamptz      NOT NULL DEFAULT ('0001-01-01 00:00:00Z'),
  "created_at"          timestamptz      NOT NULL DEFAULT (now())
);
