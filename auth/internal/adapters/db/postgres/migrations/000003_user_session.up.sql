CREATE TABLE IF NOT EXISTS "sessions"
(
  "id"            uuid PRIMARY KEY NOT NULL,
  "user_id"       uuid             NOT NULL,
  "refresh_token" varchar UNIQUE   NOT NULL,
  "is_blocked"    boolean          NOT NULL DEFAULT false,
  "user_agent"    varchar          NOT NULL,
  "client_ip"     varchar          NOT NULL,
  "expires_at"    timestamptz      NOT NULL,
  "created_at"    timestamptz      NOT NULL DEFAULT (now())
);

ALTER TABLE "sessions"
  ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "sessions"
  ADD CONSTRAINT exp_gt_create CHECK ( sessions.expires_at > sessions.created_at );
