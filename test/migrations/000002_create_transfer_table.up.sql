CREATE TABLE IF NOT EXISTS transfer (
    id                  TEXT PRIMARY KEY,
    origin_account_id   TEXT NOT NULL,
    destination_account_id TEXT NOT NULL,
    amount              INTEGER NOT NULL,
    created_at          DATETIME NOT NULL,
    FOREIGN KEY (origin_account_id) REFERENCES account (id),
    FOREIGN KEY (destination_account_id) REFERENCES account (id)
);
