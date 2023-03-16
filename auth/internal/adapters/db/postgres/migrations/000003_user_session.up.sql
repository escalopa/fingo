CREATE TABLE IF NOT EXISTS "sessions"
(
  "id"            uuid PRIMARY KEY NOT NULL,
  "user_id"       uuid             NOT NULL,
  "access_token"  varchar UNIQUE   NOT NULL,
  "refresh_token" varchar UNIQUE   NOT NULL,
  "user_agent"    varchar          NOT NULL,
  "client_ip"     varchar          NOT NULL,
  "expires_at"    timestamptz      NOT NULL,
  "updated_at"    timestamptz      NOT NULL DEFAULT (now())
);

ALTER TABLE "sessions"
  ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "sessions"
  ADD CONSTRAINT exp_gt_create CHECK ( sessions.expires_at > sessions.updated_at );
