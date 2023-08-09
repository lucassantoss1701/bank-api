CREATE TABLE account (
    id          VARCHAR(36) PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    cpf         VARCHAR(11) NOT NULL,
    secret      VARCHAR(100) NOT NULL,
    balance     INT NOT NULL,
    created_at  TIMESTAMP NOT NULL
);