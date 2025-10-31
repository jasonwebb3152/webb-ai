CREATE TYPE "source_type" AS ENUM (
  'S3',
  'ONE_DRIVE'
);

CREATE TYPE "source_status" AS ENUM (
  'SCANNING',
  'IDLE'
);

CREATE TYPE "job_status" AS ENUM (
  'RUNNING',
  'FAILED',
  'COMPLETED'
);

CREATE TYPE "file_status" AS ENUM (
  'CREATED',
  'PROCESSING',
  'PROCESSED',
  'FAILED',
  'DELETED',
  'REMOVED'
);

CREATE TABLE "users" (
  "username" varchar(50) PRIMARY KEY,
  "hashed_password" varchar NOT NULL,
  "full_name" text NOT NULL,
  "email" text UNIQUE NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT (now())

  CHECK (char_length("username") >= 3)
);

CREATE TABLE "sessions" (
  "id" uuid PRIMARY KEY,
  "username" varchar NOT NULL,
  "refresh_token" text NOT NULL,
  "user_agent" text NOT NULL,
  "client_ip" text NOT NULL,
  "is_blocked" bool NOT NULL DEFAULT false,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "verify_emails" (
  "id" UUID PRIMARY KEY,
  "username" varchar NOT NULL,
  "email" text NOT NULL,
  "secret_code" text NOT NULL,
  "is_used" bool NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "expires_at" timestamptz NOT NULL DEFAULT (now() + interval '15 minutes')
);

CREATE TABLE "user_groups" (
  "id" text PRIMARY KEY,
  "name" text NOT NULL
);

CREATE TABLE "chats" (
  "id" UUID PRIMARY KEY,
  "username" varchar NOT NULL,
  "description" text,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "chat_messages" (
  "id" UUID PRIMARY KEY,
  "chat_id" UUID NOT NULL,
  "content" text NOT NULL,
  "is_ai_response" bool NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "sources" (
  "id" UUID PRIMARY KEY,
  "name" text NOT NULL,
  "description" text,
  "type" source_type NOT NULL,
  "status" source_status NOT NULL DEFAULT 'IDLE',
  "secrets_store" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "source_jobs" (
  "id" UUID PRIMARY KEY,
  "source_id" UUID NOT NULL,
  "status" job_status NOT NULL DEFAULT 'RUNNING',
  "file_count" int NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'

  CHECK ("file_count" > 0)
);

CREATE TABLE "s3_sources" (
  "source_type" text NOT NULL DEFAULT 'S3',
  "bucket_name" text NOT NULL,
  "region" text NOT NULL

  CHECK ("source_type" = 'S3')
) INHERITS ("sources");

CREATE TABLE "onedrive_sources" (
  "source_type" text NOT NULL DEFAULT 'ONE_DRIVE',
  "tenant_id" text NOT NULL

  CHECK ("source_type" = 'ONE_DRIVE')
) INHERITS ("sources");

CREATE TABLE "files" (
  "id" UUID PRIMARY KEY,
  "source_id" UUID NOT NULL,
  "url" text UNIQUE NOT NULL,
  "status" file_status NOT NULL DEFAULT 'CREATED',
  "last_processed_at" timestamptz,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "verify_emails" ("username");

CREATE INDEX ON "chats" ("username");

CREATE INDEX ON "chat_messages" ("chat_id", "created_at");

ALTER TABLE "verify_emails" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "chats" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "sessions" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "files" ADD FOREIGN KEY ("source_id") REFERENCES "sources" ("id");

ALTER TABLE "source_jobs" ADD FOREIGN KEY ("source_id") REFERENCES "sources" ("id");

ALTER TABLE "chat_messages" ADD FOREIGN KEY ("chat_id") REFERENCES "chats" ("id");
