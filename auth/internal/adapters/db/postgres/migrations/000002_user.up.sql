CREATE TABLE IF NOT EXISTS "users"
(
  "id"                  uuid PRIMARY KEY    NOT NULL,
  "first_name"          varchar(30)         NOT NULL,
  "last_name"           varchar(30)         NOT NULL,
  "username"            varchar UNIQUE      NOT NULL,
  "email"               varchar(100) UNIQUE NOT NULL,
  "hashed_password"     varchar             NOT NULL,
  "password_changed_at" timestamptz         NOT NULL DEFAULT (now()),
  "is_verified_email"   boolean             NOT NULL DEFAULT false,
  "created_at"          timestamptz         NOT NULL DEFAULT (now())
);
