CREATE TABLE clients(
    id   VARCHAR(36)  NOT NULL PRIMARY KEY,
    name VARCHAR(1024) NOT NULL,
    document VARCHAR(11) NOT NULL UNIQUE,
    birth_date DATE NOT NULL,
    created_at DATETIME NOT NULL
);

CREATE TABLE accounts(
    number   VARCHAR(8)  NOT NULL PRIMARY KEY,
    account_type VARCHAR(14) NOT NULL,
    account_status VARCHAR(128) NOT NULL,
    client_id  VARCHAR(36) NOT NULL,
    created_at DATETIME NOT NULL,
    FOREIGN KEY (client_id) REFERENCES clients(id)
);