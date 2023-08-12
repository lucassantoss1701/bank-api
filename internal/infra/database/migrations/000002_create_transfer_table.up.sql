CREATE TABLE IF NOT EXISTS transfer (
    id                  VARCHAR(36) PRIMARY KEY,
    origin_account_id   VARCHAR(36) NOT NULL,
    destination_account_id VARCHAR(36) NOT NULL,
    amount              INT NOT NULL,
    created_at          TIMESTAMP NOT NULL,
    CONSTRAINT fk_origin_account FOREIGN KEY (origin_account_id) REFERENCES account (id),
    CONSTRAINT fk_destination_account FOREIGN KEY (destination_account_id) REFERENCES account (id)
);