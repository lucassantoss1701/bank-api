CREATE TABLE transfer (
    id                  VARCHAR(36) PRIMARY KEY,
    origin_account_id   VARCHAR(36) NOT NULL,
    destination_account_id VARCHAR(36) NOT NULL,
    amount              INT NOT NULL,
    created_at          TIMESTAMP NOT NULL,
    FOREIGN KEY (origin_account_id) REFERENCES account (id),
    FOREIGN KEY (destination_account_id) REFERENCES account (id)
);