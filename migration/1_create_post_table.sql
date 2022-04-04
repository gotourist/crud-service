CREATE TABLE IF NOT EXISTS post(
    id INTEGER,
    title varchar,
    body varchar,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT now()
)