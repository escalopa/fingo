CREATE TABLE IF NOT EXISTS "roles"
(
  "id"   SERIAL PRIMARY KEY NOT NULL,
  "name" varchar UNIQUE     NOT NULL

);

ALTER TABLE roles
  ADD CHECK ( name != '');

CREATE TABLE IF NOT EXISTS "user_roles"
(
  "id"         SERIAL PRIMARY KEY NOT NULL,
  "user_id"    uuid               NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
  "role_id"    INT                NOT NULL REFERENCES roles ("id") ON DELETE CASCADE,
  "created_at" timestamptz        NOT NULL DEFAULT (now())
);

ALTER TABLE "user_roles"
  ADD CONSTRAINT "user_roles_unique" UNIQUE ("user_id", "role_id");
