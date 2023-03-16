DO
$$
  BEGIN
    IF NOT EXISTS(SELECT 1 FROM pg_type WHERE typname = 'gender') THEN
      CREATE TYPE gender AS ENUM ('MALE', 'FEMALE');
    END IF;
  END
$$;

CREATE TABLE IF NOT EXISTS "users"
(
  "id"                  uuid PRIMARY KEY    NOT NULL,
  "first_name"          varchar(30)         NOT NULL,
  "last_name"           varchar(30)         NOT NULL,
  "username"            varchar UNIQUE      NOT NULL,
  "gender"              gender              NOT NULL,
  "email"               varchar(100) UNIQUE NOT NULL,
  "phone_number"        varchar(20) UNIQUE  NOT NULL,
  "hashed_password"     varchar             NOT NULL,
  "password_changed_at" timestamptz         NOT NULL DEFAULT (now()),
  "is_verified_email"   boolean             NOT NULL DEFAULT false,
  "is_verified_phone"   boolean             NOT NULL DEFAULT false,
  "created_at"          timestamptz         NOT NULL DEFAULT (now())
);
