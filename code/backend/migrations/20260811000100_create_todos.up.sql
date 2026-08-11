CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS todos (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  title text NOT NULL,
  is_completed boolean NOT NULL DEFAULT false,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT ck_todos_title_trimmed_not_blank CHECK (length(btrim(title)) BETWEEN 1 AND 200),
  CONSTRAINT ck_todos_title_is_trimmed CHECK (title = btrim(title))
);

CREATE INDEX IF NOT EXISTS idx_todos_created_at_id ON todos (created_at, id);
