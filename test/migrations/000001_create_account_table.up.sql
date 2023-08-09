CREATE TABLE IF NOT EXISTS account (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    cpf         TEXT NOT NULL,
    secret      TEXT NOT NULL,
    balance     INTEGER NOT NULL,
    created_at  DATETIME NOT NULL
);
